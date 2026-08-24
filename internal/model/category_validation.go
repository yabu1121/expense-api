package model

import "strings"

func (c *Category) Normalize() {
	trimmedName := strings.TrimSpace(c.Name)
	trimmedName = strings.ToLower(trimmedName)
	c.Name = trimmedName
}

func (c Category) Validate() error {
	if err := c.validateName(); err != nil {
		return err
	}
	return nil
}

func (c Category) validateName() error {
	if c.Name == "" {
		return ErrNameRequired
	}
	return nil
}
