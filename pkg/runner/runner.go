package runner

import (
	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor/output"
)

func NewRunner() *GradleRunner {
	return &GradleRunner{}
}

type GradleRunner struct {
}

func (r *GradleRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	// the Gradle executor does not support files
	if execution.Content.IsFile() {
		output.PrintEvent("using file", execution)
		// TODO implement file based test content for string, git-file, file-uri
		//      or remove if not used
	}

	if execution.Content.IsDir() {
		output.PrintEvent("using dir", execution)
		// TODO implement file based test content for git-dir
		//      or remove if not used
	}

	// TODO run executor here

	// error result should be returned if something is not ok
	// return result.Err(fmt.Errorf("some test execution related error occured"))

	// TODO return ExecutionResult
	return testkube.ExecutionResult{
		Status: testkube.StatusPtr(testkube.SUCCESS_ExecutionStatus),
		Output: "exmaple test output",
	}, nil
}
