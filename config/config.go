package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
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

	LocalRepoPath string
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

const ConfigFileName = ".veshconfig"

func getLocalRepo() string {
	local := os.Getenv("GIT_DIR")
	if local == "" {
		log.Println("WARNING: No local dir path")
	}
	return local
}

func (c *EncConfig) GetVeshConfigDir() string {
	p := filepath.Join(c.LocalRepoPath, "vesh")
	if utils.IsDebug() {
		log.Println("Vesh path:", p)
	}

	return p
}

func GetConfig() (EncConfig, error) {
	config := parseAddress(os.Args[2])
	config.LocalRepoPath = getLocalRepo()
	defaultConfig := defaultConfig()
	config.RemoteName = os.Args[1]

	config_raw, err := os.ReadFile(ConfigFileName)
	if err != nil {
	} else {
		yaml.Unmarshal(config_raw, &config)
	}

	// Default
	if config.SSHPort == 0 {
		config.SSHPort = defaultConfig.SSHPort
	}
	if config.SSHMountPath == "" {
		config.SSHMountPath = defaultConfig.SSHMountPath
	}
	if config.VeraCryptMountPath == "" {
		config.VeraCryptMountPath = defaultConfig.VeraCryptMountPath
	}
	if config.VeraCryptVaultPath == "" {
		config.VeraCryptVaultPath = defaultConfig.VeraCryptVaultPath
	}

	var mountPath string
	if runtime.GOOS == "windows" {
		if strings.HasSuffix(config.VeraCryptMountPath, ":/") {
			mountPath = config.VeraCryptMountPath
		} else {
			mountPath = string(config.VeraCryptMountPath[0]) + ":/"
		}
	}
	config.RepoPath = filepath.Join(mountPath, config.RepoPath)

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
