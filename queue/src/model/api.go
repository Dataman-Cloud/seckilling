package model

type CommonResponse struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

type CountdownData struct {
	CurTime  int64 `json:"curTime"`
	UnlockOn int64 `json:"unlockOn"`
	Locked   bool  `json:"locked"`
}

type TicketData struct {
	UID       string `json:"uid"`
	Timestamp int64  `json:"timestamp"`
}
