package util

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type IDProvider interface {
	NewID() (string, error)
	ValidateID(id string) bool
	Validator() func(fl validator.FieldLevel) bool
}

type UUIDProvider struct {
}

func NewUUIDProvider() IDProvider {
	return &UUIDProvider{}
}

func (up *UUIDProvider) NewID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}

	return id.String(), err
}

func (up *UUIDProvider) ValidateID(id string) bool {
	return uuid.Validate(id) == nil
}

func (up *UUIDProvider) Validator() func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		id, ok := fl.Field().Interface().(string)
		if ok {
			return up.ValidateID(id)
		}
		return false
	}
}
