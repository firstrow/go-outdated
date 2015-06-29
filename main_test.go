package main

import (
	"testing"
)

func TestGithubExtractUsernameAndRepository(t *testing.T) {
	imp := "github.com/username/repo/path"
	username, repo, err := extractUsernameAndRepository(imp)
	if err != nil {
		t.Error("Shoud return username and repo but returned error")
	}
	if username != "username" {
		t.Error("Shoud return username but returned", username)
	}
	if repo != "repo" {
		t.Error("Shoud return repo but returned", repo)
	}
}

func TestIsGitgubPackage(t *testing.T) {
	if isGithubPackage("github.com/foo/bar") == false {
		t.Error("Should be true")
	}
	if isGithubPackage("foohub.com/foo/bar") == true {
		t.Error("Should be false")
	}
}

func TestNewGithubPackage(t *testing.T) {
	pkg, err := newGithubPackage("github.com/foo/bar", "")
	if err != nil {
		t.Fatal("Should not return error")
	}
	if pkg.username != "foo" {
		t.Error("Expected foo, got", pkg.username)
	}
	if pkg.repository != "bar" {
		t.Error("Expected bar, got", pkg.repository)
	}
}

func TestAddUniquePackage(t *testing.T) {
	var packages = make(packagesList, 0)
	pkg1, _ := newGithubPackage("github.com/foo1/bar1", "")
	pkg2, _ := newGithubPackage("github.com/foo2/bar2", "")
	pkg3, _ := newGithubPackage("github.com/foo2/bar2", "")
	packages.Add(pkg1)
	packages.Add(pkg2)
	packages.Add(pkg3)
	if packages.Count() != 2 {
		t.Error("Expected 2 packages, got", packages.Count())
	}
}
