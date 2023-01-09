package util

import "encoding/json"

// old and new must be address
func TransformStruct(old interface{}, new interface{}) error {
	bs, err := json.Marshal(old)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, new)
	if err != nil {
		return err
	}
	return nil
}
