package utils

import (
	"dcss/pkg/utils/ssh"
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"os"
	"path"
	"path/filepath"
)

func Ssh(host, port, user, password, cmd string) (string, error) {
	c, err := ssh.NewClient(host, port, user, password)
	if err != nil {
		return "", err
	}
	defer c.Close()

	output, err := c.Output(cmd)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func SuSsh(host, port, user, password, cmd string) (string, error) {
	c, err := ssh.NewClient(host, port, user, password)
	if err != nil {
		return "", err
	}
	defer c.Close()

	output, err := c.ExecSu(cmd, password)
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func VerifySsh(host string, port int, user, password string) error {
	var config = &ssh.Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
	}

	if len(password) == 0 {
		home, err := homedir.Dir()
		if err != nil {
			return err
		}
		key := filepath.ToSlash(path.Join(home, ".ssh/id_rsa"))
		if IsFile(key) {
			return errors.New("id_rsa不存在，" + key)
		}
		config.KeyFiles = []string{key}
	}

	c, err := ssh.New(config)
	if err != nil {
		return err
	}
	defer c.Close()
	return nil

}

type SshClient struct {
	*ssh.Client
	IsRoot bool
}

func (s *SshClient) SshWriteFile(filePath, content string, perm os.FileMode) error {
	dirName := path.Dir(filePath)
	if _, err := s.SFTPClient.Stat(dirName); os.IsNotExist(err) {
		err := s.SFTPClient.MkdirAll(dirName)
		if err != nil {
			return fmt.Errorf("创建远端文件夹失败, err: %s", err.Error())
		}
	}

	remoteFile, err := s.SFTPClient.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建远端文件失败,err: %s", err)
	}
	defer remoteFile.Close()
	_, err = remoteFile.Write([]byte(content))
	if err != nil {
		return fmt.Errorf("写入远端文件失败,err: %s", err)
	}
	err = remoteFile.Chmod(perm)
	if err != nil {
		return fmt.Errorf("修改远端文件权限失败,err: %s", err)
	}
	return nil
}

func (s *SshClient) CombinedOutput(cmd string) ([]byte, error) {
	session, err := s.Client.SSHClient.NewSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return session.CombinedOutput(cmd)
}

func (s *SshClient) RunRemoteScript(filePath, content, scriptType string, perm os.FileMode, sudo bool) ([]byte, error) {
	err := s.SshWriteFile(filePath, fmt.Sprintf("set -e\n%s", content), perm)
	if err != nil {
		return nil, err
	}

	var cmd string
	if scriptType == "perl" {
		cmd = fmt.Sprintf("perl %s", filePath)
	} else if scriptType == "python" {
		cmd = fmt.Sprintf("python %s", filePath)
	} else {
		cmd = fmt.Sprintf("bash %s", filePath)
	}

	if sudo {
		cmd = "sudo -i " + cmd
	}

	return s.CombinedOutput(cmd)
}
