//go:build linux

package config

func defaultConfig() EncConfig {
	return EncConfig{
		SSHPort:            22,
		SSHMountPath:       "/tmp/vesh_sshfs",
		VeraCryptMountPath: "/tmp/vesh_veracrypt",
		VeraCryptVaultPath: "vesh.crypt",
	}
}
