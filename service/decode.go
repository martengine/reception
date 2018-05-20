package service

import (
	"errors"
	"io"
	"strconv"
)

type decoder struct {
	index   int
	dataLen int
	data    []byte
	result  map[string][]byte
}

func newDecoder(data []byte) (*decoder, error) {
	dataLen := len(data)

	if dataLen < 2 {
		return nil, errors.New("invalid data")
	}

	if data[0] != '{' || data[dataLen-1] != '}' {
		return nil, errors.New("invalid data")
	}

	data = data[1 : dataLen-1]

	return &decoder{data: data, dataLen: len(data), result: map[string][]byte{}}, nil
}

func (d *decoder) decode() error {
	var key string
	var value []byte

	var err error

	// walk through the whole data.
	for {
		// find the key first.
		key, err = d.findKey()
		if err != nil {
			if err == io.EOF {
				// in case if there was holistic key.
				d.result[key] = nil
				return nil
			}

			return err
		}

		// find the value first.
		if value, err = d.findValue(); err != nil {
			if err == io.EOF {
				// in case if there was holistic key.
				d.result[key] = nil
				return nil
			}

			return err
		}

		// save it to results.
		d.result[key] = value
	}
}

func (d *decoder) findKey() (string, error) {
	var keyStarted bool
	var startIndex int

	for {
		// check for end of data.
		if d.index >= d.dataLen {
			if keyStarted {
				return "", io.ErrUnexpectedEOF
			}

			return "", io.EOF
		}

		c := d.data[d.index]

		if !keyStarted {
			// ignore the whitespaces.
			if c == ' ' || c == '\n' || c == '\t' {
				d.index++
				continue
			}

			// key starts with doublequote.
			if c != '"' {
				return "", errors.New("unexpected '" + string(c) + "' at index " + strconv.Itoa(d.index))
			}

			// key starts right here.
			keyStarted = true
			startIndex = d.index
			d.index++
			continue
		}

		// not the end of key.
		if c != '"' {
			d.index++
			continue
		}

		// check if key is empty.
		if startIndex+1 == d.index {
			return "", errors.New("empty key at " + strconv.Itoa(startIndex))
		}

		result := d.data[startIndex+1 : d.index]
		d.index++

		return string(result), nil
	}
}

func (d *decoder) findValue() ([]byte, error) {
	var valueStarted bool
	var aboutToStart bool
	var startIndex int
	var expect byte
	var result []byte

	for {
		// check for end of data.
		if d.index >= d.dataLen {
			if valueStarted || aboutToStart {
				return nil, io.ErrUnexpectedEOF
			}

			return nil, io.EOF
		}

		c := d.data[d.index]

		// when we found result it's good to skip all the valid characters until the
		// new key starts.
		if result != nil {
			if c == ' ' || c == '\n' || c == '\t' || c == ',' {
				d.index++
				continue
			}

			return result, nil
		}

		if !valueStarted {
			if c == ' ' || c == '\n' || c == '\t' {
				d.index++
				continue
			}

			if !aboutToStart {
				switch c {
				case ',', '}':
					// there would be no value.
					d.index++
					return nil, nil
				case ':':
					aboutToStart = true
					d.index++
					continue
				}
			}

			// the symbol that marks the beggining of value
			// also has specific symbol that would mark the end of the value.
			switch c {
			case '"':
				expect = '"'
			case '{':
				expect = '}'
			case '[':
				expect = ']'
			default:
				return nil, errors.New("unexpected '" + string(c) + "' at index " + strconv.Itoa(d.index))
			}

			valueStarted = true
			startIndex = d.index
			d.index++
			continue
		}

		// value still going.
		if c != expect {
			d.index++
			continue
		}

		result = d.data[startIndex+1 : d.index]
		d.index++
	}
}

// Decode the service data to map with according references.
func Decode(data []byte) (map[string][]byte, error) {
	dataDecoder, err := newDecoder(data)
	if err != nil {
		return nil, err
	}

	if err := dataDecoder.decode(); err != nil {
		return nil, err
	}

	return dataDecoder.result, nil
}
