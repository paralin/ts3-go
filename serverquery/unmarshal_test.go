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
	outp := make([]*TestArgument, 0)
	res, err := UnmarshalArguments("thingname=\"hello world\" thingtype=2 nested=\"nested thing\"|thingname=\"goodbye\" thingtype=3", outp)
	if err != nil {
		panic(err)
	}
	outp = res.([]*TestArgument)

	if len(outp) != 2 {
		t.Fatal("expected 2 output objects")
	}

	valb := outp[0].ThingType
	if valb != 2 {
		t.Fatalf("type parsed incorrectly: %#v (%v) != %#v", valb, reflect.TypeOf(valb).Kind(), 2)
	}
	valc := outp[0].Nested
	if valc != "nested thing" {
		t.Fatalf("type parsed incorrectly: %#v (%v) != %#v", valc, reflect.TypeOf(valc).Kind(), "nested thing")
	}
}

func TestE2E(t *testing.T) {
	arg := &TestArgument{
		ThingType: 1,
		ThingName: "testing 123",
		NestedTestArgument: NestedTestArgument{
			Nested: "nested thing",
		},
		ThingBool: true,
		ThingList: []int{
			0, 1, 2, 3,
		},
	}

	str, err := MarshalArguments(arg)
	if err != nil {
		panic(err)
	}
	t.Log(str)

	argb := &TestArgument{}
	outpb, err := UnmarshalArguments(str, argb)
	if err != nil {
		panic(err)
	}

	if !reflect.DeepEqual(arg, outpb) {
		t.Fatalf("e2e mismatch: %#v != %#v", arg, outpb)
	}
}
