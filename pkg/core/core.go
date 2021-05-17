package core

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	// ID        uint       `json:"id" gorm:"primary_key"`
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4()"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `json:"deletedAt" sql:"index"`
}
