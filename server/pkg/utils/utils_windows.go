package utils

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

func ExecCommandWithTimeout(timeout int, cmdStr string, args ...string) (string, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, cmdStr, args...)
	opBytes, err := cmd.CombinedOutput()
	if err != nil {
		if ctx.Err() != nil && ctx.Err() == context.DeadlineExceeded {
			//return "", errors.New(cmdStr + strings.Join(args, " ") + "命令执行超时")
			return "", errors.New(fmt.Sprintf("命令执行超时%ds", timeout))
		}
		return string(opBytes), err
	}
	return string(opBytes), nil
}
