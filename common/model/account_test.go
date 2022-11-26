package model

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAccountWrongPath(t *testing.T) {

	Convey("Given a an email address", t, func() {
		email := EmailAddress("erik@domain.nu")

		email2 := EmailAddress("erikdomain.nu")

		Convey("When validated", func() {
			valid := email.IsValid()
			Convey("Then result should be true", func() {
				So(valid, ShouldBeTrue)
			})
			invalid := email2.IsValid()
			Convey("Then result should be false", func() {
				So(invalid, ShouldBeFalse)
			})
		})
	})
}
