package config

import (
	"fmt"
	"os"

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
		return EncConfig{}, fmt.Errorf("can't read config file: %s", err)
	}
	config := EncConfig{}
	err = yaml.Unmarshal(config_raw, &config)
	if err != nil {
		return EncConfig{}, fmt.Errorf("can't parse config file: %s", err)
	}
	if config.SSHPort == 0 {
		config.SSHPort = 22
	}
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
