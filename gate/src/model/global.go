package model

const (
	SkCookie = "DM_SK_UID"

	EventListKey = "events"

	CurrentEventKey = "cur_eid"

	EventIdKey = "sn:%s"

	OrderKey = "tr:%s:%s" //(eid, sn)

	EventInfoKey = "event:%s"

	WorkOffIndexKey = "wf:%s"
)

var (
	CurrentEventId string
)
