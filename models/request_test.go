package models_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/Jeff-All/magi/data"
	"github.com/Jeff-All/magi/mock"

	"github.com/Jeff-All/magi/models"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCreateRequest(
	t *testing.T,
) {
	Convey("models.Request.CreateRequest()", t, func() {
		Convey("Error: Already Created", func() {
			mockDB := data.Mock{Mock: mock.NewMock()}

			models.DB = &mockDB

			createdAt := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

			request := models.Request{
				BaseModel: models.BaseModel{
					ID:        100,
					CreatedAt: &createdAt,
				},
			}

			response_error := models.Requests.Create(&request)

			// Error Assertions
			So(response_error, ShouldNotBeNil)
			So(response_error.Error(), ShouldEqual, "Request was already created")

			// Mock Assertions
			mockDB.Mock.AssertEnd(t)
		})
		Convey("Error: Database Error", func() {
			mockDB := data.NewMock()
			mockDB.Mock.AddCall("Create")
			mockDB.Mock.AddCall("GetError", fmt.Errorf("Test Error"))

			models.DB = &mockDB

			request := models.Request{
				BaseModel: models.BaseModel{ID: 100},
				Agency:    models.Agency{},
			}

			response_error := models.Requests.Create(&request)

			// Error Assertions
			So(response_error, ShouldNotBeNil)
			So(response_error.Error(), ShouldEqual, "Test Error")

			// Mock Assertions
			mockDB.Mock.AssertCall(t, "Create", &request)
			mockDB.Mock.AssertCall(t, "GetError")
			mockDB.Mock.AssertEnd(t)
		})
		Convey("No Error", func() {
			mockDB := data.NewMock()
			mockDB.Mock.AddCall("Create")
			mockDB.Mock.AddCall("GetError", nil)

			models.DB = &mockDB

			request := models.Request{
				BaseModel: models.BaseModel{ID: 100},
				Agency:    models.Agency{},
			}

			response_error := models.Requests.Create(&request)

			// Error Assertions
			So(response_error, ShouldBeNil)

			// Mock Assertions
			mockDB.Mock.AssertCall(t, "Create", &request)
			mockDB.Mock.AssertCall(t, "GetError")
			mockDB.Mock.AssertEnd(t)
		})
	})
}
