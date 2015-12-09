package model

type CommonResponse struct {
	Code  int
	Data  interface{}
	Error string
}

type CountdownData struct {
	CurTime  int64
	UnlockOn int64
	Locked   bool
}
