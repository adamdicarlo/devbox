package nix

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alessio/shellescape"
	"github.com/pkg/errors"
	"go.jetpack.io/devbox/internal/boxcli/usererr"
	"go.jetpack.io/devbox/internal/debug"
)

func RunScript(nixShellFilePath string, projectDir string, scriptPath string, additionalEnv []string) error {
	if scriptPath == "" {
		return errors.New("attempted to run script but did not specify script name")
	}

	vaf, err := PrintDevEnv(nixShellFilePath)
	if err != nil {
		return err
	}

	nixEnv := []string{}
	for k, v := range vaf.Variables {
		if v.Type == "exported" {
			nixEnv = append(nixEnv, fmt.Sprintf("%s=%s", k, shellescape.Quote(v.Value.(string))))
		}
	}

	cmd := exec.Command("sh", "-c", scriptPath)
	cmd.Env = append(nixEnv, additionalEnv...)
	cmd.Dir = projectDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	debug.Log("Executing: %v", cmd.Args)
	err = cmd.Run()
	if err != nil {
		// Report error as exec error when executing scripts.
		err = usererr.NewExecError(err)
	}
	return errors.WithStack(err)
}
