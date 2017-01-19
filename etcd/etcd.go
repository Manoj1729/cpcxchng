package etcd

import (
	"fmt"
	"log"
	"strconv"
	"time"

	clientv3 "github.com/coreos/etcd/clientv3"
	cmn "github.com/ionosnetworks/cpcxchng/common"
	"golang.org/x/net/context"
)

const (
	ETCDCLIENTTIMEOUT = 30 * time.Second
)

var (
	EtcdAddress string
)

func SetEtcdAddress(etcdIP string) {
	EtcdAddress = etcdIP
}

func GetIonForCapacity(capacity int32) *cmn.Ion {

	if EtcdAddress == "" {
		return nil
	}
	ions := Get("/ION/Capacity/" + strconv.Itoa(int(capacity)))
	for _, ion := range ions {
		if ion.Capacity == int(capacity) {
			return &ion
		}
	}
	return nil
}

func Set(key, value string, timeout int64) {

	dialTimeout := time.Duration(ETCDCLIENTTIMEOUT)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{EtcdAddress},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	if timeout != 0 {
		// minimum lease TTL is 5-second
		resp, err1 := cli.Grant(context.TODO(), timeout)
		if err1 != nil {
			fmt.Println("ETCD :: Failed to grante lease ", err1)
			return
		}
		// fmt.Println("setting it for ", timeout)
		_, err = cli.Put(context.TODO(), key, value, clientv3.WithLease(resp.ID))
	} else {
		// fmt.Println("setting it for ever")
		_, err = cli.Put(context.TODO(), key, value)
	}

	if err != nil {
		log.Fatal("Put error :", err)
	}
}

func Get(key string) []cmn.Ion {

	dialTimeout := time.Duration(ETCDCLIENTTIMEOUT)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{EtcdAddress},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	requestTimeout := time.Duration(ETCDCLIENTTIMEOUT)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatal("Get Error", err)
	}
	var ret []cmn.Ion
	for _, ev := range resp.Kvs {
		var ion cmn.Ion
		cmn.ConvertByteArrayToObject(ev.Value, &ion)
		ret = append(ret, ion)
	}

	return ret
}

func Del(key string) {

	dialTimeout := time.Duration(ETCDCLIENTTIMEOUT)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{EtcdAddress},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	// delete the keys
	requestTimeout := time.Duration(ETCDCLIENTTIMEOUT)
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	_, err = cli.Delete(ctx, key, clientv3.WithPrefix())
	cancel()
	if err != nil {
		log.Fatal(err)
	} else {
		//fmt.Println(dresp)
	}
}
