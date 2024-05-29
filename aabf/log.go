package aabf

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func exeDir() string {
	return path.Dir(strings.Replace(os.Args[0], "\\", "/", -1))
}

func exeName() string {
	return cutExt(path.Base(strings.Replace(os.Args[0], "\\", "/", -1)))
}

func AbsDir(dir string) string {
	dir = strings.Replace(dir, "\\", "/", -1)
	ret, _ := filepath.Abs(dir)
	return strings.Replace(ret, "\\", "/", -1)
}

func cutExt(s string) string {
	i := len(s) - 1
	for i > 0 {
		if s[i] == '.' {
			return s[:i]
		}
		i--
	}
	return s
}

func hasexeDir(s string) bool {
	s = strings.Replace(s, "\\", "/", -1)
	return strings.Index(strings.ToLower(s), strings.ToLower(exeDir())) >= 0
}

var exeKey = ""
var ExeDir = exeDir()
var ExeName = exeName()
var silenceMode = false
var LogRoot = ExeDir

type file struct {
	log   *logger
	tname string
	fname string
	tflag string
	fhand *os.File
	hour  int
}

func (f *file) open() {
	var err error
	f.fhand, err = os.OpenFile(f.fname, os.O_APPEND|os.O_CREATE|os.O_RDWR,
		os.ModeAppend|os.ModePerm)
	if err != nil {
		fmt.Println(err, f.fname)
		return
	}
}

func (f *file) close() {
	f.fhand.Close()
	f.fhand = nil
}

func (f *file) write(tm time.Time, content string) {
	hour := tm.Hour()
	if f.hour-1 != hour {
		f.hour = hour + 1

		ek := exeKey
		if ek != "" {
			ek = "_" + ek
		}
		if strings.Index(ek, ":[") >= 0 {
			ek = strings.Replace(ek, ":", "", 1)
		} else {
			ek = strings.Replace(ek, ":", "[", 1) + "]"
		}

		if f.fhand != nil {
			f.fhand.Close()
			f.fhand = nil
		}

		path := f.log.logdir + "/" + exeName() + "_" + ek
		os.MkdirAll(path, os.ModePerm)

		f.fname = fmt.Sprintf("%s/%s%s_%s_%02d_%s_%s.log", path, ExeName, ek,
			tm.Format("2006-01-02"), tm.Hour(), f.log.logname, f.tname)
	}

	text := fmt.Sprintf("%s %s [%s] => %s", tm.Format("2006-01-02 15:04:05.000"), ExeName, f.tflag, content)

	if !silenceMode {
		fmt.Println(text)
	}

	if f.fhand == nil {
		f.open()
	}

	_, err := f.fhand.Write([]byte(text + "\r\n"))
	if err != nil {
		f.close()
		fmt.Println(err, f.fname, text)
		return
	}
}

type node struct {
	index int
	time  int64
	text  string
}

type logger struct {
	logdir  string
	logname string
	exeName string
	files   []*file
	c       chan *node
	silence bool
}

func (l *logger) init(logname string) {
	l.logname = logname
	l.c = make(chan *node, 30000)

	l.addFile("Info", "I")
	l.addFile("Waring", "W")
	l.addFile("Error", "E")
	l.addFile("Except", "EXP")

	l.I("Logger \"%s\" is ready", logname)
}

func (l *logger) Go(logdir string) {
	if logdir == "" {
		logdir = ExeDir + "/log"
	} else {
		logdir = strings.Replace(logdir, "\\", "/", -1)
		if !hasexeDir(logdir) {
			if logdir[0] != '/' {
				logdir = "/" + logdir
			}
			logdir = exeDir() + logdir
		}
	}

	l.logdir = AbsDir(logdir)

	go l.exec()
}

func (l *logger) addFile(name, flag string) int {
	f := &file{
		tname: name,
		tflag: flag,
		log:   l,
	}
	l.files = append(l.files, f)
	return len(l.files) - 1
}

func (l *logger) exec() {
	for {
		node := <-l.c
		if node.index < len(l.files) {
			l.files[node.index].write(time.Unix(0, node.time), node.text)
		}

		if len(l.c) != 0 {
			continue
		}

		for _, f := range l.files {
			f.close()
		}
	}
}

func (l *logger) I(format string, args ...any) {
	n := &node{
		index: 0,
		time:  time.Now().UnixNano(),
		text:  fmt.Sprintf(format, args...),
	}
	l.c <- n
}

func (l *logger) W(format string, args ...any) {
	n := &node{
		index: 1,
		time:  time.Now().UnixNano(),
		text:  fmt.Sprintf(format, args...),
	}
	l.c <- n
}

func (l *logger) E(format string, args ...any) {
	n := &node{
		index: 2,
		time:  time.Now().UnixNano(),
		text:  fmt.Sprintf(format, args...),
	}
	l.c <- n
}

func (l *logger) Exp(format string, args ...any) {
	n := &node{
		index: 3,
		time:  time.Now().UnixNano(),
		text:  fmt.Sprintf(format, args...),
	}
	log.Println(l)
	l.c <- n
}

var Log *logger

func SetLogParams(logroot, exekey string, silence bool) {
	LogRoot = logroot
	exeKey = exekey
	silenceMode = silence

	Log.Go(LogRoot)
}

func NewLogger(logname, logdir string) *logger {
	ret := &logger{}
	ret.init(logname)
	ret.Go(logdir)
	return ret
}

func init() {
	Log = &logger{}
	Log.init("SysLog")
}
