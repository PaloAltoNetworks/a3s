package authorizer

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.aporeto.io/elemental"
)

func TestMockOperationTransformer(t *testing.T) {

	Convey("Given a MockOperationTransformer and an elemental operation", t, func() {

		mockOperationTransformer := NewMockOperationTransformer()

		operation := elemental.OperationRetrieveMany

		Convey("Calling Transform without mock should work", func() {
			op := mockOperationTransformer.Transform(operation)
			So(op, ShouldNotBeNil)
			So(len(op), ShouldEqual, 0)
		})

		Convey("Calling Transform with mock should work", func() {
			mockOperationTransformer.MockTransform(t, func(elemental.Operation) string {
				return "get"
			})
			op := mockOperationTransformer.Transform(operation)
			So(op, ShouldNotBeNil)
			So(op, ShouldEqual, "get")
		})
	})
}
