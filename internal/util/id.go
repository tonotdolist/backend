package util

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

func NewID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), err
}

func ValidateID(id string) bool {
	return uuid.Validate(id) == nil
}

var ValidID validator.Func = func(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(string)
	if ok {
		return ValidateID(id)
	}
	return false
}
