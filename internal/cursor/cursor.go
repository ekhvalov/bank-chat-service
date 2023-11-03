package cursor

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func Encode(data any) (string, error) {
	js, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON marshal error: %v", err)
	}
	return base64.URLEncoding.EncodeToString(js), nil
}

func Decode(in string, to any) error {
	js, err := base64.URLEncoding.DecodeString(in)
	if err != nil {
		return fmt.Errorf("base64 decode error: %v", err)
	}
	if err := json.Unmarshal(js, to); err != nil {
		return fmt.Errorf("JSON unmarshal error: %v", err)
	}
	return nil
}
