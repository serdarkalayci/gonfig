package gonfig

import (
	"errors"
	"fmt"
	"strconv"
)

func convertToInt(val interface{}) (int, error) {
	var i int // your final value
	var err error
	switch t := val.(type) {
	case int:
		i = t
	case int8:
		i = int(t) // standardizes across systems
	case int16:
		i = int(t) // standardizes across systems
	case int32:
		i = int(t) // standardizes across systems
	case int64:
		i = int(t) // standardizes across systems
	case bool:
		if t {
			i = 1
		} else {
			i = 0
		}
	case float32:
		i = int(t) // standardizes across systems
	case float64:
		i = int(t) // standardizes across systems
	case uint:
		i = int(t) // standardizes across systems
	case uint8:
		i = int(t) // standardizes across systems
	case uint16:
		i = int(t) // standardizes across systems
	case uint32:
		i = int(t) // standardizes across systems
	case uint64:
		i = int(t) // standardizes across systems
	case string:
		i, err = strconv.Atoi(t)
	default:
		i = 0
		err = errors.New("Unknown type")
	}
	return i, err
}

func convertToString(val interface{}) string {
	return fmt.Sprintf("%v", val)
}

func convertToFloat(val interface{}) (float64, error) {
	var f float64 // your final value
	var err error
	switch t := val.(type) {
	case int:
		f = float64(t)
	case int8:
		f = float64(t) // standardizes across systems
	case int16:
		f = float64(t) // standardizes across systems
	case int32:
		f = float64(t) // standardizes across systems
	case int64:
		f = float64(t) // standardizes across systems
	case bool:
		if t {
			f = 1
		} else {
			f = 0
		}
	case float32:
		f = float64(t) // standardizes across systems
	case float64:
		f = t // standardizes across systems
	case uint:
		f = float64(t) // standardizes across systems
	case uint8:
		f = float64(t) // standardizes across systems
	case uint16:
		f = float64(t) // standardizes across systems
	case uint32:
		f = float64(t) // standardizes across systems
	case uint64:
		f = float64(t) // standardizes across systems
	case string:
		f, err = strconv.ParseFloat(t, 64)
	default:
		f = 0
		err = errors.New("Unknown type")
	}
	return f, err
}

func convertToBool(val interface{}) (bool, error) {
	var b bool = false // your final value
	var err error
	switch t := val.(type) {
	case int:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case int8:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case int16:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case int32:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case int64:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case bool:
		b = t
	case float32:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case float64:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case uint:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case uint8:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case uint16:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case uint32:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case uint64:
		switch t {
		case 0:
			b = false
		case 1:
			b = true
		default:
			err = errors.New("Unknown value")
		}
	case string:
		switch t {
		case "0":
			b = false
		case "1":
			b = true
		default:
			err = errors.New("Unknown value")
		}
	default:
		b = false
		err = errors.New("Unknown type")
	}
	return b, err
}
