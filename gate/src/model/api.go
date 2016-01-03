package model

type UserInfo struct {
	UID       string `json:"uid"`
	Phone     string `json:"phone"`
	EventId   string `json:"eventId"`
	Timestamp int64  `json:"timestamp"`
	Index     int64  `json:"index"`
	SerialNum string `json:"serialNum"`
}

type EventInfo struct {
	Id       string   `json:"id"`
	EffectOn int64    `json:"effectOn"`
	Duration int64    `json:"duration"`
	Describe string   `json:"describe"`
	Resource []string `json:"resource"`
	// 1 is under way
	// 2 is has not started
	// 3 is over
	Status int64 `json:"status"`
}
