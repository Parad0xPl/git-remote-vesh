//go:build linux

package config

// defaultConfig return default config related to given platform
func defaultConfig() VeshConfig {
	return VeshConfig{
		SSHPort:            22,
		SSHMountPath:       "/tmp/vesh_sshfs",
		VeraCryptMountPath: "/tmp/vesh_veracrypt",
		VeraCryptVaultPath: "vesh.crypt",
	}
}
