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
	absolute_path, err := filepath.Abs(local)
	if err == nil {
		return absolute_path
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

// absolutePath try to absolute path to the local git repository
// or CWD
func absolutePath(path string) string {
	if path == "" {
		return path
	}
	if filepath.IsAbs(path) {
		return path
	}
	ref_dir := getLocalRepo()
	if ref_dir == "" {
		ref_dir, _ = os.Getwd()
	}
	p := filepath.Join(ref_dir, path)
	return filepath.Clean(p)
}

func GetConfig() (EncConfig, error) {
	address := os.Getenv("VESH_TEST_ADDRESS")
	if address == "" {
		address = os.Args[2]
	}
	config := parseAddress(address)
	config.LocalRepoPath = getLocalRepo()
	defaultConfig := defaultConfig()

	remote_name := os.Getenv("VESH_TEST_REMOTENAME")
	if remote_name == "" {
		remote_name = os.Args[1]
	}
	config.RemoteName = remote_name

	configPath := os.Getenv("VESH_TEST_CONFIGPATH")
	if configPath == "" {
		configPath = ConfigFileName
	}
	config_raw, err := os.ReadFile(configPath)
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

	// Ensure that mount path for Veracrypt has correct form for windows
	var mountPath string
	if runtime.GOOS == "windows" {
		if strings.HasSuffix(config.VeraCryptMountPath, ":/") {
			mountPath = config.VeraCryptMountPath
		} else {
			mountPath = string(config.VeraCryptMountPath[0]) + ":/"
		}
	}
	config.RepoPath = filepath.Join(mountPath, config.RepoPath)

	// Absolutise paths
	config.LocalRepoPath = absolutePath(config.LocalRepoPath)
	config.SSHIdentityFile = absolutePath(config.SSHIdentityFile)
	if !filepath.IsAbs(config.VeraCryptVaultPath) {
		config.VeraCryptVaultPath = filepath.Join(
			config.SSHMountPath+string(filepath.Separator),
			config.VeraCryptVaultPath)
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
