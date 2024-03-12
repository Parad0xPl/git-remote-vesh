package config

import (
	"fmt"
	"os"
	"path"
	"strings"

	"gopkg.in/yaml.v2"
)

type EncConfig struct {
	SSHPort            int    `yaml:"ssh_port"`
	SSHAddress         string `yaml:"ssh_address"`
	SSHUser            string `yaml:"ssh_user"`
	SSHIdentityFile    string `yaml:"ssh_identity_file"`
	SSHMountPath       string `yaml:"ssh_mount_path"`
	SSHRemotePath      string `yaml:"ssh_remote_path"`
	VeraCryptMountPath string `yaml:"veracrypt_mount_path"`
	VeraCryptVaultPath string `yaml:"veracrypt_vault_path"`
	RepoPath           string `yaml:"repo_path"`
	RemoteName         string
}

type SSHFSParams struct {
	SSHPort         int
	SSHAddress      string
	SSHUser         string
	SSHIdentityFile string
	SSHMountPath    string
	SSHRemotePath   string
}

type VeraCryptParams struct {
	VeraCryptMountPath string
	VeraCryptVaultPath string
}

func GetConfig() (EncConfig, error) {
	config_raw, err := os.ReadFile(".gitenc")
	if err != nil {
		return EncConfig{}, fmt.Errorf("can't read config file: %v", err)
	}
	config := EncConfig{}
	config.RemoteName = os.Args[1]
	config.RepoPath = os.Args[2]
	err = yaml.Unmarshal(config_raw, &config)
	if err != nil {
		return EncConfig{}, fmt.Errorf("can't parse config file: %v", err)
	}
	if config.SSHPort == 0 {
		config.SSHPort = 22
	}

	var mountPath string
	if strings.HasSuffix(config.VeraCryptMountPath, ":") {
		mountPath = config.VeraCryptMountPath
	} else {
		mountPath = config.VeraCryptMountPath + ":"
	}
	config.RepoPath = path.Join(mountPath, config.RepoPath)

	return config, nil
}

func (c *EncConfig) ExtractSSHFSParams() SSHFSParams {
	return SSHFSParams{
		SSHPort:         c.SSHPort,
		SSHAddress:      c.SSHAddress,
		SSHUser:         c.SSHUser,
		SSHIdentityFile: c.SSHIdentityFile,
		SSHMountPath:    c.SSHMountPath,
		SSHRemotePath:   c.SSHRemotePath,
	}
}

func (c *EncConfig) ExtractVeraCryptParams() VeraCryptParams {
	return VeraCryptParams{
		VeraCryptMountPath: c.VeraCryptMountPath,
		VeraCryptVaultPath: c.VeraCryptVaultPath,
	}
}
