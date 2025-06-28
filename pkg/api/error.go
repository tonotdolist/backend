package api

var errors = make(map[uint]map[error]interface{})

func RegisterError(apiVersion uint, target interface{}, mapping error) {
	versionMapping := errors[apiVersion]

	if versionMapping == nil {
		versionMapping = make(map[error]interface{})
		errors[apiVersion] = versionMapping
	}

	versionMapping[mapping] = target
}

func getError(version ApiVersion, mapping error) interface{} {
	apiVersion := version.GetApiVersion()

	if versionMapping, ok := errors[apiVersion]; ok {
		if err, ok := versionMapping[mapping]; ok {
			return err
		}
	}

	return version.GetInternalError()
}
