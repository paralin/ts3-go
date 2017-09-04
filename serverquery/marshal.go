package serverquery

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// stringable matches an object that can be String() ed
type stringable interface {
	String() string
}

// encodeArgument encodes a specific argument
func encodeArgument(arg interface{}) (string, error) {
	if argStr, ok := arg.(string); ok {
		if len(argStr) == 0 {
			return "", nil
		}
		return strconv.QuoteToASCII(argStr) + " ", nil
	}

	typeOfArg := reflect.TypeOf(arg)
	valOfArg := reflect.ValueOf(arg)
	kindOfArg := typeOfArg.Kind()
	if kindOfArg == reflect.Ptr {
		if valOfArg.IsNil() {
			return "", nil
		}

		return encodeArgument(valOfArg.Elem().Interface())
	}

	switch kindOfArg {
	case reflect.Struct:
	case reflect.Int:
		return strconv.Itoa(arg.(int)), nil
	default:
		strble, ok := arg.(stringable)
		if !ok {
			return "", errors.Errorf("expected struct or string(able) argument but got a %v", kindOfArg)
		}
		return encodeArgument(strble.String())
	}

	var res bytes.Buffer
	for i := 0; i < typeOfArg.NumField(); i++ {
		fieldInfo := typeOfArg.Field(i)
		sqtag, ok := fieldInfo.Tag.Lookup("serverquery")
		if !ok {
			continue
		}

		res.WriteString(sqtag)
		res.WriteRune('=')
		fieldVal := valOfArg.Field(i)
		fieldStr, err := encodeArgument(fieldVal.Interface())
		if err != nil {
			return "", err
		}

		res.WriteString(fieldStr)
		res.WriteRune(' ') // there will end up with a trailing space, but whatever
	}
	return strings.TrimSpace(res.String()), nil
}

// MarshalArguments converts one or more ServerQuery arguments to a string.
func MarshalArguments(args ...interface{}) (string, error) {
	var buf bytes.Buffer

	for _, arg := range args {
		strResult, err := encodeArgument(arg)
		if err != nil {
			return "", err
		}
		if len(strResult) > 0 {
			buf.WriteString(strResult)
		}
	}

	return buf.String(), nil
}
