package lowlevel

import (
	"os"

	"github.com/daku10/concurrency-in-go-study/chapter-5/error-handling/util"
)

type LowLevelErr struct {
	error
}

func IsGloballyExec(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, LowLevelErr{util.WrapError(err, err.Error())}
	}
	return info.Mode().Perm() & 0100 == 0100, nil
}
