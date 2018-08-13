package request_test

import (
	"fmt"
	"testing"

	_http "net/http"

	"github.com/Jeff-All/magi/endpoints/request"
	http "github.com/Jeff-All/magi/mock/http"
	util "github.com/Jeff-All/magi/util"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPUT(
	t *testing.T,
) {
	Convey("endpoints/request.PUT()", t, func() {
		Convey("Error: Reading Body", func() {
			mockIOUTil := util.MockIOUtil{}
			mockIOUTil.Mock.AddCall("ReadAll", []byte{}, fmt.Errorf("test error"))
			util.IOUtil = mockIOUTil

			mockResponseWriter := http.ResponseWriter{}

			r := _http.Request{}

			toTest := request.BaseRequest{}

			returnedError := toTest.PUT(&mockResponseWriter, &r)

			So(returnedError, ShouldNotBeNil)
			So(returnedError.Error(), ShouldEqual, "test error")

			mockIOUTil.Mock.AssertCall(t, "ReadAll", r.Body)
			mockIOUTil.Mock.AssertEnd(t)

			mockResponseWriter.Mock.AssertEnd(t)
		})
		Convey("Error: Parsing Body", func() {
		})
		Convey("Error: Creating Model", func() {
		})
		Convey("Error: Response Serialization", func() {
		})
		Convey("Error: Writing Response", func() {
		})
		Convey("No Error", func() {
		})
	})
}
