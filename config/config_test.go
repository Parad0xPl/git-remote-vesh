package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigAbsolutePaths(t *testing.T) {
	os.Setenv("VESH_TEST_CONFIGPATH", "test/config")
	os.Setenv("VESH_TEST_REMOTENAME", "test.git")
	os.Setenv("VESH_TEST_ADDRESS", "bit_revolver@parad0x.pl:test.git")
	os.Setenv("GIT_DIR", "../test_repo")
	config, err := GetConfig()
	if err != nil {
		t.Errorf("Can't get config: %e", err)
	}

	// fmt.Printf("%+v\n", config)
	//
	// fmt.Println("Repo path:", config.RepoPath)
	// fmt.Println("LocalRepoPath:", config.LocalRepoPath)
	to_test := []string{config.SSHIdentityFile, config.VeraCryptVaultPath,
		config.RepoPath, config.LocalRepoPath}
	for _, v := range to_test {
		if v != "" && !filepath.IsAbs(v) {
			t.Errorf("One of path is not absolute - %s", v)
		}
	}
}

func TestMerge(t *testing.T) {
	baseConfig := VeshConfig{
		SSHPort: 1,
		SSHUser: "test",
	}

	defaults := VeshConfig{
		SSHAddress: "address",
		SSHPort:    22,
	}

	output := mergeConfig(baseConfig, defaults)

	if output.SSHPort != 1 {
		t.Errorf("Setted property SSHPort overwritten")
	}
	if output.SSHUser != "test" {
		t.Errorf("Setted property SSHUser overwritten")
	}
	if output.SSHAddress != "address" {
		t.Errorf("Unsetted property SSHAddress is not set")
	}

}
