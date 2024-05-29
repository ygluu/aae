package aagf

import (
	"flag"
	"libs/aabf"
)

func AppSGate() {
	flag.StringVar(&aabf.NetSets.CName, "cname", "AACluster", "集群名字（区服名字）")
	flag.StringVar(&aabf.NetSets.SAddr, "saddr", "127.0.0.1:10000", "服务地址(IP+Port)")
	flag.BoolVar(&aabf.NetSets.AllowAutoConnect, "aaconn", false, "是否运行其他服务主动连接")

	aabf.NetSets.SName = "s-gate"
	aabf.NetSets.IsGate = true
	aabf.SetLogParams(aabf.ExeDir+"/../log/", aabf.NetSets.SAddr, false)

	aabf.Exec()
}
