package aabf

import (
	"reflect"
)

type fieldInfo struct {
	name string
	len  int
}

var ftnameByKind = make(map[reflect.Kind]string)
var fkindByName = make(map[string]reflect.Kind)

func init() {
	ftnameByKind[reflect.Bool] = "b"
	ftnameByKind[reflect.Int8] = "i8"
	ftnameByKind[reflect.Int16] = "i16"
	ftnameByKind[reflect.Int32] = "i32"
	ftnameByKind[reflect.Int64] = "i64"
	ftnameByKind[reflect.Uint8] = "u8"
	ftnameByKind[reflect.Uint16] = "u16"
	ftnameByKind[reflect.Uint32] = "u32"
	ftnameByKind[reflect.Uint64] = "u64"
	ftnameByKind[reflect.Float32] = "f32"
	ftnameByKind[reflect.Float64] = "f64"
	ftnameByKind[reflect.String] = "s"
	ftnameByKind[reflect.Struct] = "m"

	fkindByName["b"] = reflect.Bool
	fkindByName["i8"] = reflect.Int8
	fkindByName["i16"] = reflect.Int16
	fkindByName["i32"] = reflect.Int32
	fkindByName["i64"] = reflect.Int64
	fkindByName["u8"] = reflect.Uint8
	fkindByName["u16"] = reflect.Uint16
	fkindByName["u32"] = reflect.Uint32
	fkindByName["u64"] = reflect.Uint64
	fkindByName["f32"] = reflect.Float32
	fkindByName["f64"] = reflect.Float64
	fkindByName["s"] = reflect.String
	fkindByName["m"] = reflect.Struct
}

type MsgI interface {
	GetId() MID
	SetId(id MID)
	GetH() H
	SetH(h H)
}

type Msg struct {
	id MID
	h  H
}

func (m *Msg) LugoMsg() {
}

func (m *Msg) GetH() H {
	return m.h
}

func (m *Msg) SetH(h H) {
	m.h = h
}

func (m *Msg) GetId() MID {
	return m.id
}

func (m *Msg) SetId(id MID) {
	m.id = id
}
