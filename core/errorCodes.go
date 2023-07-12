package core

type ErrCode int

type ErrorMsg struct {
	Code    ErrCode
	Message string
}
