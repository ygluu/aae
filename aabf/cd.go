package aabf

type CDSpeeding interface {
	OnSpeeding(ConnInfo, ConnAddr string)
}

/* CDSpeeding demo

type cd struct {
}

func (cd *cd) OnSpeeding(ConnInfo, ConnAddr string) {

}

*/

func SetCDTime(MsgName string, Time uint, CDS CDSpeeding) {

}
