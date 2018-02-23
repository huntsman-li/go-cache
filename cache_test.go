package cache

import (
	"encoding/gob"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Version(t *testing.T) {
	Convey("Check package version", t, func() {
		So(Version(), ShouldEqual, _VERSION)
	})
}

func Test_Cacher(t *testing.T) {
	Convey("Use cache middleware", t, func() {
		_, err := Cacher()
		So(err, ShouldBeNil)
	})

	Convey("Register invalid adapter", t, func() {
		Convey("Adatper not exists", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			Cacher(Options{
				Adapter: "fake",
			})
		})

		Convey("Provider value is nil", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			Register("fake", nil)
		})

		Convey("Register twice", func() {
			defer func() {
				So(recover(), ShouldNotBeNil)
			}()

			Register("memory", &MemoryCacher{})
		})
	})
}

func testAdapter(opt Options) {
	Convey("Basic operations", func() {
		c, err := Cacher(opt)
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

		gob.Register(opt)
		So(c.Put("struct", opt, 0), ShouldBeNil)
	})

	Convey("Increase and decrease operations", func() {
		c, err := Cacher(opt)
		So(err, ShouldBeNil)
		So(c.Incr("404"), ShouldNotBeNil)
		So(c.Decr("404"), ShouldNotBeNil)

		So(c.Put("int", 0, 0), ShouldBeNil)
		So(c.Put("int32", int32(0), 0), ShouldBeNil)
		So(c.Put("int64", int64(0), 0), ShouldBeNil)
		So(c.Put("uint", uint(0), 0), ShouldBeNil)
		So(c.Put("uint32", uint32(0), 0), ShouldBeNil)
		So(c.Put("uint64", uint64(0), 0), ShouldBeNil)
		So(c.Put("string", "hi", 0), ShouldBeNil)

		So(c.Decr("uint"), ShouldNotBeNil)
		So(c.Decr("uint32"), ShouldNotBeNil)
		So(c.Decr("uint64"), ShouldNotBeNil)

		So(c.Incr("int"), ShouldBeNil)
		So(c.Incr("int32"), ShouldBeNil)
		So(c.Incr("int64"), ShouldBeNil)
		So(c.Incr("uint"), ShouldBeNil)
		So(c.Incr("uint32"), ShouldBeNil)
		So(c.Incr("uint64"), ShouldBeNil)

		So(c.Decr("int"), ShouldBeNil)
		So(c.Decr("int32"), ShouldBeNil)
		So(c.Decr("int64"), ShouldBeNil)
		So(c.Decr("uint"), ShouldBeNil)
		So(c.Decr("uint32"), ShouldBeNil)
		So(c.Decr("uint64"), ShouldBeNil)

		So(c.Incr("string"), ShouldNotBeNil)
		So(c.Decr("string"), ShouldNotBeNil)

		So(c.Get("int"), ShouldEqual, 0)
		So(c.Get("int32"), ShouldEqual, 0)
		So(c.Get("int64"), ShouldEqual, 0)
		So(c.Get("uint"), ShouldEqual, 0)
		So(c.Get("uint32"), ShouldEqual, 0)
		So(c.Get("uint64"), ShouldEqual, 0)
	})
}
