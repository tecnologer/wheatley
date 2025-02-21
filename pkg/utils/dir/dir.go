package dir

import (
	"path"
	"runtime"
)

func CallerDir() string {
	_, filename, _, _ := runtime.Caller(1) //nolingt:dogsled

	return path.Dir(filename)
}
