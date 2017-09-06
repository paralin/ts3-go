package serverquery

import (
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
		numbers := strings.Split(val, ",")
		decCount := strings.Count(val, ".")
		isFloat := decCount > 0
		var elementType reflect.Type
		if isFloat {
			elementType = reflect.TypeOf(float32(0))
		} else {
			elementType = reflect.TypeOf(int(0))
		}

		parseElement := func(e string) reflect.Value {
			if isFloat {
				i, _ := strconv.ParseFloat(e, 32)
				return reflect.ValueOf(float32(i))
			}

			i, _ := strconv.ParseInt(e, 10, 32)
			return reflect.ValueOf(int(i))
		}

		if len(numbers) == 1 {
			return parseElement(numbers[0]).Interface(), nil
		}

		arr := reflect.MakeSlice(
			reflect.SliceOf(elementType),
			0, 0,
		)
		for _, numStr := range numbers {
			ele := parseElement(numStr)
			arr = reflect.Append(arr, ele)
		}
		return arr.Interface(), nil
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
		outpField := outpVal.Field(i)
		if !unicode.IsUpper(rune(fieldInfo.Name[0])) {
			continue
		}
		if fieldInfo.Anonymous {
			fieldType := fieldInfo.Type
			if fieldType.Kind() != reflect.Ptr {
				if !outpField.CanAddr() {
					continue
				}
				fieldType = reflect.PtrTo(fieldType)
				outpField = outpField.Addr()
			}
			if fieldType.Elem().Kind() != reflect.Struct {
				continue
			}
			unmarshalObject(str, outpField.Interface())
			continue
		}

		sqtag, ok := fieldInfo.Tag.Lookup("serverquery")
		if !ok {
			continue
		}
		argVal, ok := argMap[sqtag]
		if !ok {
			continue
		}
		delete(argMap, sqtag)
		// fmt.Printf("set: %s -> %#v\n", outpField.String(), argVal)
		v := reflect.ValueOf(argVal)
		t := reflect.TypeOf(argVal)
		ot := outpField.Type()

		if ot.Kind() == reflect.Slice && t.Kind() != reflect.Slice {
			sval := reflect.MakeSlice(ot.Elem(), 0, 1)
			sval = reflect.Append(sval, v)
			v = sval
			t = reflect.TypeOf(sval)
		} else if ot.Kind() != reflect.Slice && t.Kind() == reflect.Slice {
			v = v.Index(0)
			t = t.Elem()
		}

		if !t.AssignableTo(ot) {
			if t.ConvertibleTo(ot) {
				v = v.Convert(ot)
			} else {
				if t.Kind() == reflect.Int && ot.Kind() == reflect.Bool {
					v = reflect.ValueOf(v.Int() == 1)
				} else {
					continue
				}
			}
		}
		outpField.Set(v)
	}

	/*
		for unhandled := range argMap {
			fmt.Printf("unhandled argument: %s\n", unhandled)
		}
	*/

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
func UnmarshalArguments(result string, outp interface{}) (interface{}, error) {
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
