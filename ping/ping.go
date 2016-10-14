package ping

import (
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"
)

const Int64Max int64 = 9223372036854775807

func Ping(dhost string) (int64, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		cmd = exec.Command("ping", "-c", "1", dhost)
	case "windows":
		cmd = exec.Command("ping", "-n", "1", dhost)
	default:
		return -1, fmt.Errorf("不支持的操作系统 %v", runtime.GOOS)
	}

	o, err := cmd.Output()
	if err != nil {
		return -1, err
	}

	return findTime(string(o))
}

func findTime(t string) (int64, error) {
	re, err := regexp.Compile(`[=,<](\d+\.?\d*)\s*ms`)
	if err != nil {
		return -1, err
	}

	res := re.FindStringSubmatch(t)

	if res == nil || len(res) < 2 {
		return -1, fmt.Errorf("超时")
	}

	f, err := strconv.ParseFloat(res[1], 32)
	if err != nil {
		return -1, err
	}

	return int64(f * float64(time.Millisecond)), nil
}
