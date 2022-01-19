package intermediate

import (
	"os/exec"

	"github.com/daku10/concurrency-in-go-study/chapter-5/error-handling/lowlevel"
	"github.com/daku10/concurrency-in-go-study/chapter-5/error-handling/util"
)

type IntermediateErr struct {
	error
}

func RunJob(id string) error {
	const jobBinPath = "/bad/job/binary"
	isExecutable, err := lowlevel.IsGloballyExec(jobBinPath)
	if err != nil {
		// エラーをラップしましょう
		// return err
		return IntermediateErr{util.WrapError(err, "cannot run job %q: requisite binaries not available", id)}
	} else if isExecutable == false {
		return util.WrapError(nil, "job binary is not executable")
	}

	return exec.Command(jobBinPath, "--id="+id).Run()
}
