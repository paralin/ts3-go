package serverquery

import (
	"testing"
)

type testArgument struct {
	// ThingType is the type of the thing
	ThingType int `serverquery:"thingtype"`

	// ThingName is the name of the thing.
	ThingName string `serverquery:"thingname"`
}

// TestMarshalArguments tests marshalling arguments.
func TestMarshalArguments(t *testing.T) {
	arg := &testArgument{ThingType: 1, ThingName: "testing 123"}
	str, err := MarshalArguments(arg)
	if err != nil {
		panic(err)
	}

	expected := "thingtype=1 thingname=\"testing 123\""
	if str != expected {
		t.Fatalf("marshal returned (expected %s): %s", expected, str)
	}
}
