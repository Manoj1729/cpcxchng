package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"

	cmn "github.com/ionosnetworks/cpcxchng/common"
)

type client cmn.Client

func main() {
	newcli := NewClient("127.0.0.1", "7777")
	err := newcli.Connect()
	if err != nil {
		os.Exit(1)
	}
	cap := cmn.Capacity{Cap: 30}
	err = newcli.RequestForIonWithCapacity(cap)
	if err != nil {
		newcli.Close()
		os.Exit(1)
	}
	go newcli.ReadIncomingData()
	newcli.ProcessIncomingMsg()
	newcli.Close()

}

func NewClient(serverip, serverport string) *client {
	cli := new(client)
	cli.ServIP = serverip
	cli.ServPort = serverport
	cli.MsgChan = make(chan cmn.Msg)
	return cli
}

func (clien *client) Connect() error {
	serveraddr := clien.ServIP + ":" + clien.ServPort
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", serveraddr, conf)
	if err != nil {
		fmt.Println("Error connecting to server", err)
		return err
	}
	clien.Conn = cmn.ClientConn{Conn: conn, Lock: &sync.Mutex{}}
	return nil
}

func (clien *client) RequestForIonWithCapacity(cap cmn.Capacity) error {
	capa := cmn.Capacity{Cap: cap.Cap}
	mssg := cmn.Msg{Header: cmn.CAPACITY_REQ, Data: capa}
	data, err := mssg.Encode()
	fmt.Println("created msg", mssg, data)
	if err != nil {
		fmt.Println("msg encode error", err)
		return err
	}
	clien.Conn.Conn.Write(data)
	return nil
}

func (clien *client) ReadIncomingData() error {
	data, err := ioutil.ReadAll(clien.Conn.Conn)
	//for {
	if err != nil {
		fmt.Println("error reading data", err, *clien)
	}
	resp := new(cmn.Msg)
	err = resp.Decode(data)
	if err != nil {
		fmt.Println("Message Decode Failed", err)
		clien.Conn.Conn.Close()
		return err
	}
	clien.MsgChan <- *resp
	//conn.Close()
	//}

	return nil
}

func (clien *client) ProcessIncomingMsg() error {
	resp := <-clien.MsgChan
	switch resp.Header {
	case cmn.CAPACITY_RESP:
		iondata := resp.Data
		ion := iondata.(*cmn.Ion)
		fmt.Println("read data", ion)
	}
	return nil
}

func (clien *client) Close() error {
	return clien.Conn.Conn.Close()
}
func readIncomingData(conn net.Conn, ch chan cmn.Msg) {
	data, err := ioutil.ReadAll(conn)
	//for {
	if err != nil {
		fmt.Println("error reading data", conn)
	}
	resp := new(cmn.Msg)
	resp.Decode(data)
	if err != nil {
		fmt.Println("Wrong message")
		conn.Close()
		return
	}
	ch <- *resp
	conn.Close()
	//}
}

func processIncomingMsg(ch chan cmn.Msg) {
	resp := <-ch
	switch resp.Header {
	case cmn.CAPACITY_RESP:
		//iondata := resp.Data
		//ion := iondata.(cmn.Ion)
		fmt.Println("read data", resp)
	}
}

func requestForIonWithCapacity(cap int32, conn *tls.Conn) {
	capa := cmn.Capacity{Cap: cap}
	mssg := cmn.Msg{Header: cmn.CAPACITY_REQ, Data: capa}
	data, err := mssg.Encode()
	fmt.Println("created msg", mssg, data)
	if err != nil {
		fmt.Println("msg encode error", err)
		return
	}
	conn.Write(data)
	return
}
