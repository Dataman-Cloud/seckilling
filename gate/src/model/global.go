package model

import (
	"github.com/satori/go.uuid"
)

const (
	SkCookie = "SK-UID"
)

var UIDMap map[string]bool

func init() {
	UIDMap = make(map[string]bool)
}

func NewUUID() string {
	return uuid.NewV4().String()
}
