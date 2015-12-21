package model

import (
	"time"
)

type Order struct {
	EId    int64
	UId    string
	Status int64
	Seq    int64
	Ext    string
	Create time.Time
}
