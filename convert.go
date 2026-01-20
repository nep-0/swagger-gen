package main

import (
	"bytes"
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// convertToJSONType recursively converts YAML interface{} types to JSON-compatible types
// while preserving order using ordered JSON generation
func convertToJSONType(val interface{}) interface{} {
	switch v := val.(type) {
	case yaml.MapSlice:
		// Return the MapSlice as-is so we can handle it specially during JSON generation
		return v
	case map[interface{}]interface{}:
		result := make(map[string]interface{})
		for key, value := range v {
			if str, ok := key.(string); ok {
				result[str] = convertToJSONType(value)
			}
		}
		return result
	case []interface{}:
		result := make([]interface{}, len(v))
		for i, item := range v {
			result[i] = convertToJSONType(item)
		}
		return result
	default:
		return v
	}
}

// orderedJSONMarshal converts a value to JSON while preserving order for yaml.MapSlice
func orderedJSONMarshal(val interface{}) ([]byte, error) {
	var buf bytes.Buffer
	err := marshalValue(&buf, val)
	return buf.Bytes(), err
}

func marshalValue(buf *bytes.Buffer, val interface{}) error {
	switch v := val.(type) {
	case yaml.MapSlice:
		buf.WriteByte('{')
		for i, item := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			// Marshal the key
			keyBytes, err := json.Marshal(item.Key)
			if err != nil {
				return err
			}
			buf.Write(keyBytes)
			buf.WriteByte(':')
			// Marshal the value recursively
			err = marshalValue(buf, item.Value)
			if err != nil {
				return err
			}
		}
		buf.WriteByte('}')
		return nil
	case []interface{}:
		buf.WriteByte('[')
		for i, item := range v {
			if i > 0 {
				buf.WriteByte(',')
			}
			err := marshalValue(buf, item)
			if err != nil {
				return err
			}
		}
		buf.WriteByte(']')
		return nil
	default:
		// For other types, use standard JSON marshaling
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		buf.Write(jsonBytes)
		return nil
	}
}
