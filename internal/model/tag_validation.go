package model

import "strings"

func (t *Tag) Normalize() {
	trimmedName := strings.TrimSpace(t.Name)
	trimmedName = strings.ToLower(trimmedName)
	t.Name = trimmedName
}

func (t Tag) Validate() error {
	if err := t.validateName(); err != nil {
		return err
	}
	return nil
}

func (t Tag) validateName() error {
	if t.Name == "" {
		return ErrNameRequired
	}
	return nil
}
