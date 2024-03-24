package nmea

// FIXME verify all combinations of {Comma,}{Optional,}Type

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	errExpectedComma       = errors.New("expected comma")
	errExpectedDigit       = errors.New("expected digit")
	errExpectedEndOfData   = errors.New("expected end of data")
	errExpectedFloat       = errors.New("expected float")
	errExpectedHexDigit    = errors.New("expected hex digit")
	errExpectedRegexp      = errors.New("expected regexp")
	errUnexpectedByte      = errors.New("unexpected byte")
	errUnexpectedEndOfData = errors.New("unexpected end of data")

	floatRegexp         = regexp.MustCompile(`\A-?\d+(?:\.\d*)?`)
	unsignedFloatRegexp = regexp.MustCompile(`\A\d+(?:\.\d*)?`)
)

type SyntaxError struct {
	Data []byte
	Pos  int
	Err  error
}

func (e *SyntaxError) Error() string {
	// return fmt.Sprintf("%d: %v", e.Pos, e.Err)
	// FIXME make this a single line error? or add Humanize or modify tests show friendly errors?
	return fmt.Sprintf("%s\n%s^: %v", e.Data, strings.Repeat(" ", e.Pos), e.Err)
}

func (e *SyntaxError) Unwrap() error {
	return e.Err
}

type Tokenizer struct {
	data []byte
	pos  int
	err  error
}

func NewTokenizer(data []byte) *Tokenizer {
	return &Tokenizer{
		data: data,
	}
}

func (t *Tokenizer) AtCommaOrEndOfData() bool {
	if t.err != nil {
		return true
	}
	if t.pos == len(t.data) {
		return true
	}
	if t.data[t.pos] == ',' {
		return true
	}
	return false
}

func (t *Tokenizer) AtEndOfData() bool {
	if t.err != nil {
		return true
	}
	return t.pos == len(t.data)
}

func (t *Tokenizer) Bytes() []byte {
	if t.err != nil {
		return nil
	}
	if t.pos == len(t.data) {
		return nil
	}
	start := t.pos
	for t.pos < len(t.data) && t.data[t.pos] != ',' {
		t.pos++
	}
	return t.data[start:t.pos]
}

func (t *Tokenizer) Comma() {
	if t.err != nil {
		return
	}
	if t.pos == len(t.data) {
		t.err = errUnexpectedEndOfData
		return
	}
	if t.data[t.pos] != ',' {
		t.err = errExpectedComma
		return
	}
	t.pos++
}

func (t *Tokenizer) CommaFloat() float64 {
	t.Comma()
	return t.Float()
}

func (t *Tokenizer) CommaHex() int {
	t.Comma()
	return t.Hex()
}

func (t *Tokenizer) CommaHexBytes() []byte {
	t.Comma()
	return t.HexBytes()
}

func (t *Tokenizer) CommaInt() int {
	t.Comma()
	return t.Int()
}

func (t *Tokenizer) CommaLiteralByte(b byte) {
	t.Comma()
	t.LiteralByte(b)
}

func (t *Tokenizer) CommaOneByteOf(bytes string) byte {
	t.Comma()
	return t.OneByteOf(bytes)
}

func (t *Tokenizer) CommaOptionalFloat() Optional[float64] {
	t.Comma()
	return t.OptionalFloat()
}

func (t *Tokenizer) CommaOptionalHex() Optional[int] {
	t.Comma()
	return t.OptionalHex()
}

func (t *Tokenizer) CommaOptionalUnsignedFloat() Optional[float64] {
	t.Comma()
	return t.OptionalUnsignedFloat()
}

func (t *Tokenizer) CommaOptionalUnsignedInt() Optional[int] {
	t.Comma()
	return t.OptionalUnsignedInt()
}

func (t *Tokenizer) CommaString() string {
	t.Comma()
	return t.String()
}

func (t *Tokenizer) CommaUnsignedFloat() float64 {
	t.Comma()
	return t.UnsignedFloat()
}

func (t *Tokenizer) CommaUnsignedInt() int {
	t.Comma()
	return t.UnsignedInt()
}

func (t *Tokenizer) DecimalDigits(n int) int {
	if t.err != nil {
		return 0
	}
	value := 0
	for i := 0; i < n; i++ {
		if t.pos == len(t.data) {
			t.err = errUnexpectedEndOfData
			return 0
		}
		digit, ok := digitValue(t.data[t.pos])
		if !ok {
			t.err = errExpectedDigit
			return 0
		}
		value = 10*value + digit
		t.pos++
	}
	return value
}

func (t *Tokenizer) EndOfData() {
	if t.err != nil {
		return
	}
	if t.pos != len(t.data) {
		t.err = errExpectedEndOfData
		return
	}
}

func (t *Tokenizer) Err() error {
	if t.err == nil {
		return nil
	}
	return &SyntaxError{
		Data: t.data,
		Pos:  t.pos,
		Err:  t.err,
	}
}

func (t *Tokenizer) Float() float64 {
	if t.err != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.err = errUnexpectedEndOfData
		return 0
	}
	m := t.Regexp(floatRegexp)
	if m == nil {
		t.err = errExpectedFloat
		return 0
	}
	value, _ := strconv.ParseFloat(string(m[0]), 64)
	return value
}

func (t *Tokenizer) Fork() *Tokenizer {
	result := *t
	return &result
}

func (t *Tokenizer) Hex() int {
	if t.err != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.err = errUnexpectedEndOfData
		return 0
	}
	value, ok := hexDigitValue(t.data[t.pos])
	if !ok {
		t.err = errExpectedHexDigit
		return 0
	}
	t.pos++
	for t.pos < len(t.data) {
		hexDigit, ok := hexDigitValue(t.data[t.pos])
		if !ok {
			break
		}
		value = 16*value + hexDigit
		t.pos++
	}
	return value
}

func (t *Tokenizer) HexBytes() []byte {
	if t.err != nil {
		return nil
	}
	if t.pos == len(t.data) {
		return nil
	}
	value := []byte{}
	for t.pos+1 < len(t.data) {
		hexDigit1, ok := hexDigitValue(t.data[t.pos])
		if !ok {
			t.err = errExpectedHexDigit
			return nil
		}
		t.pos++
		hexDigit2, ok := hexDigitValue(t.data[t.pos])
		if !ok {
			t.err = errExpectedHexDigit
			return nil
		}
		t.pos++
		byteValue := hexDigit1<<4 + hexDigit2
		value = append(value, byte(byteValue))
	}
	return value
}

func (t *Tokenizer) Int() int {
	if t.err != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.err = errUnexpectedEndOfData
		return 0
	}
	sign := 1
	if t.data[t.pos] == '-' {
		sign = -1
		t.pos++
	}
	return sign * t.UnsignedInt()
}

func (t *Tokenizer) LiteralByte(b byte) {
	if t.err != nil {
		return
	}
	if t.pos == len(t.data) {
		t.err = errUnexpectedEndOfData
		return
	}
	if t.data[t.pos] != b {
		t.err = errUnexpectedByte
		return
	}
	t.pos++
}

func (t *Tokenizer) OneByteOf(bytes string) byte {
	if t.err != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.err = errUnexpectedEndOfData
		return 0
	}
	value := t.data[t.pos]
	if strings.IndexByte(bytes, value) == -1 {
		t.err = errUnexpectedByte
		return 0
	}
	t.pos++
	return value
}

func (t *Tokenizer) OptionalFloat() Optional[float64] {
	if t.err != nil {
		return Optional[float64]{}
	}
	if t.pos == len(t.data) {
		return Optional[float64]{}
	}
	if t.data[t.pos] == ',' {
		return Optional[float64]{}
	}
	return NewOptional(t.Float())
}

func (t *Tokenizer) OptionalHex() Optional[int] {
	if t.err != nil {
		return Optional[int]{}
	}
	if t.pos == len(t.data) {
		return Optional[int]{}
	}
	if t.data[t.pos] == ',' {
		return Optional[int]{}
	}
	return NewOptional(t.Hex())
}

func (t *Tokenizer) OptionalInt() int {
	if t.err != nil {
		return 0
	}
	if t.pos == len(t.data) {
		return 0
	}
	if t.data[t.pos] == ',' {
		return 0
	}
	return t.Int()
}

func (t *Tokenizer) OptionalUnsignedFloat() Optional[float64] {
	if t.err != nil {
		return Optional[float64]{}
	}
	if t.pos == len(t.data) {
		return Optional[float64]{}
	}
	if t.data[t.pos] == ',' {
		return Optional[float64]{}
	}
	return NewOptional(t.UnsignedFloat())
}

func (t *Tokenizer) OptionalUnsignedInt() Optional[int] {
	if t.err != nil {
		return Optional[int]{}
	}
	if t.pos == len(t.data) {
		return Optional[int]{}
	}
	if t.data[t.pos] == ',' {
		return Optional[int]{}
	}
	return NewOptional(t.UnsignedInt())
}

func (t *Tokenizer) Peek() (byte, bool) {
	if t.err != nil {
		return 0, false
	}
	if t.pos == len(t.data) {
		return 0, false
	}
	return t.data[t.pos], true
}

func (t *Tokenizer) Regexp(regexp *regexp.Regexp) [][]byte {
	if t.err != nil {
		return nil
	}
	m := regexp.FindSubmatch(t.data[t.pos:])
	if m == nil {
		t.err = errExpectedRegexp
		return nil
	}
	t.pos += len(m[0])
	return m
}

func (t *Tokenizer) String() string {
	bytes := t.Bytes()
	if bytes == nil {
		return ""
	}
	return string(bytes)
}

func (t *Tokenizer) UnsignedFloat() float64 {
	if t.err != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.err = errUnexpectedEndOfData
		return 0
	}
	m := t.Regexp(unsignedFloatRegexp)
	if m == nil {
		t.err = errExpectedFloat
		return 0
	}
	value, _ := strconv.ParseFloat(string(m[0]), 64)
	return value
}

func (t *Tokenizer) UnsignedInt() int {
	if t.err != nil {
		return 0
	}
	if t.pos == len(t.data) {
		t.err = errUnexpectedEndOfData
		return 0
	}
	if t.data[t.pos] < '0' || '9' < t.data[t.pos] {
		t.err = errExpectedDigit
		return 0
	}
	value := int(t.data[t.pos] - '0')
	t.pos++
	for t.pos < len(t.data) {
		digit, ok := digitValue(t.data[t.pos])
		if !ok {
			break
		}
		value = 10*value + digit
		t.pos++
	}
	return value
}

func digitValue(c byte) (int, bool) {
	if '0' <= c && c <= '9' {
		return int(c - '0'), true
	}
	return 0, false
}

func hexDigitValue(c byte) (int, bool) {
	switch {
	case '0' <= c && c <= '9':
		return int(c - '0'), true
	case 'A' <= c && c <= 'F':
		return int(c - 'A' + 10), true
	case 'a' <= c && c <= 'f':
		return int(c - 'a' + 10), true
	default:
		return 0, false
	}
}