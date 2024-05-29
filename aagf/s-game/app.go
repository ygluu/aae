package aagf

import (
	"flag"
	"libs/aaframe"
)

func AppSGame() {
	aaf.NetSets.SName = "s-game"

	flag.StringVar(&aaf.NetSets.CName, "cname", "AACluster", "集群名字（区服名字）")
	flag.StringVar(&aaf.NetSets.SAddr, "saddr", "127.0.0.1:30000", "服务地址(IP+Port)")

	aaf.SetLogParams(aaf.ExeDir+"/../log/", aaf.NetSets.SAddr, false)

	aaf.Exec()
}
