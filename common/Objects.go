package common

import (
	"net"
	"sync"
)

type MsgHeader int32
type ConnState int32

const (
	CAPACITY_REQ MsgHeader = iota
	CAPACITY_RESP
)

const (
	NOT_VALID ConnState = iota
	CLOSED
	CONNECTED
)

type Msg struct {
	Header MsgHeader
	Data   interface{}
}

type Capacity struct {
	Cap int32
}

type Server struct {
	IP         string
	Port       string
	Type       string
	Clients    map[string]Client
	ClientLock *sync.RWMutex
	EtcdIP     string
}

type ClientConn struct {
	Conn net.Conn
	Lock *sync.Mutex
}

type Client struct {
	ServIP   string
	ServPort string
	Conn     ClientConn
	State    ConnState
	MsgChan  chan Msg
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
