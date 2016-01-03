package model

type UserInfo struct {
	UID       string `json:"uid"`
	Phone     string `json:"phone"`
	EventId   string `json:"eventId"`
	Timestamp int64  `json:"timestamp"`
	Index     int64  `json:"index"`
	SerialNum string `json:"serialNum"`
}
