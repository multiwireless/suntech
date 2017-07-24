package st600

import (
	"github.com/multiwireless/suntech/lexer"
	"github.com/multiwireless/suntech/st"
)

type AliveReport struct {
	Hdr   MsgType
	DevID string
}

func parseALVAscii(lex *lexer.Lexer, msg *Msg) {
	msg.Type = ALVReport

	alv := &AliveReport{
		Hdr: ALVReport,
	}
	msg.ALV = alv

	devID, token, err := st.AsciiDevIDAtEnd(lex)
	msg.Frame = append(msg.Frame, token.Literal...)
	if err != nil {
		msg.ParsingError = err
		return
	}
	alv.DevID = devID

	return
}
