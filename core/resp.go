package core

import "errors"

// reads length typically first integer of the string until hit by non-digit byte and returns
// the integer and delta = length + 2 (\r\n)
func readLength(data []byte) (int, int) {
	pos, length := 0, 0
	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}
	return 0, 0
}

// reads RESP encoded simple string from data and returns -> string, delta, error
func readSimpleString(data []byte) (string, int, error) {
	// first character +
	pos := 1

	for ; data[pos] != '\r'; pos++ {
	}
	return string(data[1:pos]), pos + 2, nil
}

// reads RESP encoded errorfrom data and returns -> error string, delta, error
func readError(data []byte) (string, int, error) {
	return readSimpleString(data)
}

// reads RESP encoded integerfrom data and returns -> integer value, delta, error
func readInt64(data []byte) (int64, int, error) {
	// first character :
	pos := 1
	var value int64 = 0

	for ; data[pos] != '\r'; pos++ {
		value = value*10 + int64(data[pos]-'0')
	}
	return value, pos + 2, nil
}

// reads RESP encoded  string from data and returns -> string, delta, error
func readBulkString(data []byte) (string, int, error) {
	// first character $
	pos := 1

	// reading length and forwarding pos by length of integer + first special character
	len, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos:(pos + len)]), pos + len + 2, nil
}

// reads RESP encoded  array from data and returns -> array, delta, error
func readArray(data []byte) (interface{}, int, error) {
	// first character *
	pos := 1

	// reading the length
	count, delta := readLength(data[pos:])
	pos += delta

	var elems []interface{} = make([]interface{}, count)
	for i := range elems {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("no data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '-':
		return readError(data)
	case ':':
		return readInt64(data)
	case '$':
		return readBulkString(data)

	case '*':
		return readArray(data)
	}
	return nil, 0, nil
}

func Decode(data []byte) (interface{}, error) {
	if len(data) == 0 {
		return nil, errors.New("no data")
	}

	value, _, err := DecodeOne(data)
	return value, err
}
