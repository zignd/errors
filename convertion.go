package errors

import "errors"

// toMapsSlice converts an error and its causes to flat slice of maps where each map represents an error.
func toMapsSlice(err error) []map[string]any {
	errMaps := make([]map[string]any, 0)

	if err == nil {
		return errMaps
	}

	currentErr := err
	for {
		errMap, errCause := toMapAndCause(currentErr)
		errMaps = append(errMaps, errMap)
		if errCause == nil {
			break
		}
		currentErr = errCause
	}

	return errMaps
}

// toMapAndCause converts an error to a map and extracts the cause.
func toMapAndCause(err error) (map[string]any, error) {
	errMap := make(map[string]any)
	var errCause error

	if e, ok := err.(*Err); ok {
		errMap["message"] = e.Message
		if e.Data != nil {
			errMap["data"] = e.Data
		}
		errMap["stack"] = e.Stack
		errCause = e.Cause
	} else {
		errMap["message"] = err.Error()
		errCause = errors.Unwrap(err)
	}

	return errMap, errCause
}
