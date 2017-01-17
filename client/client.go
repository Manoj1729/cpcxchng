package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	cmn "github.com/ionosnetworks/cpcxchng/common"
)

func main() {
	msgchan := make(chan cmn.Msg)
	service := "127.0.0.1:7777"
	//tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	//checkError(err)
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := tls.Dial("tcp", service, conf)
	if err != nil {
		fmt.Println("Error connecting", err)
		return
	}
	requestForIonWithCapacity(30, conn)
	go readIncomingData(conn, msgchan)
	go processIncomingMsg(msgchan)
	time.Sleep(20 * time.Second)
	//result, err := ioutil.ReadAll(conn)
	//if err != nil {
	//	fmt.Println("result read error", err)
	//}
	//fmt.Println(string(result))
	//conn.Close()
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

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
