package common

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
)

func CheckErr(err erro) {

}

func (msg *Msg) Encode() ([]byte, error) {
	w := new(bytes.Buffer)
	err := binary.Write(w, binary.LittleEndian, &msg.Header)
	if err != nil {
		return nil, err
	}
	encoder := gob.NewEncoder(w)
	switch msg.Header {
	case CAPACITY_REQ:
		cap := msg.Data.(Capacity)
		err = encoder.Encode(cap)
	case CAPACITY_RESP:
		ion := msg.Data.(Ion)
		err = encoder.Encode(ion)
	}
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil

	return nil, nil
}

func (msg *Msg) Decode(buf []byte) error {
	var err error
	w := bytes.NewBuffer(buf)
	err = binary.Read(w, binary.LittleEndian, &msg.Header)
	if err != nil {
		return err
	}
	decoder := gob.NewDecoder(w)
	fmt.Println("msg id", msg.Header)
	switch msg.Header {
	case CAPACITY_REQ:
		cap := new(Capacity)
		err = decoder.Decode(cap)
		if err != nil {
			return err
		}
		fmt.Println("capacity", cap)
		msg.Data = cap
	case CAPACITY_RESP:
		ion := new(Ion)
		err = decoder.Decode(ion)
		if err != nil {
			fmt.Println("decode error", err)
			return err
		}
		fmt.Println("capacity", ion)
		msg.Data = ion
	}
	return err

}
