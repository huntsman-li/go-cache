package cache

import (
	"os"
	"path"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_FileCacher(t *testing.T) {
	Convey("Test file cache adapter", t, func() {
		dir := path.Join(os.TempDir(), "data/caches")
		os.RemoveAll(dir)
		testAdapter(Options{
			Adapter:       "file",
			AdapterConfig: dir,
			Interval:      2,
		})
	})
}