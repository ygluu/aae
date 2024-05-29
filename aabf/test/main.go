package main

import (
	"libs/afgo"
)

func main() {
	af.SetLogParams(af.ExeDir+"/log/", "123.12.12.00:1234", false)

	af.Exec()
}
