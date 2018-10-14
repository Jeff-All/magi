package input_test

import (
	"encoding/json"

	"testing"

	"github.com/Jeff-All/magi/input"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Unit_Nullable(t *testing.T) {
	Convey("Unit Tests for Nullable", t, func() {
		Convey("Base Test", func() {
			type TestStruct struct {
				FieldA input.Nullable
				FieldB input.Nullable
				FieldC input.Nullable
			}
			t := TestStruct{}
			data := []byte(`{"FieldA":89,"FieldB":"testB"}`)
			err := json.Unmarshal(data, &t)
			So(err, ShouldBeNil)
			So(t.FieldA.Data, ShouldEqual, 89)
			So(t.FieldB.Data, ShouldEqual, "testB")
			So(t.FieldC.Data, ShouldBeNil)

			So(t.FieldA.Wrote, ShouldBeTrue)
			So(t.FieldB.Wrote, ShouldBeTrue)
			So(t.FieldC.Wrote, ShouldBeFalse)
		})
		Convey("Get Object Map", func() {
			type TestStruct struct {
				FieldA input.Nullable
				FieldB input.Nullable
				FieldC input.Nullable
			}
			t := TestStruct{}
			data := []byte(`{"FieldA":89,"FieldB":"testB"}`)
			err := json.Unmarshal(data, &t)
			So(err, ShouldBeNil)
			objMap := input.GetObjectMap(t)
			So(objMap, ShouldContainKey, "FieldA")
			So(objMap, ShouldContainKey, "FieldB")
			So(objMap, ShouldNotContainKey, "FieldC")

			So(objMap["FieldA"], ShouldEqual, 89)
			So(objMap["FieldB"], ShouldEqual, "testB")
		})
	})
}
