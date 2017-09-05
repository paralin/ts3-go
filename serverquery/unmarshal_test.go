package serverquery

import (
	"reflect"
	"testing"
)

func TestParseArgumentList(t *testing.T) {
	res, err := ParseArgumentList("test=\"hello world\" type=2")
	if err != nil {
		panic(err)
	}

	val := res["test"]
	if valStr, ok := val.(string); !ok || valStr != "hello world" {
		t.Fatalf("test parsed incorrectly: %#v", val)
	}

	valb := res["type"]
	if valb != 2 {
		t.Fatalf("type parsed incorrectly: %#v (%v) != %#v", valb, reflect.TypeOf(valb).Kind(), 2)
	}
}

func TestParseObjectList(t *testing.T) {
	outp := make([]*testArgument, 0)
	res, err := Unmarshal("thingname=\"hello world\" thingtype=2|thingname=\"goodbye\" thingtype=3", outp)
	if err != nil {
		panic(err)
	}
	outp = res.([]*testArgument)

	if len(outp) != 2 {
		t.Fatal("expected 2 output objects")
	}

	valb := outp[0].ThingType
	if valb != 2 {
		t.Fatalf("type parsed incorrectly: %#v (%v) != %#v", valb, reflect.TypeOf(valb).Kind(), 2)
	}
}

func TestE2E(t *testing.T) {
	arg := &testArgument{ThingType: 1, ThingName: "testing 123"}

	str, err := MarshalArguments(arg)
	if err != nil {
		panic(err)
	}

	argb := &testArgument{}
	outpb, err := Unmarshal(str, argb)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(arg, outpb) {
		t.Fatalf("e2e mismatch: %#v != %#v", arg, outpb)
	}
}
