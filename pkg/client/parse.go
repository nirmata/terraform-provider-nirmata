package client

import (
	"encoding/json"
	"fmt"
)

// ParseCollection parses JSON bytes to a slice of generic maps
func ParseCollection(b []byte) ([]map[string]interface{}, error) {
	var collection []map[string]interface{}
	err := json.Unmarshal(b, &collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// ParseObject parses JSON bytes to a generic map
func ParseObject(b []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal(b, &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

// ParseID parses JSON bytes to an ID
func ParseID(b []byte) (ID, error) {
	var jsonData []map[string]interface{}
	err := json.Unmarshal(b, &jsonData)
	if err != nil {
		return nil, err
	}

	return ParseIDFromMap(jsonData[0])
}

// ParseIDFromMap parses a JSON data map to an ID
func ParseIDFromMap(data map[string]interface{}) (ID, error) {
	if data["id"] == nil {
		return nil, fmt.Errorf("Failed to parse ID. Missing attribute 'id' in %+v", data)
	}

	if data["modelIndex"] == nil {
		return nil, fmt.Errorf("Failed to parse ID. Missing attribute 'modelIndex' in %+v", data)
	}

	if data["service"] == nil {
		return nil, fmt.Errorf("Failed to parse ID. Missing attribute 'service' in %+v", data)
	}

	uuid := data["id"].(string)
	modelIndex := data["modelIndex"].(string)
	s := data["service"].(string)

	service, err := ParseService(s)
	if err != nil {
		return nil, err
	}

	return NewID(service, modelIndex, uuid), nil
}
