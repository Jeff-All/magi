package endpoints_test

import (
	"testing"

	"github.com/Jeff-All/magi/endpoints"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_Regression_ParseWhereParam(t *testing.T) {
	Convey("Regression Tests for endpoints", t, func() {
		type TestModel struct {
			FieldA string
			FieldB string
		}
		Convey("Batch=vrgr:ID=1", func() {
			type TestModel struct {
				Batch string
				ID    int
			}
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "Batch=vrgr:ID=1",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "Batch=? AND ID=?")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 2)
			So(out_argument[0], ShouldEqual, "vrgr")
			So(out_argument[1], ShouldEqual, "1")
		})
	})
}

func Test_Unit_ParseWhereParam(t *testing.T) {
	Convey("Unit Tests for endpoints", t, func() {
		type TestModel struct {
			FieldA string
			FieldB string
		}
		Convey("=", func() {
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "FieldA=test",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "FieldA=?")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 1)
			So(out_argument[0], ShouldEqual, "test")
		})
		Convey("!=", func() {
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "FieldA!=test",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "FieldA!=?")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 1)
			So(out_argument[0], ShouldEqual, "test")
		})
		Convey("Like", func() {
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "FieldA%test%",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "FieldA LIKE ?")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 1)
			So(out_argument[0], ShouldEqual, "%test%")
		})
		Convey("Not Like", func() {
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "FieldA!%test%what",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "FieldA NOT LIKE ?")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 1)
			So(out_argument[0], ShouldEqual, "%test%what")
		})
		Convey("In", func() {
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "FieldA[test,test2]",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "FieldA IN(?,?)")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 2)
			So(out_argument[0], ShouldEqual, "test")
			So(out_argument[1], ShouldEqual, "test2")
		})
		Convey("Not In", func() {
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "FieldA![test,test2]",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "FieldA NOT IN(?,?)")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 2)
			So(out_argument[0], ShouldEqual, "test")
			So(out_argument[1], ShouldEqual, "test2")
		})
		Convey("And", func() {
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "FieldA![test,test2]:FieldB=test3",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "FieldA NOT IN(?,?) AND FieldB=?")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 3)
			So(out_argument[0], ShouldEqual, "test")
			So(out_argument[1], ShouldEqual, "test2")
			So(out_argument[2], ShouldEqual, "test3")
		})
		Convey("Or", func() {
			out_clause, out_argument, out_error := endpoints.ParseWhereParam(
				TestModel{}, "FieldA![test,test2]|FieldB=test3",
			)
			So(out_error, ShouldBeNil)
			So(out_clause, ShouldEqual, "FieldA NOT IN(?,?) OR FieldB=?")
			So(out_argument, ShouldNotBeNil)
			So(out_argument, ShouldHaveLength, 3)
			So(out_argument[0], ShouldEqual, "test")
			So(out_argument[1], ShouldEqual, "test2")
			So(out_argument[2], ShouldEqual, "test3")
		})
	})
}
