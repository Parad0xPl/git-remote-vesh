package config

import (
	"reflect"
	"testing"
)

func TestParserAll(t *testing.T) {
	input := "user@domain.com:repo.git"
	c := parseAddress(input)
	expected := EncConfig{
		SSHUser:    "user",
		SSHAddress: "domain.com",
		RepoPath:   "repo.git",
	}
	if !reflect.DeepEqual(c, expected) {
		t.Fatalf("Parsed value doesn't match\n---Value---\n%v\n---Expected---\n%v",
			c, expected)
	}
}

func TestParserAddPlusRepo(t *testing.T) {
	input := "domain.com:repo.git"
	c := parseAddress(input)
	expected := EncConfig{
		SSHAddress: "domain.com",
		RepoPath:   "repo.git",
	}
	if !reflect.DeepEqual(c, expected) {
		t.Fatalf("Parsed value doesn't match\n---Value---\n%v\n---Expected---\n%v",
			c, expected)
	}
}

func TestParserRepo(t *testing.T) {
	input := "repo.git"
	c := parseAddress(input)
	expected := EncConfig{
		RepoPath: "repo.git",
	}
	if !reflect.DeepEqual(c, expected) {
		t.Fatalf("Parsed value doesn't match\n---Value---\n%v\n---Expected---\n%v",
			c, expected)
	}
}
