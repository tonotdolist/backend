package api

import (
	"fmt"
	"reflect"
)

var requestToInternal = make(map[uint]map[reflect.Type]reflect.Type)
var internalToResponse = make(map[uint]map[reflect.Type]reflect.Type)

type VersionedRequest interface {
	ToInternalRequest() interface{}
}

type VersionedResponse interface {
	FromInternalResponse(data interface{}) interface{}
}

func RegisterRequest[TInternal any, TVersion any](apiVersion uint) {
	requestMapping := requestToInternal[apiVersion]

	if requestMapping == nil {
		requestMapping = make(map[reflect.Type]reflect.Type)
		requestToInternal[apiVersion] = requestMapping
	}

	internalType := reflect.TypeOf((*TInternal)(nil))
	versionedType := reflect.TypeOf((*TVersion)(nil))

	_, ok := versionedType.(VersionedRequest)
	if !ok {
		panic(fmt.Sprintf("%T does not implement VersionedRequest interface", versionedType))
	}

	requestToInternal[apiVersion][internalType] = versionedType
}

func RegisterResponse[TInternal any, TVersion any](apiVersion uint) {
	responseMapping := internalToResponse[apiVersion]

	if responseMapping == nil {
		responseMapping = make(map[reflect.Type]reflect.Type)
		internalToResponse[apiVersion] = responseMapping
	}

	internalType := reflect.TypeOf((*TInternal)(nil))
	versionedType := reflect.TypeOf((*TVersion)(nil))
	_, ok := versionedType.(VersionedResponse)
	if !ok {
		panic(fmt.Sprintf("%T does not implement VersionedResponse interface", versionedType))
	}

	internalToResponse[apiVersion][versionedType] = internalType
}

func GetRequest(internalType reflect.Type, version uint) (VersionedRequest, error) {
	requestMapping, ok := requestToInternal[version]
	if !ok {
		return nil, fmt.Errorf("api version %d not registered on the versioned request to internal request mapping", version)
	}

	requestType, ok := requestMapping[internalType]
	if !ok {
		return nil, fmt.Errorf("no versioned request mapping for internal type %q with api version %d", internalType, version)
	}

	value := reflect.New(requestType.Elem()).Interface()
	versionedRequest, _ := value.(VersionedRequest)

	return versionedRequest, nil
}

func GetResponse(resp interface{}, version uint) (interface{}, error) {
	responseMapping, ok := internalToResponse[version]
	if !ok {
		return nil, fmt.Errorf("api version %d not registered on the versioned request to internal request mapping", version)
	}

	incomingType := reflect.TypeOf(resp)

	responseType, ok := responseMapping[incomingType]
	if !ok {
		return nil, fmt.Errorf("no versioned response mapping for internal type %q with api version %d", incomingType, version)
	}

	value := reflect.New(responseType.Elem()).Interface()
	versionedRequest, _ := value.(VersionedResponse)

	return versionedRequest.FromInternalResponse(resp), nil
}
