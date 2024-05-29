package aagf

import (
	"flag"
	"libs/aabf"
)

func AppSLogin() {
	aabf.NetSets.SName = "s-login"

	flag.StringVar(&aabf.NetSets.CName, "cname", "gecluster", "集群名字（区服名字）")
	flag.StringVar(&aabf.NetSets.SAddr, "saddr", "127.0.0.1:20000", "服务地址(IP+Port)")

	aabf.SetLogParams(aabf.ExeDir+"/../log/", aabf.NetSets.SAddr, false)

	aabf.Exec()
}
