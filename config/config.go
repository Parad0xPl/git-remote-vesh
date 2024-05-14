package config

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/debug"
	"github.com/Parad0xpl/git-remote-vesh/v2/utils"
	"gopkg.in/yaml.v2"
)

// VeshConfig holds all data related to connectio and encryption. It is based
// on remote url and config inside repo for more customization. Most of it
// should be left to default for easier usage as all important data can be
// provided in url.
// TODO: Add user config in home directory
type VeshConfig struct {
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

// SSHFSParams is extracted part of config related only to SSHFS mounting
type SSHFSParams struct {
	SSHPort         int
	SSHAddress      string
	SSHUser         string
	SSHIdentityFile string
	SSHMountPath    string
	SSHRemotePath   string
}

// VeraCryptParams is extracted part of config related only to VeraCrypt
// mounting
type VeraCryptParams struct {
	VeraCryptMountPath string
	VeraCryptVaultPath string
}

// ConfigFileName is a default filename of the config inside of the repo
const ConfigFileName = ".veshconfig"

// getLocalRepo should return absolute path to the repositorie holded inside
// GIT_DIR environment variable
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

// GetVeshConfigDir return vesh "inside repo" directory for mark files and other
// local config
func (c *VeshConfig) GetVeshConfigDir() string {
	p := filepath.Join(c.LocalRepoPath, "vesh", c.RemoteName)
	if utils.IsDebug() {
		log.Println("Vesh path:", p)
	}

	return p
}

// absolutePath try to absolute path to the local git repository
// or CWD if there is no repository path
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

// Merge one config into another replacing only empty fields
func mergeConfig(to, from VeshConfig) VeshConfig {
	valueTo := reflect.ValueOf(&to)
	valueFrom := reflect.ValueOf(from)

	fields := reflect.VisibleFields(valueFrom.Type())
	numberOfFields := len(fields)

	for i := 0; i < numberOfFields; i++ {
		leftValue := valueTo.Elem().FieldByIndex([]int{i})
		if leftValue.IsZero() {
			rightValue := valueFrom.FieldByIndex([]int{i})
			leftValue.Set(rightValue)
		}
	}

	return to
}

// Read given file and parse it as a VeshConfig
func readConfig(configPath string) (VeshConfig, error) {
	var config VeshConfig
	configRaw, err := os.ReadFile(configPath)
	if err != nil {
		return VeshConfig{}, err
	} else {
		yaml.Unmarshal(configRaw, &config)
	}
	return config, nil
}

// Return config based on test variables:
// VESH_TEST_ADDRESS
// VESH_TEST_REMOTENAME
// VESH_TEST_CONFIGPATH
func testConfig() VeshConfig {
	config := parseAddress(os.Getenv("VESH_TEST_ADDRESS"))
	config.RemoteName = os.Getenv("VESH_TEST_REMOTENAME")
	configPath := os.Getenv("VESH_TEST_CONFIGPATH")
	fileConfig, _ := readConfig(configPath)
	return mergeConfig(fileConfig, config)
}

// Return config parsed from arguments
func argumentConfig() VeshConfig {
	config := parseAddress(os.Args[2])
	config.RemoteName = os.Args[1]
	return config

}

// Return config optionaly putted inside users home directory. Usefull
// if SSH has custom unencrypted key for easier access
func homeConfig() VeshConfig {
	homePath, err := os.UserHomeDir()
	if err != nil {
		debug.Println("Can't get home path: ", err)
		return VeshConfig{}
	}

	config, _ := readConfig(filepath.Join(homePath, ConfigFileName))
	return config
}

// Return config inserted in CWD
func localConfig() VeshConfig {
	config, _ := readConfig(ConfigFileName)
	return config
}

// GetConfig return parsed config from all sources. It is based on defaultConfig()
// and filled with custom populated properties.
//
// For testing some options can be passed with environment variables:
// VESH_TEST_ADDRESS
// VESH_TEST_REMOTENAME
// VESH_TEST_CONFIGPATH
func GetConfig() (VeshConfig, error) {
	var config VeshConfig
	config.LocalRepoPath = getLocalRepo()

	config = mergeConfig(config, testConfig())
	config = mergeConfig(config, argumentConfig())
	config = mergeConfig(config, localConfig())
	config = mergeConfig(config, homeConfig())
	config = mergeConfig(config, defaultConfig())

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

// ExtractSSHFSParams return sshfs isolated properties
func (c *VeshConfig) ExtractSSHFSParams() SSHFSParams {
	return SSHFSParams{
		SSHPort:         c.SSHPort,
		SSHAddress:      c.SSHAddress,
		SSHUser:         c.SSHUser,
		SSHIdentityFile: c.SSHIdentityFile,
		SSHMountPath:    c.SSHMountPath,
		SSHRemotePath:   c.SSHRemotePath,
	}
}

// ExtractVeraCryptParam return VeraCrypt isolated properties
func (c *VeshConfig) ExtractVeraCryptParams() VeraCryptParams {
	return VeraCryptParams{
		VeraCryptMountPath: c.VeraCryptMountPath,
		VeraCryptVaultPath: c.VeraCryptVaultPath,
	}
}
