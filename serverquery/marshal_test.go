package serverquery

import (
	"testing"
)

type NestedTestArgument struct {
	// Nested is a nested string
	Nested string `serverquery:"nested"`
}

type TestArgument struct {
	NestedTestArgument

	// ThingType is the type of the thing
	ThingType int `serverquery:"thingtype"`
	// ThingName is the name of the thing.
	ThingName string `serverquery:"thingname"`
	// ThingBool is a boolean field
	ThingBool bool `serverquery:"thing_bool"`
	// ThingList is a list of things.
	ThingList []int `serverquery:"thing_list"`
}

// TestMarshalArguments tests marshalling arguments.
func TestMarshalArguments(t *testing.T) {
	arg := &TestArgument{ThingBool: true, ThingType: 1, ThingName: "testing 123", NestedTestArgument: NestedTestArgument{
		Nested: "nested thing",
	}, ThingList: []int{5, 4, 2}}
	str, err := MarshalArguments(arg)
	if err != nil {
		panic(err)
	}

	expected := "nested=\"nested thing\" thingtype=1 thingname=\"testing 123\" thing_bool=1 thing_list=5,4,2"
	if str != expected {
		t.Fatalf("marshal returned (expected %s): %s", expected, str)
	}
}

// TestMarshalCommand tries to marshal a command.
func TestMarshalCommand(t *testing.T) {
	cmd := &GetClientInfoCommand{ClientId: 1}
	str, err := MarshalCommand(cmd)
	if err != nil {
		t.Fatal(err.Error())
	}
	expected := "clientinfo clid=1"
	if str != expected {
		t.Fatalf("expected %s, got: %s", expected, str)
	}
}
