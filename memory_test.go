package cache

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_MemoryCacher(t *testing.T) {
	Convey("Test memory cache adapter", t, func() {
		testAdapter(Options{
			Interval: 2,
		})
	})
}