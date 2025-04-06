package config

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	"github.com/Parad0xpl/git-remote-vesh/v2/debug"
	"gopkg.in/yaml.v2"
)

// VeshConfig holds all data related to connection and encryption. It is based
// on the remote URL and the configuration inside the repository for more customization.
// Most of it should be left to default for easier usage, as all important data can be
// provided in the URL.
// TODO: Add user configuration in the home directory
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

// SSHFSParams contains configuration related only to SSHFS mounting.
type SSHFSParams struct {
	SSHPort         int
	SSHAddress      string
	SSHUser         string
	SSHIdentityFile string
	SSHMountPath    string
	SSHRemotePath   string
}

// VeraCryptParams contains configuration related only to VeraCrypt mounting.
type VeraCryptParams struct {
	VeraCryptMountPath string
	VeraCryptVaultPath string
}

// ConfigFileName is the default filename of the configuration file inside the repository.
const ConfigFileName = ".veshconfig"

// getLocalRepo returns the absolute path to the repository held inside
// the GIT_DIR environment variable.
func getLocalRepo() string {
	local := os.Getenv("GIT_DIR")
	if local == "" {
		log.Println("WARNING: No local dir path")
	}
	absolutePath, err := filepath.Abs(local)
	if err == nil {
		return absolutePath
	}
	return local
}

// GetVeshConfigDir returns the "inside repo" directory for Vesh, used for mark files
// and other local configurations.
func (c *VeshConfig) GetVeshConfigDir() string {
	p := filepath.Join(c.LocalRepoPath, "vesh", c.RemoteName)
	debug.Println("Vesh path:", p)

	return p
}

// absolutePath converts a given path to an absolute path relative to the local
// git repository or the current working directory if no repository path is found.
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

// mergeConfig merges one configuration into another, replacing only empty fields.
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

// readConfig reads the given file and parses it as a VeshConfig.
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

// testConfig returns a configuration based on test variables:
// VESH_TEST_ADDRESS, VESH_TEST_REMOTENAME, and VESH_TEST_CONFIGPATH.
func testConfig() VeshConfig {
	config := parseAddress(os.Getenv("VESH_TEST_ADDRESS"))
	config.RemoteName = os.Getenv("VESH_TEST_REMOTENAME")
	configPath := os.Getenv("VESH_TEST_CONFIGPATH")
	fileConfig, _ := readConfig(configPath)
	return mergeConfig(fileConfig, config)
}

// argumentConfig returns a configuration parsed from command-line arguments.
func argumentConfig() VeshConfig {
	config := parseAddress(os.Args[2])
	config.RemoteName = os.Args[1]
	return config

}

// homeConfig returns a configuration optionally located in the user's home directory.
// Useful if SSH has a custom unencrypted key for easier access.
func homeConfig() VeshConfig {
	homePath, err := os.UserHomeDir()
	if err != nil {
		debug.Println("Can't get home path: ", err)
		return VeshConfig{}
	}

	config, _ := readConfig(filepath.Join(homePath, ConfigFileName))
	return config
}

// localConfig returns a configuration located in the current working directory.
func localConfig() VeshConfig {
	config, _ := readConfig(ConfigFileName)
	return config
}

// GetConfig returns a parsed configuration from all sources. It is based on defaultConfig()
// and filled with custom populated properties.
//
// For testing, some options can be passed with environment variables:
// VESH_TEST_ADDRESS, VESH_TEST_REMOTENAME, and VESH_TEST_CONFIGPATH.
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
	} else {
		mountPath = config.VeraCryptMountPath
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

// ExtractSSHFSParams extracts SSHFS-related properties from the configuration.
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

// ExtractVeraCryptParams extracts VeraCrypt-related properties from the configuration.
func (c *VeshConfig) ExtractVeraCryptParams() VeraCryptParams {
	return VeraCryptParams{
		VeraCryptMountPath: c.VeraCryptMountPath,
		VeraCryptVaultPath: c.VeraCryptVaultPath,
	}
}
