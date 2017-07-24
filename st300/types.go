package st300

import (
	"github.com/multiwireless/suntech/st"
)

// MsgType is the type of command or report of a Msg
type MsgType int

const (
	UnknownMsg MsgType = iota
	NTWCmd
	RPTCmd
	EVTCmd
	GSMCmd
	SVCCmd
	MBVCmd
	MSRCmd
	CGFCmd
	ADPCmd
	NPTCmd
	LTMCmd
	PLGCmd
	PLSCmd
	PLCCmd
	CTRCmd
	STRCmd
	GTRCmd

	STTReport
	EMGReport
	EVTReport
	ALTReport
	ALVReport
	UEXReport
	DEXReport

	CMD
)

type Msg struct {
	// Model is the model version. Could be Unknown (some messages don't contain this field)
	Model st.Model

	Type MsgType

	CGF *ST300CGF

	STT *StatusReport
	EMG *EmergencyReport
	EVT *EventReport
	ALT *AlertReport

	Frame []byte

	ParsingError error
}
