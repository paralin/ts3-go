package serverquery

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

// stringable matches an object that can be String() ed
type stringable interface {
	String() string
}

// encodeArgument encodes a specific argument
func encodeArgument(arg interface{}) (rStr string, rErr error) {
	defer func() {
		rStr = strings.TrimSpace(rStr)
	}()
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
	case reflect.Slice:
		var buf bytes.Buffer
		l := valOfArg.Len()
		for i := 0; i < l; i++ {
			str, err := encodeArgument(valOfArg.Index(i).Interface())
			if err != nil {
				return "", err
			}
			str = strings.TrimSpace(str)
			if i != 0 {
				buf.WriteRune(',')
			}
			buf.WriteString(str)
		}
		return buf.String(), nil
	case reflect.Struct:
	case reflect.Int:
		return strconv.Itoa(arg.(int)), nil
	case reflect.Bool:
		if arg.(bool) {
			return "1", nil
		} else {
			return "0", nil
		}
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
		fieldVal := valOfArg.Field(i)
		if !unicode.IsUpper(rune(fieldInfo.Name[0])) {
			continue
		}
		if fieldInfo.Anonymous {
			fieldType := fieldInfo.Type
			if !unicode.IsUpper(rune(fieldInfo.Name[0])) {
				continue
			}
			if fieldType.Kind() == reflect.Ptr {
				if fieldVal.IsNil() {
					continue
				}
				fieldType = fieldType.Elem()
				fieldVal = fieldVal.Elem()
			}
			str, err := encodeArgument(fieldVal.Interface())
			if err != nil {
				return "", err
			}
			res.WriteString(str)
			res.WriteRune(' ')
			continue
		}

		sqtag, ok := fieldInfo.Tag.Lookup("serverquery")
		if !ok {
			continue
		}

		res.WriteString(sqtag)
		res.WriteRune('=')
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
