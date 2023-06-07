package main

import (
	"bytes"
	"encoding/binary"
)

const FORM_ORDER_ID = 2
const CONTROL_ORDER_ID = 3

type Order interface {
	Bytes() []byte //FIXME: Add read? And pass as a pointer
}

type FormOrder struct {
	Kind    string  `json:"kind"`
	Payload float64 `json:"payload"`
}

func (order FormOrder) Bytes() []byte {
	idBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(idBuf, FORM_ORDER_ID)
	kindBuf := []byte(order.Kind)
	payloadBuf := make([]byte, 8)
	binary.LittleEndian.PutUint64(payloadBuf, uint64(order.Payload))

	resultBuf := append(idBuf, kindBuf...)
	return append(resultBuf, payloadBuf...)
}

func (order *FormOrder) Read(data []byte) {
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, order)
}

type ControlOrder struct {
	Id    uint8 `json:"id"`
	State bool  `json:"state"`
}

func (order ControlOrder) Bytes() []byte {
	head := make([]byte, 2)
	binary.LittleEndian.PutUint16(head, CONTROL_ORDER_ID)

	var booleanValue uint8
	if order.State {
		booleanValue = 1
	} else {
		booleanValue = 0
	}
	return append(head, order.Id, booleanValue)
}

func (order *ControlOrder) Read(data []byte) {
	reader := bytes.NewReader(data)
	binary.Read(reader, binary.LittleEndian, order)
}
