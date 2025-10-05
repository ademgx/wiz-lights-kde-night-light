package dbusclient

import (
	"errors"
	"reflect"

	"github.com/godbus/dbus/v5"
)

// GetCurrentTemperature extracts the current color temperature value from a D-Bus signal.
// It expects the signal body to contain at least 2 elements, where the second element
// Returns the temperature as an integer or an error if the signal format is invalid
// or the required data is missing/malformed. Note: This may break if KWin night light signal format is changed
// in the future
func GetCurrentTemperature(signal *dbus.Signal) (int, error) {
	if len(signal.Body) < 2 {
		return 0, errors.New("invalid number of values in body of signal, expected >= 2")
	}
	params := signal.Body[1]
	paramType := reflect.TypeOf(params).Kind()
	if paramType != reflect.Map {
		return 0, errors.New("expected second parameter of signal body to be a map, got: " + paramType.String())
	}

	paramsMap, ok := params.(map[string]dbus.Variant)
	if !ok {
		return 0, errors.New("failed to cast params to map[string]interface{}")
	}

	value, ok := paramsMap["currentTemperature"]
	if !ok {
		return 0, errors.New("key 'currentTemperature' not present in signal body")
	}
	val, ok := value.Value().(uint32)
	if !ok {
		return 0, errors.New("could not convert temperature value to an integer")
	}
	return int(val), nil
}
