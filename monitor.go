package ssstatus

import (
	"fmt"
	"os"
	"time"

	"github.com/SoulSec/ssstatus/dial"
	"github.com/SoulSec/ssstatus/logger"
	"github.com/SoulSec/ssstatus/notify"
	"github.com/SoulSec/ssstatus/track"
)

type Monitor struct {
	config              *Config
	checkerV            chan *Server
	notifiers           notify.Notifiers
	notifierV           chan *Server
	notificationTracker map[*Server]*track.TimeTracker
	dialer              *dial.Dialer
	stop                chan struct{}
	serverStatusData    *ServerStatusData
}

func NewMonitor(c *Config) *Monitor {
	z := &Monitor{
		config:              c,
		checkerV:            make(chan *Server),
		notifiers:           c.Settings.Notifications.GetNotifiers(),
		notifierV:           make(chan *Server),
		notificationTracker: make(map[*Server]*track.TimeTracker),
		dialer:              dial.NewDialer(c.Settings.Monitor.MaxConnection),
		stop:                make(chan struct{}),
		serverStatusData:    NewServerStatusData(c.Servers),
	}
	z.initialize()
	return z
}

func (z *Monitor) initialize() {
	for _, notifier := range z.notifiers {
		if initializer, ok := notifier.(notify.Initializer); ok {
			logger.Logln("Initializing", initializer)
			initializer.Initialize()
		}
	}
	for _, server := range z.config.Servers {
		z.notificationTracker[server] = NewTrackerWithExpBackoff(z.config.Settings.Monitor.ExponentialBackoffSeconds)
		switch {
		case server.CheckInterval <= 0:
			server.CheckInterval = z.config.Settings.Monitor.CheckInterval
		case server.Timeout <= 0:
			server.Timeout = z.config.Settings.Monitor.Timeout
		}
	}
}

func NewTrackerWithExpBackoff(expBackoffSeconds int) *track.TimeTracker {
	return track.NewTracker(track.NewExpTrack(expBackoffSeconds))
}

func (z *Monitor) Run() {
	z.RunForSeconds(0)
}

func (z *Monitor) RunForSeconds(runningSeconds int) {
	if runningSeconds != 0 {
		go func() {
			runningSecondsTime := time.Duration(runningSeconds) * time.Second
			<-time.After(runningSecondsTime)
			z.stop <- struct{}{}
		}()
	}
	for _, server := range z.config.Servers {
		go z.scheduleServer(server)
	}
	logger.Logln("Starting monitor.")
	z.monitor()
}

func (z *Monitor) scheduleServer(s *Server) {
	z.checkerV <- s
	tickerSeconds := time.NewTicker(time.Duration(s.CheckInterval) * time.Second)
	for range tickerSeconds.C {
		z.checkerV <- s
	}
}

func (z *Monitor) monitor() {
	go z.listenForChecks()
	go z.listenForNotifications()

	<-z.stop
	logger.Logln("Terminating.....")
	os.Exit(0)
}

func (z *Monitor) listenForChecks() {
	for server := range z.checkerV {
		z.checkServerStatus(server)
	}
}

func (z *Monitor) checkServerStatus(server *Server) {
	worker, output := z.dialer.NewWorker()
	go func() {
		logger.Logln("Checking", server)
		formattedAddress := fmt.Sprintf("%s:%d", server.IPAddress, server.Port)
		timeoutSeconds := time.Duration(server.Timeout) * time.Second
		worker <- dial.NetAddressTimeout{NetAddress: dial.NetAddress{Network: server.Protocol, Address: formattedAddress}, Timeout: timeoutSeconds}
		dialerStatus := <-output
		z.serverStatusData.SetStatusAtTimeForServer(server, time.Now(), dialerStatus.Ok)
		if !dialerStatus.Ok {
			logger.Logln(dialerStatus.Err)
			logger.Logln("ERROR", server)
			go func() {
				z.notifierV <- server
			}()
			return
		}
		logger.Logln("OK", server)
		if z.notificationTracker[server].HasBeenRan() {
			z.notificationTracker[server] = NewTrackerWithExpBackoff(z.config.Settings.Monitor.ExponentialBackoffSeconds)
		}
	}()
}

func (z *Monitor) listenForNotifications() {
	for server := range z.notifierV {
		timeTracker := z.notificationTracker[server]
		if timeTracker.IsReady() {
			nextDelay, nextTime := timeTracker.SetNext()
			logger.Logln("Sending notifications for", server)
			go z.notifiers.NotifyAll(fmt.Sprintf("%s (%s)", server.Name, server))
			logger.Logln("Next available notification for", server.String(), "in", nextDelay, "at", nextTime)
		}
	}
}
