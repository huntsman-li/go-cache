package cache

import (
	"testing"
	"time"

	"github.com/Unknwon/com"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/huntsman-li/go-cache"
)

func Test_MemcacheCacher(t *testing.T) {
	Convey("Test memcache cache adapter", t, func() {
		opt := cache.Options{
			Adapter:       "memcache",
			AdapterConfig: "127.0.0.1:9090",
		}

		Convey("Basic operations", func() {
			c, err := cache.Cacher(opt)
			So(err, ShouldBeNil)

			So(c.Put("uname", "unknwon", 1), ShouldBeNil)
			So(c.Put("uname2", "unknwon2", 1), ShouldBeNil)
			So(c.IsExist("uname"), ShouldBeTrue)

			So(c.Get("404"), ShouldBeNil)
			So(c.Get("uname").(string), ShouldEqual, "unknwon")

			time.Sleep(1 * time.Second)
			So(c.Get("uname"), ShouldBeNil)
			time.Sleep(1 * time.Second)
			So(c.Get("uname2"), ShouldBeNil)

			So(c.Put("uname", "unknwon", 0), ShouldBeNil)
			So(c.Delete("uname"), ShouldBeNil)
			So(c.Get("uname"), ShouldBeNil)

			So(c.Put("uname", "unknwon", 0), ShouldBeNil)
			So(c.Flush(), ShouldBeNil)
			So(c.Get("uname"), ShouldBeNil)
		})

		Convey("Increase and decrease operations", func() {
			c, err := cache.Cacher(opt)
			So(err, ShouldBeNil)

			So(c.Incr("404"), ShouldNotBeNil)
			So(c.Decr("404"), ShouldNotBeNil)

			So(c.Put("int", 0, 0), ShouldBeNil)
			So(c.Put("int64", int64(0), 0), ShouldBeNil)
			So(c.Put("string", "hi", 0), ShouldBeNil)

			So(c.Incr("int"), ShouldBeNil)
			So(c.Incr("int64"), ShouldBeNil)

			So(c.Decr("int"), ShouldBeNil)
			So(c.Decr("int64"), ShouldBeNil)

			So(c.Incr("string"), ShouldNotBeNil)
			So(c.Decr("string"), ShouldNotBeNil)

			So(com.StrTo(c.Get("int").(string)).MustInt(), ShouldEqual, 0)
			So(com.StrTo(c.Get("int64").(string)).MustInt64(), ShouldEqual, 0)

			So(c.Flush(), ShouldBeNil)
		})
	})
}