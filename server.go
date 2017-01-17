package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	cmn "github.com/ionosnetworks/cpcxchng/common"
)

func main() {
	service := "127.0.0.1:7777"
	cer, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Println(err)
		return
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{cer},
		ClientAuth:   tls.NoClientCert,
	}
	listener, err := tls.Listen("tcp", service, config) //net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go handleClient(conn)
	}

}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		//os.Exit(1)
	}
}

func handleClient(conn net.Conn) {
	//conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
	defer conn.Close() // close connection before exit
	//_, ok := conn.(*tls.Conn)
	request := make([]byte, 1024) // set maximum request length to 128B to prevent flood based attacks
	//if ok {
	//	for {
	_, err := conn.Read(request)
	checkError(err)
	msg := new(cmn.Msg)
	err = msg.Decode(request)
	if err != nil {
		//fmt.Println("msg decode error", err)
		//break
	}
	fmt.Println("read request", msg.Header)
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
		//daytime := strconv.FormatInt(time.Now().Unix(), 10)
		conn.Write(respdata)
	}
	time.Sleep(5 * time.Second)
	//}
	//}
}

func findIonForCapacity(cap cmn.Capacity) (cmn.Ion, error) {
	ion := cmn.Ion{IP: "192.168.1.164", Port: "3030", Capacity: 30, SyncID: 1, ID: 1}
	return ion, nil
}
