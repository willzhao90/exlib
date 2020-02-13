package exutil

import (
	"bytes"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// FloatFloat64 converts a float64 to big float
func FloatFloat64(val float64) (out *big.Float, err error) {
	out = new(big.Float).SetFloat64(val)
	return
}

// FloatInt converts a big int to float
func FloatInt(val *big.Int) (out *big.Float, err error) {
	if err != nil {
		return
	}
	out = new(big.Float).SetInt(val)
	return
}

// FloatIntStr converts a string presentation of a big int to a big float
func FloatIntStr(valInt string) (out *big.Float, err error) {
	i, err := DecodeBigInt(valInt)
	if err != nil {
		return
	}
	out, err = FloatInt(i)
	return
}

// Float converts string representation to a big float
func Float(val string) (out *big.Float, err error) {
	value, ret := new(big.Float).SetString(val)
	if !ret {
		err = fmt.Errorf("invalid va")
		return
	}
	return value, err
}

// Int converts a big.Float to big.Int
func Int(val *big.Float) (out *big.Int) {
	out = new(big.Int)
	val.Int(out)
	return
}

// SubInt x - y
func SubInt(x, y *big.Int) *big.Int {
	return new(big.Int).Sub(x, y)
}

// AddInt x + y
func AddInt(x, y *big.Int) *big.Int {
	return new(big.Int).Add(x, y)
}

// SubFloat x - y
func SubFloat(x, y *big.Float) *big.Float {
	return new(big.Float).Sub(x, y)
}

// AddFloat x + y
func AddFloat(x, y *big.Float) *big.Float {
	return new(big.Float).Add(x, y)
}

// QuoInt x / y
func QuoInt(x, y *big.Int) *big.Int {
	return new(big.Int).Quo(x, y)
}

// QuoFloat x / y
func QuoFloat(x, y *big.Float) *big.Float {
	return new(big.Float).Quo(x, y)
}

// MulInt x * y
func MulInt(x, y *big.Float) *big.Float {
	return new(big.Float).Mul(x, y)
}

// MulFloat x * y
func MulFloat(x, y *big.Float) *big.Float {
	return new(big.Float).Mul(x, y)
}

// MinUint64 returns the smaller value of a and b
func MinUint64(a, b uint64) uint64 {
	if a < b {
		return a
	}
	return b
}

// MaxUint64 returns the larger value of a and b
func MaxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

// DecodeBigInt ascii to big.Int
func DecodeBigInt(txt string) (*big.Int, error) {
	if txt == "" {
		return new(big.Int), nil // Defaults to 0
	}
	res, success := new(big.Int).SetString(txt, 10)
	if !success {
		return nil, fmt.Errorf("cannot decode %v into big.Int", txt)
	}
	return res, nil
}

// DecodeTwoBigInts at the same time
func DecodeTwoBigInts(a, b string) (*big.Int, *big.Int, error) {
	aInt, err := DecodeBigInt(a)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot decode first string %v: %v", a, err)
	}
	bInt, err := DecodeBigInt(b)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot decode second string %v: %v", b, err)
	}
	return aInt, bInt, err
}

// AddBigInt adds two string represented big.Int
func AddBigInt(a, b string) (*big.Int, error) {
	aInt, bInt, err := DecodeTwoBigInts(a, b)
	if err != nil {
		return nil, err
	}

	return aInt.Add(aInt, bInt), nil
}

// BigPow10 give 10**n
func BigPow10(n int) *big.Int {
	return new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(n)), new(big.Int))
}

const (
	// MaxPreciseDigitNum give the max precise digit number
	MaxPreciseDigitNum = 20
)

// ImprFloat define an imprecised float type between string_decimal and big float
type ImprFloat struct {
	data    []byte //original data
	err     error
	pos     int
	endFlag bool
}

// FromString set itself with string
func (imf *ImprFloat) FromString(valStr string) *ImprFloat {
	imf.resetFlag()
	if !ValidPositiveFloatStr(valStr) {
		err := fmt.Errorf("ImprFloat warning: Bad float string %v, use default value", valStr)
		log.Println(err)
		imf.err = err
		valStr = "0"
	}
	valStr = strings.TrimSpace(valStr)
	imf.data = []byte(valStr)
	//fill '.'
	size := len(imf.data)
	pos := bytes.Index(imf.data, []byte("."))
	if pos == -1 { // no '.' so assume '.' in the end
		imf.data = append(imf.data, '.')
		pos = size
	}
	imf.pos = pos
	return imf
}

// Shift data * 10**n
func (imf *ImprFloat) Shift(n int) *ImprFloat {
	imf.resetFlag()
	if n > 200 || n < -200 {
		err := fmt.Errorf("ImprFloat warning: Shift %v out of range, ignore shift", n)
		log.Println(err)
		imf.err = err
		return imf
	}
	if n == 0 {
		return imf
	}

	size := len(imf.data)
	pos := imf.pos
	if pos == -1 {
		err := fmt.Errorf("ImprFloat warning: Intend to shift %v before assigning data, ignore shift", n)
		log.Println(err)
		imf.err = err
		return imf
	}
	//fill
	aim := pos + n
	var fillSize int
	if n > 0 {
		fillSize = aim - size + 1
		if fillSize <= 0 {
			fillSize = 0
		} else {
			imf.data = append(imf.data, strings.Repeat("0", fillSize)...)
		}

		for i := pos; i < aim; i++ {
			imf.data[i], imf.data[i+1] = imf.data[i+1], imf.data[i]
		}
	} else {
		fillSize = -aim + 1
		if fillSize <= 0 {
			fillSize = 0
		} else {
			temp := imf.data
			imf.data = []byte(strings.Repeat("0", fillSize))
			imf.data = append(imf.data, temp...)
			pos += fillSize
			aim = 1
		}
		for i := pos; i > aim; i-- {
			imf.data[i], imf.data[i-1] = imf.data[i-1], imf.data[i]
		}
	}
	pos = aim
	imf.pos = pos
	return imf
}

// ToString convert itself to string in float format
func (imf *ImprFloat) ToString() string {
	defer imf.endOp()
	str := strings.TrimRight(string(imf.data), " 0")
	return strings.TrimRight(str, ".")
}

// TrimDecimal remove all decimal from data
func (imf *ImprFloat) TrimDecimal() *ImprFloat {
	imf.resetFlag()
	imf.data = imf.data[0 : imf.pos+1]
	return imf
}

// ToIntString convert itself to string in float format
func (imf *ImprFloat) ToIntString() string {
	defer imf.endOp()
	pos := imf.pos
	str := string(imf.data[:pos])
	str = strings.TrimLeft(str, "0")
	if len(str) == 0 {
		return "0"
	}
	return str
}

// FromFloat set itself with float
func (imf *ImprFloat) FromFloat(v float64) *ImprFloat {
	imf.resetFlag()
	fstr := new(big.Float).SetFloat64(v).Text('f', MaxPreciseDigitNum)
	imf = imf.FromString(fstr)
	return imf
}

// ToFloat convert itself to float64
func (imf *ImprFloat) ToFloat() float64 {
	defer imf.endOp()
	bigf, ok := new(big.Float).SetString(imf.ToString())
	if !ok {
		err := fmt.Errorf("ToFloat: error float format of: %v", imf.ToString())
		log.Println(err)
		imf.err = err
		return 0.0
	}
	f, _ := bigf.Float64()
	return f
}

// GetLastErr return error in operation
func (imf *ImprFloat) GetLastErr() error {
	err := imf.err
	imf.err = nil
	return err
}

// when ImprFloat output, call endOp
func (imf *ImprFloat) endOp() {
	imf.endFlag = true
}

func (imf *ImprFloat) resetFlag() {
	if imf.endFlag == true {
		imf.err = nil
		imf.endFlag = false
	}
}

// ValidFloatStr check a float string
func ValidFloatStr(val string) bool {
	var validFloat = regexp.MustCompile(`^[-+]?([0-9]+(\.[0-9]+)?)$`)
	return validFloat.MatchString(val)
}

// ValidPositiveFloatStr check a positive float string
func ValidPositiveFloatStr(val string) bool {
	var validFloat = regexp.MustCompile(`^[+]?([0-9]+(\.[0-9]+)?)$`)
	return validFloat.MatchString(val)
}
