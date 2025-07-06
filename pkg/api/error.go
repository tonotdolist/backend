package api

import errors2 "errors"

var errors = make(map[uint]map[error]interface{})

func RegisterError(apiVersion uint, target interface{}, mapping error) {
	versionMapping := errors[apiVersion]

	if versionMapping == nil {
		versionMapping = make(map[error]interface{})
		errors[apiVersion] = versionMapping
	}

	errors[apiVersion][mapping] = target
}

func GetError(version ApiVersionHandler, mapping error) interface{} {
	apiVersion := version.GetApiVersion()

	unwrappedErr := errors2.Unwrap(mapping)
	if unwrappedErr == nil {
		unwrappedErr = mapping
	}

	if versionMapping, ok := errors[apiVersion]; ok {
		if err, ok := versionMapping[unwrappedErr]; ok {
			return err
		}
	}

	return version.GetInternalError()
}
