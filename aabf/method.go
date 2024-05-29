package aabf

import (
	//"log"
	"reflect"
	"strings"
	"unsafe"
)

type method func(msg, param unsafe.Pointer)

type eface struct {
	typ  uintptr
	data unsafe.Pointer
}

func isStructPrt(typ reflect.Type) bool {
	return (typ.Kind() == reflect.Ptr) && (typ.Elem().Kind() == reflect.Struct)
}

func isMsgObj(typ reflect.Type) bool {
	_, ret1 := typ.MethodByName("LugoMsg")
	_, ret2 := typ.MethodByName("GetId")
	_, ret3 := typ.MethodByName("SetId")
	_, ret4 := typ.MethodByName("GetH")
	_, ret5 := typ.MethodByName("SetH")
	return ret1 && ret2 && ret3 && ret4 && ret5
}

// servive method info

type sminfo struct {
	methodName string
	msgId      MID
	method     method
	paramCnt   int
	recvQ      *q
}
type sminfos []*sminfo

func regMethod(moduleName string, q *q, methodStruct reflect.Method, methodValue reflect.Value) {
	//log.Println(methodStruct, moduleName)

	methodIface := methodValue.Interface()
	methodType := reflect.TypeOf(methodIface)

	if methodType.Kind() != reflect.Func {
		Panic("registerMethod => 参数methodValue不是方法函数")
	}

	//methodName := runtime.FuncForPC(methodStruct.Func.Pointer()).Name()
	methodName := methodStruct.Name
	if (len(methodName) > 2) && (strings.Index(methodName, "On") != 0) {
		return
	}
	methodName = moduleName + "." + methodName

	pn := methodType.NumIn()
	if (pn == 0) || (pn > 2) {
		return
	}

	if methodType.NumOut() != 0 {
		return
	}

	mt := methodType.In(0)
	if (!isStructPrt(mt)) || (!isMsgObj(mt)) {
		return
	}
	msgName := mt.Elem().Name()

	if pn == 2 {
		pt := methodType.In(1)
		if (pt.Kind() != reflect.Uintptr) && (!isStructPrt(pt)) {
			return
		}
	}

	lmsgName := strings.ToLower(msgName)
	minfo := minfoByName.Load()[lmsgName]
	if minfo == nil {
		msg := reflect.New(mt.Elem()).Interface().(MsgI)
		regMsg(msg, false, false)
		minfo = minfoByName.Load()[lmsgName]
		//Panic("This message has not been registered yet >> " + msgName)
	}

	methodElem := (*eface)(unsafe.Pointer(&methodIface)).data
	method := *(*method)(unsafe.Pointer(&methodElem))

	//log.Println("**************", methodName, msgName, minfo)

	sminfo := &sminfo{
		methodName: methodName,
		msgId:      minfo.id,
		recvQ:      q,
		method:     method,
		paramCnt:   pn,
	}

	minfo.sminfos = append(minfo.sminfos, sminfo)

	return
}

func regMethods(m mi) {
	mt := reflect.TypeOf(m)
	mv := reflect.ValueOf(m)
	mn := mt.Elem().Name()

	q := m.getQ()

	for i := 0; i < mt.NumMethod(); i++ {
		regMethod(mn, q, mt.Method(i), mv.Method(i))
	}
}
