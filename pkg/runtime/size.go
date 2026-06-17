package runtime

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"strconv"
)

// Size represents any amount of memory occupied by something, like 4 bytes or
// 12 Megabytes (12 MB), useful to handle amounts of memory like you handle
// time durations.
type Size uint64

// Size units
const (
	Byte     Size = 1
	Kilobyte      = 1024 * Byte
	Megabyte      = 1024 * Kilobyte
	Gigabyte      = 1024 * Megabyte
	Terabyte      = 1024 * Gigabyte
)

var units = []string{
	"B",
	"KB",
	"MB",
	"GB",
	"TB",
}

// String returns the size representation as the number of bytes represented by
// the instance with the bigger unit that value can be represent.
//
// For any value below 1024 Bytes, just the number will be represented.
func (s Size) String() string {
	return string(s.fmt(2, 0))
}

// Format implements fmt.Formatter interface.
//
// This ensures that a Size can be printed using %v, %+v, %#v, %q, %d, %s and %f
func (s Size) Format(f fmt.State, c rune) {
	width, ok := f.Width()
	if !ok {
		width = -1
	}

	prec, ok := f.Precision()
	if !ok {
		prec = 0
	}

	switch c {
	case 'v':
		var buf []byte
		switch {
		case f.Flag('#'):
			buf = []byte("runtime.Size(" + strconv.FormatUint(uint64(s), 10) + ")")
		case f.Flag('+'):
			buf = []byte(strconv.FormatUint(uint64(s), 10))
		default:
			buf = s.fmt(prec, width)
		}
		f.Write(buf)
	case 'q':
		f.Write([]byte(`"`))
		f.Write(s.fmt(prec, width))
		f.Write([]byte(`"`))
	case 'd', 's':
		f.Write(s.fmt(prec, width))
	case 'f':
		f.Write(s.fmt(prec, width))
	}
}

func (s Size) fmt(precision, width int) []byte {
	var (
		value float64
		unit  string
	)
	switch {
	case s == 0:
		value = 0
		unit = ""
	case s >= Terabyte:
		value = float64(s) / float64(Terabyte)
		unit = "TB"
	default:
		i, f := math.Modf(float64(s))
		base := math.Log(i) / math.Log(1024)
		value = math.Pow(1024, base-math.Floor(base))
		prec := int(math.Pow(10, float64(precision)))
		digit := float64(prec) * value
		var round float64
		_, div := math.Modf(digit)
		if div >= .5 {
			round = math.Ceil(digit)
		} else {
			round = math.Floor(digit)
		}
		value = round / float64(prec)
		value += f
		unit = units[int(math.Floor(base))]
	}

	res := []byte(strconv.FormatFloat(value, 'f', precision, 64) + unit)
	if len(res) < width {
		res = append(bytes.Repeat([]byte(" "), width-len(res)), res...)
	}
	return res
}

// ParseSize transforms a size string into a Size.
//
// A size string is a sequence of decimal numbers with an optional unit,
// possible units are: "" (bytes), "KB", "MB", "GB", "TB".
//
// It can be seen as the opposite to String, hence the next holds:
//
//	ParseSize(size.String()) == size
func ParseSize(str string) (Size, error) {
	// [0-9]*(\.[0-9]*)?[A-Z]+
	orig := str

	if str == "0" {
		return 0, nil
	}
	if str == "" {
		return 0, errors.New("runtime: invalid size " + orig)
	}

	// The next character must be [0-9.]
	if !(str[0] == '.' || '0' <= str[0] && str[0] <= '9') {
		return 0, errors.New("runtime: invalid size " + orig)
	}

	var (
		v, f  uint64 // integers before, after decimal point
		err   error
		s     uint64
		scale float64 = 1 // value = v + f/scale
	)

	// Consume [0-9]*
	pl := len(str)
	v, str, err = leadingInt(str)
	if err != nil {
		return 0, errors.New("runtime: invalid size " + orig)
	}
	pre := pl != len(str) // whether we consumed anything before a period

	// Consume (\.[0-9]*)?
	post := false
	if str != "" && str[0] == '.' {
		str = str[1:]
		pl := len(str)
		f, scale, str = leadingFraction(str)
		post = pl != len(str)
	}
	if !pre && !post {
		// no digits (e.g. ".s" or "-.s")
		return 0, errors.New("runtime: invalid size " + orig)
	}

	// Consume unit. (optional)
	i := 0
	for ; i < len(str); i++ {
		c := str[i]
		if c == '.' || '0' <= c && c <= '9' {
			break
		}
	}

	u := str[:i]
	unit, ok := unitMap[u]
	if !ok {
		return 0, errors.New("runtime: unknown unit " + u + " in size " + orig)
	}
	if v > (1<<63-1)/unit {
		// overflow
		return 0, errors.New("runtime: invalid size " + orig)
	}
	v *= unit
	if f > 0 {
		v += uint64(float64(f) * (float64(unit) / scale))
		if v < 0 {
			// overflow
			return 0, errors.New("runtime: invalid size " + orig)
		}
	}
	s += v
	if s < 0 {
		// overflow
		return 0, errors.New("runtime: invalid size " + orig)
	}

	return Size(s), nil
}

// leadingInt consumes the leading [0-9]* from s.
func leadingInt(s string) (x uint64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > (1<<63-1)/10 {
			// overflow
			return 0, "", errors.New("runtime: bad [0-9]*")
		}
		x = x*10 + uint64(c) - '0'
		if x < 0 {
			// overflow
			return 0, "", errors.New("runtime: bad [0-9]*")
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x uint64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + uint64(c) - '0'
		if y < 0 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}

var unitMap = map[string]uint64{
	"B":  uint64(Byte),
	"KB": uint64(Kilobyte),
	"MB": uint64(Megabyte),
	"GB": uint64(Gigabyte),
	"TB": uint64(Terabyte),
}
