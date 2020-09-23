package expensify

import (
	"encoding/json"
	"reflect"
)

const (
	jobTypeCreate = "create"
)

type jobRequest struct {
	Type        string `json:"type"`
	Credentials struct {
		PartnerUserID     string `json:"partnerUserID"`
		PartnerUserSecret string `json:"partnerUserSecret"`
	} `json:"credentials"`
	InputSettings *inputSettings `json:"inputSettings"`
}

type inputSettings struct {
	Type string `json:"type"`
	data interface{}
}

// MarshalJSON implements json.Marshaler.
func (i inputSettings) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["type"] = i.Type

	if i.data != nil {
		for k, v := range structToMap(i.data) {
			m[k] = v
		}
	}

	return json.Marshal(m)
}

func structToMap(item interface{}) map[string]interface{} {
	res := map[string]interface{}{}
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get("json")
		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMap(field)
			} else {
				res[tag] = field
			}
		}
	}
	return res
}
