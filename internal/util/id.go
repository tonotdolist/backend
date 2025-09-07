package util

import "github.com/google/uuid"

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
