package input

import (
	"encoding/json"
	"reflect"

	"github.com/sirupsen/logrus"

	validator "gopkg.in/go-playground/validator.v9"
)

type Nullable struct {
	Data  interface{}
	Wrote bool
}

func (u *Nullable) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Data)
}

func (u *Nullable) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &u.Data)
	u.Wrote = true
	return err
}

func GetObjectMap(obj interface{}) map[string]interface{} {
	objVal := reflect.ValueOf(obj)
	objTyp := reflect.TypeOf(obj)
	toReturn := make(map[string]interface{})
	for i := 0; i < objVal.NumField(); i++ {
		cur := objVal.Field(i)
		if cur.FieldByName("Wrote").Bool() {
			toReturn[objTyp.Field(i).Name] = cur.FieldByName("Data").Interface()
		}
	}
	return toReturn
}

func ValidateObject(valid *validator.Validate, obj interface{}) error {
	objVal := reflect.ValueOf(obj)
	objTyp := reflect.TypeOf(obj)
	for i := 0; i < objVal.NumField(); i++ {
		cur := objVal.Field(i)
		curType := cur.Type()
		fieldMap := make(map[string]interface{})
		for j := 0; j < curType.NumField(); j++ {
			fieldMap[curType.Field(j).Name] = cur.Field(j).Interface()
		}
		fieldMap["name"] = objTyp.Field(i).Name
		logrus.WithFields(logrus.Fields(fieldMap)).Debug("ValiedateObject")
		if cur.FieldByName("Wrote").Bool() {
			if err := valid.Var(
				cur.FieldByName("Data").Interface(),
				objTyp.Field(i).Tag.Get("validate"),
			); err != nil {
				return err
			}
		}
	}
	return nil
}
