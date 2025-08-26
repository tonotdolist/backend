package api

import (
	"reflect"
	"testing"
)

type testInternalStruct struct {
	field int
}

type testVersionedStruct struct {
	fieldd int16
}

type testVersionedResStruct struct {
	fieldd int16
}

func (tis *testVersionedStruct) ToInternalRequest() interface{} {
	return nil
}

func (tis *testVersionedResStruct) FromInternalResponse(_ interface{}) {
}

func TestRequestMapping(t *testing.T) {
	internalTyp := reflect.TypeOf(testInternalStruct{})
	var version uint = 0
	RegisterRequest[testInternalStruct, testVersionedStruct](version)
	req, err := GetRequest(internalTyp, version)
	if err != nil {
		t.Errorf("error getting request: %v", err)
		return
	}

	gotType := reflect.TypeOf(req)
	expectedTyp := reflect.TypeOf((*testVersionedStruct)(nil))
	if gotType != expectedTyp {
		t.Errorf("incorrect GetRequest returned type. expected: %s, got: %s", expectedTyp, gotType)
	}
}

func TestResponseMapping(t *testing.T) {
	var version uint = 0
	RegisterResponse[testInternalStruct, testVersionedResStruct](version)
	req, err := GetResponse(testInternalStruct{}, version)
	if err != nil {
		t.Errorf("error getting response: %v", err)
		return
	}

	gotType := reflect.TypeOf(req)
	expectedTyp := reflect.TypeOf((*testVersionedResStruct)(nil))
	if gotType != expectedTyp {
		t.Errorf("incorrect GetResponse returned type. expected: %s, got: %s", expectedTyp, gotType)
	}
}

var _ VersionedRequest = (*testVersionedStruct)(nil)
var _ VersionedResponse = (*testVersionedResStruct)(nil)
