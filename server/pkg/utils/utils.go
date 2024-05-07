package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var Random = rand.New(rand.NewSource(time.Now().UnixNano()))

func ExecCommand(cmdStr string, args ...string) (string, error) {
	timeout := viper.GetInt("CMD.TIMEOUT")
	if timeout < 5 {
		timeout = 5
	}
	return ExecCommandWithTimeout(timeout, cmdStr, args...)
}

func GetNowTime() string {
	return time.Now().Format("2006/01/02 15:04:05")
}

func GetRandomStrNum(num int) string {
	randomNum := Random.Intn(num)
	return strconv.Itoa(randomNum)
}

func CheckFileOrDirIsExist(path string) (bool, error) {
	if len(strings.TrimSpace(path)) == 0 {
		return false, errors.New("path is empty")
	}

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CheckUserIsExist(user string) (bool, error) {
	_, err := ExecCommand("id", user)
	if err != nil {
		return false, err
	}
	return true, nil
}

func CheckMailProcessIsRunning(user string) (bool, error) {
	numStr, err := ExecCommand("pgrep", "-u", user, "-c")
	if err != nil {
		return false, err
	}

	numStr = strings.Trim(numStr, "\n")
	num, err := strconv.Atoi(numStr)

	if err != nil {
		log.Println(err)
		return false, err
	}

	if num > 0 {
		return true, nil
	}

	return false, nil
}

func GetAllIPList() ([]string, error) {
	ipList := make([]string, 5)
	stdout, err := ExecCommand("hostname", "-I")
	if err != nil {
		return nil, err
	}
	stdout = strings.Trim(stdout, "\n")
	stdout = strings.TrimSpace(stdout)

	ipList = strings.Split(stdout, " ")
	return ipList, nil
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func IsIpv4(ip string) bool {
	ipRex, err := regexp.Compile(`^((\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])\.){3}(\d|[1-9]\d|1\d\d|2[0-4]\d|25[0-5])(?::(?:[0-9]|[1-9][0-9]{1,3}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5]))?$`)
	if err != nil {
		return false
	}
	return ipRex.MatchString(ip)
}

func ReverseAny(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

func GetCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func GetJson(in interface{}) string {
	v, _ := json.Marshal(in)
	var dst bytes.Buffer
	_ = json.Indent(&dst, v, "", "")

	return strings.ReplaceAll(dst.String(), "\n", "")
}
