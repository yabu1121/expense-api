package model

import "uuid"

type Tag struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
