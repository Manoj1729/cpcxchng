package common

type ReqType int
type MsgHeader int32

const (
	REQUEST ReqType = iota
	RESPONSE
)

const (
	CAPACITY_REQ MsgHeader = iota
	CAPACITY_RESP
)

type Msg struct {
	Header MsgHeader
	Data   interface{}
}

type Capacity struct {
	Cap int32
}

type Ion struct {
	IP       string
	Port     string
	Capacity int
	SyncID   int32
	ID       int
}

type erro struct {
	Level string
	Err   error
}
