package model

import (
	"github.com/satori/go.uuid"
)

const (
	SkCookie = "SK-UID"
)

func NewUUID() string {
	return uuid.NewV4().String()
}
