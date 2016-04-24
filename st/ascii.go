package st

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/larixsource/suntech/lexer"
)

const (
	Separator  = ';'
	EndOfFrame = '\r'

	tsLayout = "20060102;15:04:05;"
)

var (
	ErrSeparator        = errors.New("invalid separator, a ';' was expected")
	ErrEndOfFrame       = errors.New("invalid end of frame, a CR was expected")
	ErrInvalidDevID     = errors.New("invalid DevID")
	ErrInvalidModel     = errors.New("invalid Model")
	ErrInvalidSwVer     = errors.New("invalid SwVer")
	ErrInvalidDate      = errors.New("invalid Date")
	ErrInvalidTime      = errors.New("invalid Time")
	ErrInvalidCell      = errors.New("invalid Cell")
	ErrInvalidLat       = errors.New("invalid Latitude")
	ErrInvalidLng       = errors.New("invalid Longitude")
	ErrInvalidSpeed     = errors.New("invalid Speed")
	ErrInvalidCourse    = errors.New("invalid Course")
	ErrInvalidSatt      = errors.New("invalid Satt")
	ErrInvalidFix       = errors.New("invalid Fix")
	ErrInvalidDist      = errors.New("invalid Dist")
	ErrInvalidPowerVolt = errors.New("invalid PowerVolt")
	ErrInvalidMode      = errors.New("invalid Mode")
	ErrInvalidMsgNum    = errors.New("invalid MsgNum")
	ErrInvalidHMeter    = errors.New("invalid HMeter")
	ErrInvalidMsgType   = errors.New("invalid MsgType")
)

func AsciiDevID(lex *lexer.Lexer) (devID string, token lexer.Token, err error) {
	token, err = lex.NextFixed(10)
	if err != nil {
		return
	}
	if !token.OnlyDigits() {
		err = ErrInvalidDevID
		return
	}
	if !token.EndsWith(Separator) {
		err = ErrSeparator
		return
	}
	devID = string(token.WithoutSuffix())
	return
}

func AsciiModel(lex *lexer.Lexer) (model Model, token lexer.Token, err error) {
	token, err = lex.NextFixed(3)
	if err != nil {
		return
	}
	if !token.OnlyDigits() {
		err = ErrInvalidModel
		return
	}
	if !token.EndsWith(Separator) {
		err = ErrSeparator
		return
	}
	var md uint64
	md, err = strconv.ParseUint(string(token.WithoutSuffix()), 10, 8)
	if err != nil {
		return
	}
	model = Model(md)
	return
}

func AsciiSwVer(lex *lexer.Lexer) (swVer uint16, token lexer.Token, err error) {
	token, err = lex.NextFixed(4)
	if err != nil {
		return
	}
	if !token.OnlyDigits() {
		err = ErrInvalidSwVer
		return
	}
	if !token.EndsWith(Separator) {
		err = ErrSeparator
		return
	}
	var swv uint64
	swv, err = strconv.ParseUint(string(token.WithoutSuffix()), 10, 16)
	if err != nil {
		return
	}
	swVer = uint16(swv)
	return
}

func AsciiTimestamp(lex *lexer.Lexer) (ts time.Time, tokens []lexer.Token, err error) {
	tokens = make([]lexer.Token, 0, 2)

	// date
	dateToken, dateErr := lex.NextFixed(9)
	tokens = append(tokens, dateToken)
	if dateErr != nil {
		err = dateErr
		return
	}
	if !dateToken.OnlyDigits() {
		err = ErrInvalidDate
		return
	}
	if !dateToken.EndsWith(Separator) {
		err = ErrSeparator
		return
	}

	// time
	timeToken, timeErr := lex.NextFixed(9)
	tokens = append(tokens, timeToken)
	if timeErr != nil {
		err = timeErr
		return
	}
	if timeToken.Type != lexer.DataToken {
		err = ErrInvalidTime
		return
	}
	if !timeToken.EndsWith(Separator) {
		err = ErrSeparator
		return
	}

	// to timestamp
	buf := bytes.NewBuffer(make([]byte, 0, 18))
	buf.Write(dateToken.Literal)
	buf.Write(timeToken.Literal)
	ts, err = time.Parse(tsLayout, buf.String())
	return
}

func AsciiCell(lex *lexer.Lexer) (cell string, token lexer.Token, err error) {
	token, err = lex.Next(7, Separator)
	if err != nil {
		return
	}
	if !token.IsHex() {
		err = ErrInvalidCell
		return
	}
	cell = string(token.WithoutSuffix())
	return
}

func AsciiLat(lex *lexer.Lexer) (lat float32, token lexer.Token, err error) {
	token, err = lex.Next(11, Separator)
	if err != nil {
		return
	}
	if token.Type != lexer.FloatToken {
		err = ErrInvalidLat
		return
	}
	lat64, parseErr := strconv.ParseFloat(string(token.WithoutSuffix()), 32)
	if parseErr != nil {
		err = parseErr
		return
	}
	lat = float32(lat64)
	return
}

func AsciiLon(lex *lexer.Lexer) (lng float32, token lexer.Token, err error) {
	token, err = lex.Next(12, Separator)
	if err != nil {
		return
	}
	if token.Type != lexer.FloatToken {
		err = ErrInvalidLng
		return
	}
	lng64, parseErr := strconv.ParseFloat(string(token.WithoutSuffix()), 32)
	if parseErr != nil {
		err = parseErr
		return
	}
	lng = float32(lng64)
	return
}

func AsciiSpeed(lex *lexer.Lexer) (speed float32, token lexer.Token, err error) {
	token, err = lex.Next(8, Separator)
	if err != nil {
		return
	}
	if token.Type != lexer.FloatToken {
		err = ErrInvalidSpeed
		return
	}
	spd, parseErr := strconv.ParseFloat(string(token.WithoutSuffix()), 32)
	if parseErr != nil {
		err = parseErr
		return
	}
	speed = float32(spd)
	return
}

func AsciiCourse(lex *lexer.Lexer) (speed float32, token lexer.Token, err error) {
	token, err = lex.Next(7, Separator)
	if err != nil {
		return
	}
	if token.Type != lexer.FloatToken {
		err = ErrInvalidCourse
		return
	}
	crs, parseErr := strconv.ParseFloat(string(token.WithoutSuffix()), 32)
	if parseErr != nil {
		err = parseErr
		return
	}
	speed = float32(crs)
	return
}

func AsciiSatellites(lex *lexer.Lexer) (satellites uint8, token lexer.Token, err error) {
	token, err = lex.Next(3, Separator)
	if err != nil {
		return
	}
	if !token.OnlyDigits() {
		err = ErrInvalidSatt
		return
	}
	sat, parseErr := strconv.ParseUint(string(token.WithoutSuffix()), 10, 8)
	if parseErr != nil {
		err = parseErr
		return
	}
	satellites = uint8(sat)
	return
}

func AsciiFix(lex *lexer.Lexer) (fix bool, token lexer.Token, err error) {
	token, err = lex.Next(3, Separator)
	if err != nil {
		return
	}
	if token.Type != lexer.BitsToken {
		err = ErrInvalidFix
		return
	}
	switch token.Literal[0] {
	case '0':
		fix = false
	case '1':
		fix = true
	default:
		err = ErrInvalidFix
	}
	return
}

func AsciiDistance(lex *lexer.Lexer) (distance uint32, token lexer.Token, err error) {
	token, err = lex.Next(11, Separator)
	if err != nil {
		return
	}
	if !token.OnlyDigits() {
		err = ErrInvalidDist
		return
	}
	dist, parseErr := strconv.ParseUint(string(token.WithoutSuffix()), 10, 32)
	if parseErr != nil {
		err = parseErr
		return
	}
	distance = uint32(dist)
	return
}

func AsciiPowerVolt(lex *lexer.Lexer) (powerVolt float32, token lexer.Token, err error) {
	token, err = lex.Next(11, Separator)
	if err != nil {
		return
	}
	if token.Type != lexer.FloatToken {
		err = ErrInvalidPowerVolt
		return
	}
	pv, parseErr := strconv.ParseFloat(string(token.WithoutSuffix()), 32)
	if parseErr != nil {
		err = parseErr
		return
	}
	powerVolt = float32(pv)
	return
}

func AsciiMode(lex *lexer.Lexer) (mode ModeType, token lexer.Token, err error) {
	token, err = lex.NextFixed(2)
	if err != nil {
		return
	}
	if !token.OnlyDigits() {
		err = ErrInvalidMode
		return
	}
	if !token.EndsWith(Separator) {
		err = ErrSeparator
		return
	}
	switch token.Literal[0] {
	case '1':
		mode = IdleMode
	case '2':
		mode = ActiveMode
	case '4':
		mode = DistanceMode
	case '5':
		mode = AngleMode
	default:
		err = fmt.Errorf("invalid mode value: %v", token.Literal[0])
	}
	return
}

func AsciiMsgNum(lex *lexer.Lexer) (msgNum uint16, token lexer.Token, err error) {
	token, err = lex.NextFixed(5)
	if err != nil {
		return
	}
	if !token.IsHex() {
		err = ErrInvalidMsgNum
		return
	}
	if !token.EndsWith(Separator) {
		err = ErrSeparator
		return
	}
	mnum, parseErr := strconv.ParseUint(string(token.WithoutSuffix()), 10, 16)
	if parseErr != nil {
		err = parseErr
	}
	msgNum = uint16(mnum)
	return
}

func AsciiDrivingHourMeter(lex *lexer.Lexer) (hmeter uint32, token lexer.Token, err error) {
	token, err = lex.Next(8, Separator)
	if err != nil {
		return
	}
	if !token.OnlyDigits() {
		err = ErrInvalidHMeter
		return
	}
	hm, parseErr := strconv.ParseUint(string(token.WithoutSuffix()), 10, 32)
	if parseErr != nil {
		err = parseErr
		return
	}
	hmeter = uint32(hm)
	return
}

func AsciiBackupVolt(lex *lexer.Lexer) (backupVolt float32, token lexer.Token, err error) {
	token, err = lex.Next(11, Separator)
	if err != nil {
		return
	}
	if token.Type != lexer.FloatToken {
		err = ErrInvalidHMeter
		return
	}
	bv, parseErr := strconv.ParseFloat(string(token.WithoutSuffix()), 32)
	if parseErr != nil {
		err = parseErr
		return
	}
	backupVolt = float32(bv)

	return
}

func AsciiMsgType(lex *lexer.Lexer, last bool) (realTime bool, token lexer.Token, err error) {
	token, err = lex.NextFixed(2)
	if err != nil {
		return
	}
	if token.Type != lexer.BitsToken {
		err = ErrInvalidMsgType
		return
	}
	switch {
	case last && !token.EndsWith(EndOfFrame):
		err = ErrEndOfFrame
		return
	case !last && !token.EndsWith(Separator):
		err = ErrSeparator
		return
	}
	switch token.Literal[0] {
	case '0':
		realTime = false
	case '1':
		realTime = true
	default:
		err = fmt.Errorf("invalid MsgType value: %v", token.Literal[0])
	}
	return
}

func AsciiUnknownTail(lex *lexer.Lexer, max int) (token lexer.Token, err error) {
	token, err = lex.Next(max, EndOfFrame)
	return
}