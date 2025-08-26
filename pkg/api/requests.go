package api

import (
	"fmt"
	"reflect"
)

var versionedRequestType = reflect.TypeOf((*VersionedRequest)(nil)).Elem()
var versionedResponseType = reflect.TypeOf((*VersionedResponse)(nil)).Elem()

var requestToInternal = make(map[uint]map[reflect.Type]reflect.Type)
var internalToResponse = make(map[uint]map[reflect.Type]reflect.Type)

type VersionedRequest interface {
	ToInternalRequest() interface{}
}

type VersionedResponse interface {
	FromInternalResponse(data interface{})
}

func RegisterRequest[TInternal any, TVersion any](apiVersion uint) {
	requestMapping := requestToInternal[apiVersion]

	if requestMapping == nil {
		requestMapping = make(map[reflect.Type]reflect.Type)
		requestToInternal[apiVersion] = requestMapping
	}

	internalType := reflect.TypeOf((*TInternal)(nil))
	versionedType := reflect.TypeOf((*TVersion)(nil))

	if internalType.Kind() == reflect.Ptr {
		internalType = internalType.Elem()
	}

	if versionedType.Kind() == reflect.Struct {
		versionedType = reflect.PointerTo(versionedType)
	}

	if !versionedType.Implements(versionedRequestType) {
		panic(fmt.Sprintf("%s does not implement VersionedRequest interface", versionedType))
	}

	versionedType = versionedType.Elem()

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

	if internalType.Kind() == reflect.Ptr {
		internalType = internalType.Elem()
	}

	if versionedType.Kind() == reflect.Struct {
		versionedType = reflect.PointerTo(versionedType)
	}

	if !versionedType.Implements(versionedResponseType) {
		panic(fmt.Sprintf("%s does not implement VersionedResponse interface", versionedType))
	}

	versionedType = versionedType.Elem()

	internalToResponse[apiVersion][internalType] = versionedType
}

func GetRequest(internalType reflect.Type, version uint) (VersionedRequest, error) {
	requestMapping, ok := requestToInternal[version]
	if !ok {
		return nil, fmt.Errorf("api version %d not registered on the versioned request to internal request mapping", version)
	}

	if internalType.Kind() == reflect.Ptr {
		internalType = internalType.Elem()
	}

	requestType, ok := requestMapping[internalType]
	if !ok {
		return nil, fmt.Errorf("no versioned request mapping for internal type %q with api version %d", internalType, version)
	}

	value := reflect.New(requestType).Interface()
	versionedRequest, _ := value.(VersionedRequest)

	return versionedRequest, nil
}

func GetResponse(resp interface{}, version uint) (interface{}, error) {
	responseMapping, ok := internalToResponse[version]
	if !ok {
		return nil, fmt.Errorf("api version %d not registered on the versioned request to internal request mapping", version)
	}

	commonRespType := reflect.TypeOf(resp)

	if commonRespType.Kind() == reflect.Ptr {
		commonRespType = commonRespType.Elem()
	}

	versionedRespType, ok := responseMapping[commonRespType]
	if !ok {
		return nil, fmt.Errorf("no versioned response mapping for internal type %q with api version %d", commonRespType, version)
	}

	value := reflect.New(versionedRespType).Interface()
	versionedRes, _ := value.(VersionedResponse)

	versionedRes.FromInternalResponse(resp)

	return versionedRes, nil
}
