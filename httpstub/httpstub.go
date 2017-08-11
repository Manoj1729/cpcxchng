// This package provides a REST API frontend to bolt db/leveldb
package httpstub

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ionosnetworks/qfx_cmn/blog"
	kr "github.com/ionosnetworks/qfx_cmn/keyreader"

	"github.com/gorilla/mux"
	o "github.com/ionosnetworks/qfx_cmn/capcli/common"
	etcd "github.com/ionosnetworks/qfx_cp/capmgr/etcd"
)

const (
	CRT = "/keys/logsvr.crt" //"/keys/logsvr.crt"
	KEY = "/keys/logsvr.key" //"/keys/logsvr.key"
)

func init() {
	key := kr.New(KEY)
	var logger blog.Logger
	if logger = blog.New("127.0.0.1:2000", key.Key, key.Secret); logger == nil {
		fmt.Println("Logger failed ")
		return
	}

	//etcd.InitEtcd("127.0.0.1:2379", "testing", logger)

}

func getNextHopRegion(w http.ResponseWriter, req *http.Request) {
	var reqData o.RouteReqIn
	var rd o.RouteInfoResp
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		fmt.Println("error reading body from request", err.Error())
	} else if err = json.Unmarshal(body, &reqData); err != nil {
		fmt.Println("error unmarshalling data", err.Error())
	}
	routeData := o.RouteInfo{NextHopId: "127.0.0.1"}
	rd.RouteInfo = routeData
	apiresp := o.ApiResp{ErrorCode: "NONE"}
	rd.ApiResp = apiresp
	//sd.SrcCsId = src
	//sd.DstCsIds = dsts
	//sd.BW, _ = strconv.Atoi(bw)
	if err = json.NewEncoder(w).Encode(rd); err != nil {
		fmt.Println("error encoding json data", err.Error(), rd)
	}

}

func getSyncInfo(w http.ResponseWriter, req *http.Request) {
	fmt.Println("yo buddy")
	//var sd cmn.SyncDetail
	//var err error
	var reqData o.SyncReqIn
	var sd o.SyncReqOut
	var dstcps []o.DstCpeSyncInfo
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, 1048576))
	if err != nil {
		fmt.Println("error reading body from request", err.Error())
	} else if err = json.Unmarshal(body, &reqData); err != nil {
		fmt.Println("error unmarshalling data", err.Error())
	}
	dsts := etcd.GetDestsForSync(reqData.SyncId)
	src := etcd.GetSrcForSync(reqData.SyncId)
	bw := etcd.GetBWforSync(reqData.SyncId)
	if src == "" || len(dsts) == 0 {
		err = errors.New("sync info not found")
	} else if bw == "" {
		err = errors.New("band width not found for sync")
	}
	for _, dst := range dsts {
		dstcp := o.DstCpeSyncInfo{DstCpeId: dst}
		dstcps = append(dstcps, dstcp)
	}
	syncRel := o.SyncReln{SrcCpeId: src, DstCpesInfo: dstcps}
	sd.SyncReln = syncRel
	apiresp := o.ApiResp{ErrorCode: "NONE"}
	sd.ApiResp = apiresp
	//sd.SrcCsId = src
	//sd.DstCsIds = dsts
	//sd.BW, _ = strconv.Atoi(bw)
	if err = json.NewEncoder(w).Encode(sd); err != nil {
		fmt.Println("error encoding json data", err.Error(), sd)
	}
	fmt.Println("sync info", sd, src, dsts)

}

func configTLS() *tls.Config {
	var caCertPool *x509.CertPool = nil
	// Load CA cert
	/*if httpcli.caFile != "" {
		caCert, err := ioutil.ReadFile(httpcli.caFile)
		if err == nil {

			fmt.Println("Loaded CA certificates")
			caCertPool = x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)
		}
	}*/

	TLSConfig := &tls.Config{RootCAs: caCertPool}
	TLSConfig.Rand = rand.Reader
	TLSConfig.MinVersion = tls.VersionTLS10
	TLSConfig.SessionTicketsDisabled = false
	TLSConfig.InsecureSkipVerify = true

	return TLSConfig

}

func Start() {
	router := mux.NewRouter()
	router.HandleFunc("/slc/v1/getsyncinfo", getSyncInfo).Methods("POST")
	router.HandleFunc("/slc/v1/getnexthop", getNextHopRegion).Methods("POST")
	//go func() {
	s := &http.Server{
		Addr:           ":12345",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig:      configTLS(),
	}
	fmt.Println("Starting capdisc http stub...")
	//s.ListenAndServe()
	err := s.ListenAndServeTLS(CRT, KEY)
	panic(err)
	//}()
}
