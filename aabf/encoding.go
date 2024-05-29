package aabf

import (
	"encoding/json"
	//"log"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

func unmarshalFields(mi *minfo, p uintptr, bytes []byte, index, end int) int {
	// log.Println("unmarshalFields", index, end)
	for _, fi := range mi.finfos {
		addr := p + fi.offset
		space := end - index

		if space == 0 {
			return end
		}

		// log.Println("Unmarshal >> ", index, p, fi.offset, space, fi.kind, addr, mi.name, fi.fname)
		if fi.isarr {
			if space < 2 {
				Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
				return 0
			}
			alen := int(uint32(bytes[index]) + uint32(bytes[index+1])<<8)
			index += 2

			switch fi.kind {
			case reflect.Bool, reflect.Int8, reflect.Uint8:
				if space < alen {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				dest := make([]uint8, alen)
				*(*[]uint8)(unsafe.Pointer(addr)) = dest
				for i := 0; i < alen; i++ {
					dest[i] = bytes[index]
					index++
				}
				break
			case reflect.Int16, reflect.Uint16:
				if space < alen*2 {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				dest := make([]uint16, alen)
				*(*[]uint16)(unsafe.Pointer(addr)) = dest
				for i := 0; i < alen; i++ {
					v := uint16(bytes[index])
					index++
					v = v + uint16(bytes[index])<<8
					index++
					dest[i] = v
				}
				break
			case reflect.Int32, reflect.Uint32, reflect.Float32:
				if space < alen*4 {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				dest := make([]uint32, alen)
				*(*[]uint32)(unsafe.Pointer(addr)) = dest
				for i := 0; i < alen; i++ {
					v := uint32(bytes[index])
					index++
					v = v + uint32(bytes[index])<<8
					index++
					v = v + uint32(bytes[index])<<16
					index++
					v = v + uint32(bytes[index])<<24
					index++
					dest[i] = v
				}
				break
			case reflect.Int64, reflect.Uint64, reflect.Float64:
				if space < alen*8 {
					// log.Println(index, space, alen)
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				dest := make([]uint64, alen)
				*(*[]uint64)(unsafe.Pointer(addr)) = dest
				for i := 0; i < alen; i++ {
					v := uint64(bytes[index])
					index++
					v = v + uint64(bytes[index])<<8
					index++
					v = v + uint64(bytes[index])<<16
					index++
					v = v + uint64(bytes[index])<<24
					index++
					v = v + uint64(bytes[index])<<32
					index++
					v = v + uint64(bytes[index])<<40
					index++
					v = v + uint64(bytes[index])<<48
					index++
					v = v + uint64(bytes[index])<<56
					index++
					dest[i] = v
				}
				break
			case reflect.String:
				dest := make([]string, alen)
				*(*[]string)(unsafe.Pointer(addr)) = dest
				for i := 0; i < alen; i++ {
					slen := int(bytes[index])
					index++
					slen = slen + int(uint32(bytes[index])<<8)
					index++
					if slen == 0 {
						continue
					}
					if end-index < slen {
						Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
						return 0
					}

					dest[i] = string(bytes[index : index+slen])
					index += slen
				}
				break
			case reflect.Struct:
				dest := make([]uintptr, alen)
				*(*[]uintptr)(unsafe.Pointer(addr)) = dest
				for i := 0; i < alen; i++ {
					msg, ret := unmarshalMsg(bytes, index, end)
					if ret <= 0 {
						Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
						return 0
					}
					dest[i] = msg.Pointer()
					index = ret
				}
				break
			}
		} else {
			switch fi.kind {
			case reflect.Bool, reflect.Int8, reflect.Uint8:
				if space < 1 {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				*(*byte)(unsafe.Pointer(addr)) = bytes[index]
				index++
				break
			case reflect.Int16, reflect.Uint16:
				if space < 2 {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				*(*uint16)(unsafe.Pointer(addr)) = uint16(bytes[index]) + uint16(bytes[index+1])<<8
				index += 2
				// log.Println("Unmarshal16", index, index)
				break
			case reflect.Int32, reflect.Uint32, reflect.Float32:
				if space < 4 {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				v := uint32(bytes[index])
				index++
				v = v + uint32(bytes[index])<<8
				index++
				v = v + uint32(bytes[index])<<16
				index++
				v = v + uint32(bytes[index])<<24
				index++
				*(*uint32)(unsafe.Pointer(addr)) = v
				// log.Println("Unmarshal32", v, index)
				break
			case reflect.String:
				if space < 2 {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}

				slen := int(bytes[index])
				index++
				slen = slen + int(uint32(bytes[index])<<8)
				index++
				if slen == 0 {
					continue
				}
				// log.Println("Unmarshal string", end-index, slen)
				if end-index < slen {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}

				// log.Println("Unmarshal string", index, slen, string(bytes[index:index+slen]))
				*(*string)(unsafe.Pointer(addr)) = string(bytes[index : index+slen])
				index += slen
				break
			case reflect.Int64, reflect.Uint64, reflect.Float64:
				if space < 8 {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				v := uint64(bytes[index])
				index++
				v = v + uint64(bytes[index])<<8
				index++
				v = v + uint64(bytes[index])<<16
				index++
				v = v + uint64(bytes[index])<<24
				index++
				v = v + uint64(bytes[index])<<32
				index++
				v = v + uint64(bytes[index])<<40
				index++
				v = v + uint64(bytes[index])<<48
				index++
				v = v + uint64(bytes[index])<<56
				index++
				*(*uint64)(unsafe.Pointer(addr)) = v
				// log.Println("Unmarshal64", v, index)
				break
			case reflect.Struct:
				// log.Println("Unmarshal Struct", index, end)
				msg, ret := unmarshalMsg(bytes, index, end)
				if ret <= 0 {
					Log.E("Unmarshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}
				if msg != nil {
					*(*uintptr)(unsafe.Pointer(addr)) = msg.Pointer()
				}
				index = ret
				break
			}
		}
	}

	return end
}

func unmarshalMsg(bytes []byte, index, end int) (*reflect.Value, int) {
	if len(bytes) < 8 {
		return nil, 0
	}

	msgid := MID(bytes[index])
	index++
	msgid = msgid + MID(bytes[index])<<8
	index++
	msgid = msgid + MID(bytes[index])<<16
	index++
	msgid = msgid + MID(bytes[index])<<24
	index++

	// log.Println("Unmarshal MsgId", index-4, msgid)

	mi := minfoById.Load()[msgid]
	if mi == nil {
		return nil, 0
	}

	datalen := uint32(bytes[index])
	index++
	datalen = datalen + uint32(bytes[index])<<8
	index++
	datalen = datalen + uint32(bytes[index])<<16
	index++
	datalen = datalen + uint32(bytes[index])<<24
	index++
	// log.Println("Unmarshal MsgLen", index-4, datalen, end, mi.name)

	ret := 0
	if int(datalen) > end-index {
		return nil, 0
	}
	if datalen == 0 {
		return nil, index
	}

	ins := reflect.New(mi.mt)
	ret = unmarshalFields(mi, ins.Pointer(), bytes, index, index+int(datalen))

	return &ins, ret
}

func Unmarshal(bytes []byte) MsgI {
	msg, ret := unmarshalMsg(bytes, 0, len(bytes))
	if ret <= 0 {
		Log.E("Unmarshal failed, invalid data packet")
		return nil
	}
	return msg.Interface().(MsgI)
}

func strToUtf(s string, d []byte) int {
	ret := 0

	for _, r := range s {
		dest := d[ret : ret+4]
		ret += utf8.EncodeRune(dest, r)
	}

	return ret
}

func marshal(mi *minfo, p uintptr, bytes []byte, capacity, index int) int {
	if capacity < 8 {
		Log.E("Marshal failed >> %s", mi.name)
		return 0
	}
	start := index

	id := mi.id
	bytes[index] = byte(id)
	index++
	bytes[index] = byte(id >> 8)
	index++
	bytes[index] = byte(id >> 16)
	index++
	bytes[index] = byte(id >> 24)
	index++

	// log.Println("Marshal MsgId", id, "Id Index", index-4, mi.name)

	msglenindex := index
	index += 4

	if p != 0 {
		for _, fi := range mi.finfos {
			addr := p + fi.offset
			space := capacity - index
			// log.Println("Marshal >> ", index, p, fi.offset, fi.kind, addr, mi.name, fi.fname)
			if fi.isarr {
				if space < 2 {
					Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
					return 0
				}

				switch fi.kind {
				case reflect.Bool, reflect.Int8, reflect.Uint8:
					arr := *(*[]byte)(unsafe.Pointer(addr))
					alen := len(arr)
					if (space < alen) || (alen > 0xFFFF) {
						Log.E("Marshal failed, insufficient memory space or excessively long array length >> %s.%s", mi.name, fi.fname)
						return 0
					}
					bytes[index] = byte(alen)
					index++
					bytes[index] = byte(uint32(alen) >> 8)
					index++

					copy(bytes[index:index+alen], arr)
					index += alen
					break
				case reflect.Int16, reflect.Uint16:
					arr := *(*[]uint16)(unsafe.Pointer(addr))
					alen := len(arr)
					ablen := alen * 2
					if (space < ablen) || (alen > 0xFFFF) {
						Log.E("Marshal failed, insufficient memory space or array length is too long >> %s.%s", mi.name, fi.fname)
						return 0
					}

					bytes[index] = byte(alen)
					index++
					bytes[index] = byte(uint32(alen) >> 8)
					index++

					for i := 0; i < alen; i++ {
						v := arr[i]
						bytes[index] = byte(v)
						index++
						bytes[index] = byte(v >> 8)
						index++
					}
					break
				case reflect.String:
					arr := *(*[]string)(unsafe.Pointer(addr))
					alen := len(arr)
					if alen > 0xFFFF {
						Log.E("Marshal failed, array length is too long >> %s.%s", mi.name, fi.fname)
						return 0
					}

					bytes[index] = byte(alen)
					index++
					bytes[index] = byte(uint32(alen) >> 8)
					index++

					for _, s := range arr {
						src := []byte(s)
						slen := len(src)
						if space < slen+2 {
							Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
							return 0
						}
						bytes[index] = byte(slen)
						index++
						bytes[index] = byte(uint16(slen) >> 8)
						index++
						copy(bytes[index:index+slen], src)
						index += slen
					}
					break
				case reflect.Int32, reflect.Uint32, reflect.Float32:
					arr := *(*[]uint32)(unsafe.Pointer(addr))
					alen := len(arr)
					ablen := alen * 4
					if (space < ablen) || (alen > 0xFFFF) {
						Log.E("Marshal failed, insufficient memory space or array length is too long >> %s.%s", mi.name, fi.fname)
						return 0
					}

					bytes[index] = byte(alen)
					index++
					bytes[index] = byte(uint32(alen) >> 8)
					index++

					for i := 0; i < alen; i++ {
						v := arr[i]
						bytes[index] = byte(v)
						index++
						bytes[index] = byte(v >> 8)
						index++
						bytes[index] = byte(v >> 16)
						index++
						bytes[index] = byte(v >> 24)
						index++
					}
					break
				case reflect.Int64, reflect.Uint64, reflect.Float64:
					arr := *(*[]uint64)(unsafe.Pointer(addr))
					alen := len(arr)
					ablen := alen * 8
					if (space < ablen) || (alen > 0xFFFF) {
						Log.E("Marshal failed, insufficient memory space or array length is too long >> %s.%s", mi.name, fi.fname)
						return 0
					}

					bytes[index] = byte(alen)
					index++
					bytes[index] = byte(uint32(alen) >> 8)
					index++

					for i := 0; i < alen; i++ {
						v := arr[i]
						bytes[index] = byte(v)
						index++
						bytes[index] = byte(v >> 8)
						index++
						bytes[index] = byte(v >> 16)
						index++
						bytes[index] = byte(v >> 24)
						index++
						bytes[index] = byte(v >> 32)
						index++
						bytes[index] = byte(v >> 40)
						index++
						bytes[index] = byte(v >> 48)
						index++
						bytes[index] = byte(v >> 56)
						index++
					}
					break
				case reflect.Struct:
					arr := *(*[]uintptr)(unsafe.Pointer(addr))
					alen := len(arr)
					if alen > 0xFFFF {
						Log.E("Marshal failed, array length is too long >> %s.%s", mi.name, fi.fname)
						return 0
					}

					bytes[index] = byte(alen)
					index++
					bytes[index] = byte(uint32(alen) >> 8)
					index++

					for _, ins := range arr {
						ret := marshal(fi.submi, ins, bytes, capacity, index)
						if ret <= 0 {
							// log.Println("Marshal failed", capacity-index, ret, index)
							Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
							return 0
						}
						index += ret
					}
					break
				}
			} else {
				switch fi.kind {
				case reflect.Bool, reflect.Int8, reflect.Uint8:
					if space < 1 {
						Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
						return 0
					}
					bytes[index] = *(*byte)(unsafe.Pointer(addr))
					index++
					break
				case reflect.Int16, reflect.Uint16:
					if space < 2 {
						Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
						return 0
					}
					v := *(*uint16)(unsafe.Pointer(addr))
					bytes[index] = byte(v)
					index++
					bytes[index] = byte(v >> 8)
					index++
					break
				case reflect.Int32, reflect.Uint32, reflect.Float32:
					if space < 4 {
						Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
						return 0
					}
					v := *(*uint32)(unsafe.Pointer(addr))
					bytes[index] = byte(v)
					index++
					bytes[index] = byte(v >> 8)
					index++
					bytes[index] = byte(v >> 16)
					index++
					bytes[index] = byte(v >> 24)
					index++
					break
				case reflect.String:
					src := []byte(*(*string)(unsafe.Pointer(addr)))
					slen := len(src)
					if space < slen+2 {
						Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
						return 0
					}
					bytes[index] = byte(slen)
					index++
					bytes[index] = byte(uint16(slen) >> 8)
					index++
					copy(bytes[index:index+slen], src)
					index += slen
					break
				case reflect.Int64, reflect.Uint64, reflect.Float64:
					if space < 8 {
						Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
						return 0
					}
					v := *(*uint64)(unsafe.Pointer(addr))
					bytes[index] = byte(v)
					index++
					bytes[index] = byte(v >> 8)
					index++
					bytes[index] = byte(v >> 16)
					index++
					bytes[index] = byte(v >> 24)
					index++
					bytes[index] = byte(v >> 32)
					index++
					bytes[index] = byte(v >> 40)
					index++
					bytes[index] = byte(v >> 48)
					index++
					bytes[index] = byte(v >> 56)
					index++
					break
				case reflect.Struct:
					// log.Println("struct len", index)
					addr = *(*uintptr)(unsafe.Pointer(addr))
					ret := marshal(fi.submi, addr, bytes, capacity, index)
					if ret <= 0 {
						Log.E("Marshal failed >> %s.%s", mi.name, fi.fname)
						return 0
					}
					// log.Println("struct len", index, ret)
					index += ret
					break
				}
			}
		}
	}

	msglen := index - start - 8

	// log.Println("Marshal MsgLen", index-start, "DataLen", msglen, "LenIndex", msglenindex, index, mi.name)
	bytes[msglenindex] = byte(uint32(msglen))
	msglenindex++
	bytes[msglenindex] = byte(uint32(msglen) >> 8)
	msglenindex++
	bytes[msglenindex] = byte(uint32(msglen) >> 16)
	msglenindex++
	bytes[msglenindex] = byte(uint32(msglen) >> 24)

	return index - start
}

func Marshal(msg MsgI, bytes []byte) []byte {
	mt := reflect.TypeOf(msg)
	minfo := minfoByType.Load()[mt.Name()]
	if minfo == nil {
		Log.E("The message object has not been registered yet >> " + mt.Name())
		return nil
	}

	ins := (uintptr)((*eface)(unsafe.Pointer(&msg)).data)
	ret := marshal(minfo, ins, bytes, len(bytes), 0)
	if ret <= 0 {
		Log.E("Marshal failed due to insufficient memory space")
		return nil
	}

	return bytes[:ret]
}

func ToJson(msg MsgI) string {
	data, err := json.Marshal(msg)
	if err != nil {
		return err.Error()
	}
	return string(data)
}
