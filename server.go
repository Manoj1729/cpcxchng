package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"

	cmn "github.com/ionosnetworks/cpcxchng/common"
	etcd "github.com/ionosnetworks/cpcxchng/etcd"
)

type serv cmn.Server

const (
	CERT = "server.crt"
	KEY  = "server.key"
)

func main() {
	server := Create()
	err := server.Start()
	if err != nil {
		fmt.Println("Error starting server", err)
	}
}

func Create() *serv {

	server := new(serv)
	server.Clients = make(map[string]cmn.Client)
	server.ClientLock = &sync.RWMutex{}
	server.readParams()
	return server
}

func (server serv) readParams() {

	if ip := os.Getenv("IP_TO_USE"); ip != "" {
		server.IP = ip
	}

	if port := os.Getenv("CPC_PORT"); port != "" {
		server.Port = port
	} else {
		server.Port = "3000"
	}

	if etcdIP := os.Getenv("ETCD_IP"); etcdIP != "" {
		server.EtcdIP = etcdIP
		etcd.SetEtcdAddress(server.EtcdIP)
	} else {
		server.EtcdIP = ""
	}

	if servtype := os.Getenv("SERV_TYPE"); servtype != "" {
		server.Type = servtype
	} else {
		server.Type = "tcp"
	}

}

func (serv *serv) Start() error {
	var servaddr string
	if serv.IP != "" {
		servaddr = serv.IP + ":" + serv.Port
	} else {
		servaddr = "127.0.0.1:7777" //":" + serv.Port
	}

	config := cmn.TlsConfig(CERT, KEY)
	ln, err := tls.Listen("tcp", servaddr, config)
	if err != nil {
		// handle error
		fmt.Println("Not able to listen on port", serv.Port)
		os.Exit(1)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			// TODO:: This is a serious error.
		}
		newconn := cmn.ClientConn{Conn: conn, Lock: &sync.Mutex{}}
		go serv.handleClient(newconn)
	}

	return nil
}

func (server *serv) handleClient(conn cmn.ClientConn) {
	conn.Conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
	defer conn.Conn.Close()
	request := make([]byte, 1024) // set maximum request length to 128B to prevent flood based attacks
	_, err := conn.Conn.Read(request)
	if err != nil {
		fmt.Println("conn read error", err)
	}
	msg := new(cmn.Msg)
	err = msg.Decode(request)
	if err != nil {
		fmt.Println("msg decode error", err)
	}
	switch msg.Header {
	case cmn.CAPACITY_REQ:
		inf := msg.Data
		cap := inf.(*cmn.Capacity)
		ion, err := findIonForCapacity(*cap)
		if err != nil {
			fmt.Println("Error Getting ion for capacity")
		}
		respmsg := cmn.Msg{Header: cmn.CAPACITY_RESP, Data: ion}
		respdata, err := respmsg.Encode()
		if err != nil {
			fmt.Println("Error encoding ion data")
		}
		conn.Conn.Write(respdata)
	}

}

func findIonForCapacity(cap cmn.Capacity) (cmn.Ion, error) {
	ion := etcd.GetIonForCapacity(cap.Cap) //cmn.Ion{IP: "192.168.1.164", Port: "3030", Capacity: 30, SyncID: 1, ID: 1}
	return *ion, nil
}
