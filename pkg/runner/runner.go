package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/executor/output"
)

type Params struct {
	Datadir string // RUNNER_DATADIR
}

func NewRunner() *GradleRunner {
	params := Params{
		Datadir: os.Getenv("RUNNER_DATADIR"),
	}

	runner := &GradleRunner{
		params: params,
	}

	return runner
}

type GradleRunner struct {
	params Params
}

func (r *GradleRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {
	// check that the datadir exists
	_, err = os.Stat(r.params.Datadir)
	if errors.Is(err, os.ErrNotExist) {
		return result, err
	}

	// the Gradle executor does not support files
	if execution.Content.IsFile() {
		return result.Err(fmt.Errorf("executor only support git-dir based test")), nil
	}

	// check settings.gradle or build.gradle files exist
	directory := filepath.Join(r.params.Datadir, "repo")
	settingsGradle := filepath.Join(directory, "settings.gradle")
	_, err = os.Stat(settingsGradle)
	if errors.Is(err, os.ErrNotExist) {
		return result.Err(fmt.Errorf("no settings.gradle file found for test")), nil
	}

	// determine the Gradle command to use
	gradleCommand := "gradle"
	gradleWrapper := filepath.Join(directory, "gradlew")
	_, err = os.Stat(gradleWrapper)
	if err == nil {
		// then we use the wrapper instead
		gradleCommand = "./gradlew"
	}

	// prepare the ENVs to use during Gradle execution
	envs := ""
	for key, value := range execution.Envs {
		env := fmt.Sprintf("%s=%s", key, value)
		envs = envs + env + " "
	}

	// pass additional executor arguments/flags to Gradle
	args := []string{"--no-daemon"}
	args = append(args, execution.Args...)

	output.PrintEvent("Running", directory, envs+gradleCommand, args)
	output, err := executor.Run(directory, envs+gradleCommand, args...)
	if err != nil {
		return result.Err(err), nil
	}

	return testkube.ExecutionResult{
		Status:     testkube.StatusPtr(testkube.SUCCESS_ExecutionStatus),
		Output:     string(output),
		OutputType: "text/plain",
	}, nil
}
