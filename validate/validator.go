package validate

type Validator interface {
	Validate() error
}

func ValidateAll(validators ...Validator) error {
	for _, validator := range validators {
		if err := validator.Validate(); err != nil {
			return err
		}
	}
	return nil
}
