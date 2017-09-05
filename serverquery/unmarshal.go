package serverquery

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/kballard/go-shellquote"
	"github.com/pkg/errors"
)

// ParseArgumentValue attempts to determine the type and parses a argument value.
func ParseArgumentValue(val string) (interface{}, error) {
	if len(val) == 0 {
		return nil, errors.New("cannot have empty argument value")
	}

	val = strings.Replace(val, "\\s", " ", -1)
	valRunes := []rune(val)
	firstRune := valRunes[0]
	if firstRune == '"' {
		return string(valRunes[1 : len(valRunes)-1]), nil
	}
	if unicode.IsDigit(firstRune) && !strings.Contains(val, " ") {
		decCount := strings.Count(val, ".")
		switch decCount {
		case 0:
			i, err := strconv.ParseInt(val, 10, 32)
			return int(i), err
		case 1:
			i, err := strconv.ParseFloat(val, 32)
			return float32(i), err
		default:
		}
	}

	return val, nil
}

// ParseArgumentList parses an args string to a map.
func ParseArgumentList(args string) (map[string]interface{}, error) {
	parts, err := shellquote.Split(args)
	if err != nil {
		return nil, err
	}

	res := make(map[string]interface{})
	for _, pt := range parts {
		ptEqParts := strings.SplitN(pt, "=", 2)
		key := ptEqParts[0]
		if len(ptEqParts) != 2 || len(key) < 1 {
			return nil, errors.Errorf("malformed arg: %s", pt)
		}

		val := ptEqParts[1]
		ptVal, err := ParseArgumentValue(val)
		if err != nil {
			return nil, err
		}
		res[key] = ptVal
	}

	return res, nil
}

// unmarshalObject unmarshals an argument list to an object.
func unmarshalObject(str []rune, outp interface{}) error {
	outpType := reflect.TypeOf(outp)
	outpVal := reflect.ValueOf(outp)

	if outpType.Kind() != reflect.Ptr {
		return errors.New("expected to unmarshal an object to a struct pointer")
	}

	outpType = outpType.Elem()
	if outpType.Kind() != reflect.Struct {
		return errors.New("expected to unmarshal an object to a struct pointer")
	}
	outpVal = outpVal.Elem()

	argMap, err := ParseArgumentList(string(str))
	if err != nil {
		return err
	}

	for i := 0; i < outpType.NumField(); i++ {
		fieldInfo := outpType.Field(i)
		sqtag, ok := fieldInfo.Tag.Lookup("serverquery")
		if !ok {
			continue
		}
		argVal, ok := argMap[sqtag]
		if !ok {
			continue
		}
		delete(argMap, sqtag)
		outpField := outpVal.Field(i)
		// fmt.Printf("set: %s -> %#v\n", outpField.String(), argVal)
		outpField.Set(reflect.ValueOf(argVal))
	}

	for unhandled := range argMap {
		fmt.Printf("unhandled argument: %s\n", unhandled)
	}

	return nil
}

// unmarshalArray unmarshals an encoded array into an output array.
func unmarshalArray(str []rune, outp interface{}) (interface{}, error) {
	outpType := reflect.TypeOf(outp)
	if outpType.Kind() != reflect.Slice {
		return nil, errors.New("expected slice output when decoding array")
	}

	elemType := outpType.Elem()
	elemKind := elemType.Kind()
	isPtr := elemKind == reflect.Ptr
	if !isPtr {
		return nil, errors.New("expected to output a slice of struct pointers")
	}

	elemType = elemType.Elem()
	elemKind = elemType.Kind()
	if elemKind != reflect.Struct {
		return nil, errors.New("expected to output a slice of struct pointers")
	}

	outpVal := reflect.ValueOf(outp)
	pts := strings.Split(string(str), "|")
	for _, part := range pts {
		elemVal := reflect.New(elemType)

		err := unmarshalObject([]rune(part), elemVal.Interface())
		if err != nil {
			return nil, err
		}

		outpVal = reflect.Append(outpVal, elemVal)
	}

	return outpVal.Interface(), nil
}

// Unmarshal processes a result into an output interface.
func Unmarshal(result string, outp interface{}) (interface{}, error) {
	outpType := reflect.TypeOf(outp)
	if strings.ContainsRune(result, '|') || outpType.Kind() == reflect.Slice {
		return unmarshalArray([]rune(result), outp)
	}

	isPtr := outpType.Kind() == reflect.Ptr
	if !isPtr {
		return nil, errors.New("unmarshal must be given a pointer")
	}

	err := unmarshalObject([]rune(result), outp)
	if err != nil {
		return nil, err
	}
	return outp, nil
}
