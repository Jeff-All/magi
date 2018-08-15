package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Jeff-All/magi/errors"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/Jeff-All/magi/middleware"
)

func TestHandleError(t *testing.T) {
	Convey("Unit Tests for middleware.HandleError()", t, func() {
		Convey("No Error", func() {
			mockWriter := httptest.NewRecorder()
			mockRequest := httptest.NewRequest("POST", "/test", strings.NewReader(""))
			handler := middleware.HandleError("test", func(http.ResponseWriter, *http.Request) error {
				return nil
			})
			handler.ServeHTTP(mockWriter, mockRequest)
			So(mockWriter.Body, ShouldNotBeNil)
			So(string(mockWriter.Body.Bytes()), ShouldEqual, "")
		})
		Convey("CodedError", func() {
			mockWriter := httptest.NewRecorder()
			mockRequest := httptest.NewRequest("POST", "/test", strings.NewReader(""))
			handler := middleware.HandleError("test", func(http.ResponseWriter, *http.Request) error {
				return errors.CodedError{
					Code:     100,
					Message:  "test message",
					HTTPCode: 101,
					Err:      fmt.Errorf("test error message"),
				}
			})
			handler.ServeHTTP(mockWriter, mockRequest)
			So(mockWriter.Body, ShouldNotBeNil)
			So(string(mockWriter.Body.Bytes()), ShouldEqual, `{"Message":"test message","Code":100,"HTTPCode":101}`)
			So(mockWriter.Code, ShouldEqual, 101)
		})
		Convey("!CodedError", func() {
			mockWriter := httptest.NewRecorder()
			mockRequest := httptest.NewRequest("POST", "/test", strings.NewReader(""))
			handler := middleware.HandleError("test", func(http.ResponseWriter, *http.Request) error {
				return fmt.Errorf("Test")
			})
			handler.ServeHTTP(mockWriter, mockRequest)
			So(mockWriter.Body, ShouldNotBeNil)
			So(string(mockWriter.Body.Bytes()), ShouldEqual, `{"Message":"Internal Server Error.","Code":0,"HTTPCode":500}`)
			So(mockWriter.Code, ShouldEqual, 500)
		})
	})
}
