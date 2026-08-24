package model

import "uuid"

type Category struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
