package model

const (
	SkCookie = "DM_SK_UID"

	EventListKey = "events"

	CurrentEventKey = "cur_eid"

	EventIdKey = "SN:%s"

	OrderKey = "TR:%s:%s" //(eid, sn)

	EventInfoKey = "event:%s"
)

var (
	CurrentEventId string
)
