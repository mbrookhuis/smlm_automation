package cmdexecutor

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

const (
	// failedCMD failed os command
	failedCMD = "failed to execute command "

	// executeScript generic identifier for os methods
	executeScript = "Execute Script"

	// failedCreateDir failed os command
	failedCreateDir = "failed to create directory"

	// executeScript generic identifier for os methods
	createDirectory = "create directory"
)

// cmdUtils contains use case configuration
type cmdUtils struct {
	logger      *zap.Logger
	execCommand func(name string, arg ...string) *exec.Cmd
}

// NewCMDExecutor new command line executor with configuration
//
//	@param logger
//	@return ICMDExecutor
func NewCMDExecutor(logger *zap.Logger) ICMDExecutor {
	return &cmdUtils{
		logger:      logger,
		execCommand: exec.Command,
	}
}

// ExecuteCommand executes os command
//
//	@receiver os
//	@param name
//	@param args
//	@return []string
//	@return error
func (cmdUtils *cmdUtils) ExecuteCommand(binaryName string, args []string) ([]string, error) {
	cmdUtils.logger.Info(fmt.Sprintf("%v Started", executeScript))

	executeCMD := cmdUtils.execCommand(binaryName, args...)

	// Create StdOut and StdErr Multi Writer
	var stdoutBuf, stderrBuf bytes.Buffer
	executeCMD.Stdout = io.MultiWriter(os.Stdout, &stdoutBuf)
	executeCMD.Stderr = io.MultiWriter(os.Stderr, &stderrBuf)

	// Execute Command
	err := executeCMD.Run()

	// Convert StdOut and StdErr bytes to string
	stdOut, stdErr := stdoutBuf.String(), stderrBuf.String()
	cmdUtils.logger.Debug(fmt.Sprintf("%v /n StdOut-%v /n StdErr-%v",
		executeScript, stdOut, stdErr))

	if err != nil {
		cmdUtils.logger.Error(fmt.Sprintf("%v error- %v",
			executeScript, err))
		return nil, errors.New(failedCMD + err.Error() + ";stdErr:" + stdErr)
	}

	if stdErr != "" {
		cmdUtils.logger.Error(fmt.Sprintf("%v %v error- %v",
			executeScript, failedCMD, stdErr))
		return nil, errors.New(failedCMD + stdErr)
	}

	cmdUtils.logger.Info(fmt.Sprintf("%v Completed", executeScript))
	// Sending StdOut Response
	return strings.Split(string(stdOut), "\n"), nil
}

func (cmdUtils *cmdUtils) CreateDirectory(path string) error {
	cmdUtils.logger.Info(fmt.Sprintf("%v Started", createDirectory))
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			cmdUtils.logger.Error(fmt.Sprintf("%v error- %v", failedCreateDir, err))
			return fmt.Errorf("failed to create directory: %s", failedCreateDir)
		}
	}
	return nil
}
