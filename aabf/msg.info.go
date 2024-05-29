package aabf

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"reflect"
	"strings"
	"sync"
)

type codfunc func(addr uintptr, bytes []byte, index int)

type finfo struct {
	fname  string
	tname  string
	isarr  bool
	kind   reflect.Kind
	offset uintptr
	submi  *minfo
	mfunc  codfunc
	unfunc codfunc
}

// msg info

type minfo struct {
	proto   bool
	client  bool
	name    string
	id      MID
	sminfos sminfos
	finfos  []*finfo
	mt      reflect.Type
}

func (i *minfo) toJson() string {
	ret := "{"
	for i, fi := range i.finfos {
		if i == 0 {
			ret = ret + fmt.Sprintf("{\"%s\",\"%s\"}", fi.fname, fi.tname)
		} else {
			ret = ret + "," + fmt.Sprintf("{\"%s\",\"%s\"}", fi.fname, fi.tname)
		}
	}
	ret = ret + "}"
	return ret
}

var minfoByName Map[string, *minfo]
var minfoById Map[MID, *minfo]
var minfoByType Map[string, *minfo]
var msgMutex sync.Mutex

func StrToMd5(str string) string {
	h := md5.New()
	io.WriteString(h, str)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func StrToId(str string) MID {
	return MID(crc32.ChecksumIEEE([]byte(StrToMd5(strings.ToLower(str)))))
}

func regMsgInfo(id MID, name string, feilds [][]string, proto, client bool,
	offsets []uintptr, typ reflect.Type) (ret MID) {

	ret = 0
	if id == 0 {
		id = StrToId(name)
	}

	lname := strings.ToLower(name)
	ret = id

	MsgI := minfoByName.Load()[lname]
	if MsgI == nil {
		MsgI = &minfo{}
	}
	MsgI.id = id
	MsgI.name = name
	MsgI.proto = proto
	MsgI.client = client
	MsgI.finfos = nil

	for i, fi := range feilds {
		fname := fi[0]
		tname := strings.ToLower(fi[1])

		isarr := false
		btname := tname
		if tname[0] == 'a' {
			isarr = true
			btname = tname[1:]
		}

		var submi *minfo = nil
		kind := fkindByName[btname]
		if kind == reflect.Invalid {
			submi = minfoByName.Load()[btname]
			if submi == nil {
				Panic("Invalid field type, undefined or non-existent >> FieldName: %s, FieldType: %s", fname, tname)
			}
			kind = reflect.Struct
		}

		offset := uintptr(0)
		if offsets != nil {
			offset = offsets[i]
		}

		fi := &finfo{
			fname:  fname,
			tname:  tname,
			kind:   kind,
			offset: offset,
			submi:  submi,
			isarr:  isarr,
		}
		MsgI.finfos = append(MsgI.finfos, fi)
	}

	newMIInfoById := minfoById.New()
	for k, v := range minfoById.Load() {
		newMIInfoById[k] = v
	}
	newMIInfoById[id] = MsgI

	newMIInfoByName := minfoByName.New()
	for k, v := range minfoByName.Load() {
		newMIInfoByName[k] = v
	}
	newMIInfoByName[lname] = MsgI

	minfoById.Store(newMIInfoById)
	minfoByName.Store(newMIInfoByName)

	log.Println("Register msg info ", name, minfoByName, MsgI.toJson())

	return
}

func RegMsgBy(id MID, name string, feilds [][]string, client bool) MID {
	msgMutex.Lock()
	defer msgMutex.Unlock()

	return regMsgInfo(id, name, feilds, true, client, []uintptr{}, nil)
}

func RegMsg(msg MsgI, client bool) MID {
	msgMutex.Lock()
	defer msgMutex.Unlock()

	return regMsg(msg, true, client)
}

func regMsg(msg MsgI, proto, client bool) MID {
	mtp := reflect.TypeOf(msg)
	mt := mtp.Elem()
	name := mt.Name()

	mi := minfoByType.Load()[name]
	if mi != nil {
		mi.proto = proto
		mi.client = client
		return mi.id
	}

	id := StrToId(name)

	var feilds [][]string
	var offsets []uintptr
	for i := 0; i < mt.NumField(); i++ {
		fs := mt.Field(i)
		ft := fs.Type
		kind := ft.Kind()
		tname := ""
		subname := ""

		if kind == reflect.Slice {
			tname = "a"
			ft = ft.Elem()
			kind = ft.Kind()
		}

		if kind == reflect.Struct {
			if (ft.Name() == reflect.TypeOf(Msg{}).Name()) {
				continue
			}
			if !proto {
				continue
			}
			Panic("The message field type cannot be a non structured pointer >> %s.%s(%s)", name, fs.Name, ft.Name())
		}

		if kind == reflect.Ptr {
			ft = ft.Elem()
			if ft.Kind() != reflect.Struct {
				if !proto {
					continue
				}
				Panic("The message field type cannot be a non structured pointer >> %s.%s(%s)", name, fs.Name, ft.Name())
			}
			kind = reflect.Struct

			subm := reflect.New(ft)
			if !isMsgObj(subm.Type()) {
				if !proto {
					continue
				}
				Panic("The message type must be a struct pointer based on lugo.msg >> %s.%s(%s)", name, fs.Name, ft.Name())
			}

			subname = strings.ToLower(ft.Name())
			regMsg(subm.Interface().(MsgI), proto, client)
		}

		s := ftnameByKind[kind]
		if s == "" {
			Panic("Unsupported message types >> %s.%s(%s)", name, fs.Name, ft.Name())
		}
		if s == "m" {
			s = subname
		}

		tname = tname + s
		feilds = append(feilds, []string{fs.Name, tname})
		offsets = append(offsets, fs.Offset)
	}

	ret := regMsgInfo(id, name, feilds, proto, client, offsets, mt)
	if ret > 0 {
		newMInfoByType := minfoByType.New()
		for k, v := range minfoByType.Load() {
			newMInfoByType[k] = v
		}

		MsgI := minfoByName.Load()[strings.ToLower(name)]
		MsgI.mt = mt
		newMInfoByType[mtp.Name()] = MsgI

		minfoByType.Store(newMInfoByType)

		log.Println(MsgI)
	}

	return ret
}
