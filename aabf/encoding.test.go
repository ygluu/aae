package aabf

import (
	"hash/crc32"
)

type SubMsg struct {
	Msg
	S  string
	AS []string
	P1 int32
	P2 uint32
}

type TestMsg struct {
	Msg
	B    bool
	AB   []bool
	I8   int8
	AI8  []int8
	I16  int16
	AI16 []int16
	I32  int32
	AI32 []int32
	I64  int64
	AI64 []int64
	AU8  []uint8
	U8   uint8
	U16  uint16
	AU16 []uint16
	U32  uint32
	AU32 []uint32
	U64  uint64
	AU64 []uint64

	S    string
	F32  float32
	AS   []string
	F64  float64
	M    *SubMsg
	AF32 []float32
	AM   []*SubMsg
	AF64 []float64
}

type TestM struct {
	M
}

func (m *TestM) Ready() {
	m.AsynOper(
		func() {
			Log.I("TestM asynoper execing: %d", CurrGoID())
		},
		func() {
			Log.I("TestM asynoper execed: %d", CurrGoID())
		},
		"",
	)

	m.Timer().Add(
		3000, 1000,
		func() {
			Log.I("TestM timer event: %d", CurrGoID())
		},
		func() {
			Log.I("TestM timer die: %d", CurrGoID())
		},
	)
}

func (m *TestM) OnTestMsg(msg *TestMsg) {
	Log.I("TestM.OnTestMsg: %d, %s", m.GID(), ToJson(msg))
}

type TestG struct {
	G
}

func (g *TestG) OnTestMsg(msg *TestMsg) {
	Log.I("TestGo.OnTestMsg: %d, %s", g.GID(), ToJson(msg))
}

func (g *TestG) Ready() {
	g.AsynOper(
		func() {
			Log.I("TestGo asynoper execing: %d", CurrGoID())
		},
		func() {
			Log.I("TestGo asynoper execed: %d", CurrGoID())
		},
		"",
	)

	g.Timer().Add(
		3000, 1000,
		func() {
			Log.I("TestGo timer event: %d", CurrGoID())
		},
		func() {
			Log.I("TestGo timer die: %d", CurrGoID())
		},
	)
}

type test1 struct {
	AU8 []uint8
}

func init() {
	RegModule(&TestM{})
	RegModule(&TestG{})
}

func testEncoding() {
	tm := &TestMsg{
		B:    true,
		AB:   []bool{true, false, true},
		I8:   -8,
		AI8:  []int8{1, 2, -3, -4},
		I16:  -16,
		AI16: []int16{0x1601, 0x1602, -0x1603, -0x1604},
		I32:  -32,
		AI32: []int32{0x3201, 0x3202, -0x3203, -0x3204},
		I64:  -64,
		AI64: []int64{0x6401, 0x6402, -0x6403, -0x6404},
		U8:   8,
		AU8:  []uint8{1, 2, 3, 4},
		U16:  16,
		AU16: []uint16{0x1601, 0x1602, 0x1603, 0x1604},
		U32:  32,
		AU32: []uint32{0x3201, 0x3202, 0x3203, 0x3204},
		U64:  64,
		AU64: []uint64{0x6401, 0x6402, 0x6403, 0x6404},
		F32:  32.32,
		S:    "这是测试信息",
		AS:   []string{"这是测试信息1", "", "这是test信息3", "这是test信息4"},
		AF64: []float64{64.31, 64.32, 64.33, 64.34},
		M:    &SubMsg{S: "This is sub msg", AS: []string{"这是子信息1", "这是子信息2"}, P1: -123, P2: 123},
		AF32: []float32{32.31, 32.32, 32.33, 32.34},
		AM: []*SubMsg{&SubMsg{S: "", AS: []string{"这是数组子信息1", ""}, P1: -123, P2: 123},
			&SubMsg{S: "This is array sub msg", AS: []string{"这是数组子信息1", "这是数组子信息2"}, P1: -123, P2: 123}},
		F64: 64.32,
	}

	datas := make([]byte, 10*1024*1024)

	md := Marshal(tm, datas)
	if md == nil {
		Log.I("Marshal error")
		return
	}

	mi := Unmarshal(md)
	if mi == nil {
		Log.I("Unmarshal error")
		return
	}

	c1 := crc32.ChecksumIEEE(md)
	md = Marshal(mi, datas)
	if md == nil {
		Log.I("Remarshal error")
	}

	c2 := crc32.ChecksumIEEE(md)
	if c1 != c2 {
		Log.E("Marshall and Unmarshal errors")
		return
	}

	MainG.PostMsg(mi)
}
