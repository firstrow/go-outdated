package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Unknwon/com"
	"github.com/mgutz/ansi"
)

type branchInfo struct {
	Commit struct {
		Commit struct {
			Commiter struct {
				Date string `json:"date"`
			} `json:"committer"`
		} `json:"commit"`
	} `json:"commit"`
}

type githubPackage struct {
	ImportName string
	parent     string
	username   string // Github username
	repository string // Github repository
}

type packagesList map[string]*githubPackage

func (p packagesList) Add(pkg *githubPackage) {
	if _, ok := p[pkg.ImportName]; !ok {
		p[pkg.ImportName] = pkg
	}
}

func (p packagesList) Count() int {
	return len(p)
}

// Exclude package by import path pattern
func (p packagesList) Exclude(pattern string) packagesList {
	var packages = make(packagesList, 0)
	for _, pkg := range p {
		if !strings.Contains(pkg.ImportName, pattern) {
			packages.Add(pkg)
		}
	}
	return packages
}

// Get last commit date from Github
func (g *githubPackage) GithubLastCommitDate() (time.Time, error) {
	branch := "master"
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/branches/%s?access_token=%s",
		g.username, g.repository, branch, accessToken)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Time-Zone", "UTC")
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		return time.Time{}, errors.New("Github request error.")
	}

	if resp.StatusCode != 200 {
		return time.Time{}, errors.New("Error. Page not found.")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return time.Time{}, errors.New("Error reading response body")
	}

	branchInfo := &branchInfo{}
	err = json.Unmarshal(body, branchInfo)
	if err != nil {
		return time.Time{}, errors.New("Error parse json")
	}

	t, err := time.Parse(time.RFC3339, branchInfo.Commit.Commit.Commiter.Date)
	if err != nil {
		return time.Time{}, errors.New("Error parse github response date")
	}
	return t, nil
}

// Get last commit date from local repository
func (g *githubPackage) LocalLastCommitDate() (time.Time, error) {
	stdout, _, _ := com.ExecCmdDir(g.Dir(), "git", "log", "-1", "--date=rfc2822", "--pretty=format:%cd")
	layout := "Mon, _2 Jan 2006 15:04:05 -0700"
	t, err := time.Parse(layout, stdout)
	if err != nil {
		return time.Time{}, errors.New("Error parse local time")
	}
	return t.UTC(), nil
}

func (g *githubPackage) DisplayResult() {
	localTime, lerr := g.LocalLastCommitDate()
	githubTime, gerr := g.GithubLastCommitDate()

	fmt.Println("Package:", g.ImportName)
	if lerr != nil {
		fmt.Println("Local:  Error")
	} else {
		fmt.Println("Local:  ", localTime)
	}

	if gerr != nil {
		fmt.Println("Github:  Error")
	} else {
		fmt.Println("Github: ", githubTime)
	}

	if lerr == nil && gerr == nil {
		delta := localTime.Sub(githubTime)
		if delta.Minutes() == 0 {
			green := ansi.ColorCode("green")
			reset := ansi.ColorCode("reset")
			fmt.Println("Status:", green, "OK", reset)
		} else {
			green := ansi.ColorCode("red")
			reset := ansi.ColorCode("reset")
			fmt.Println("Status:", green, "Outdated", reset)
		}
	}
	fmt.Println(strings.Repeat("-", len(g.ImportName)))
}

// Get full path to package directory
func (g *githubPackage) Dir() string {
	p, err := build.Import(g.ImportName, "", 0)
	if err != nil {
		return ""
	}
	return p.Dir
}

// Get parent package import name
func (g *githubPackage) Parent() string {
	return g.parent
}

func newGithubPackage(imp, parent string) (*githubPackage, error) {
	username, repository, err := extractUsernameAndRepository(imp)
	if err != nil {
		return nil, err
	}
	return &githubPackage{imp, parent, username, repository}, nil
}

// Find all imported packages
func findImports(packages packagesList, name string) {
	pkg, err := build.Import(name, "", build.AllowBinary)
	if err != nil {
		log.Println("Package import error:", err)
		return
	}

	for _, imp := range pkg.Imports {
		if isGithubPackage(imp) {
			gpkg, _ := newGithubPackage(imp, name)
			packages.Add(gpkg)
			findImports(packages, imp)
		}
	}
}

// Check if package hosted on Github
func isGithubPackage(pkg string) bool {
	return strings.HasPrefix(pkg, "github.com")
}

// extractUsernameAndRepository extracts username and repository from
// package name. e.g: github.com/foo/bar
// Returns: username, repository, error
func extractUsernameAndRepository(imp string) (string, string, error) {
	if !isGithubPackage(imp) {
		return "", "", errors.New("Package name should have github.com prefix")
	}
	parts := strings.Split(imp, "/")
	if len(parts) >= 3 {
		return parts[1], parts[2], nil
	}
	return "", "", errors.New("Error extracting github username and repository")
}

func getGihubTokenFromConfig() string {
	stdout, _, _ := com.ExecCmd("git", "config", "--global", "github.token")
	return strings.TrimSpace(stdout)
}

var accessToken string

func setAccessToken() {
	flag.StringVar(&accessToken, "token", "", "GitHub Access Token")
	flag.Parse()

	// If token not present, try to use token from git config
	if accessToken == "" {
		accessToken = getGihubTokenFromConfig()
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("CWD error:", err)
	}
	pkg, err := build.ImportDir(cwd, build.AllowBinary)
	if err != nil {
		log.Fatalln("Error importing current package:", err)
	}

	setAccessToken()

	var packages = make(packagesList, 0)
	findImports(packages, pkg.ImportPath)

	external := packages.Exclude(pkg.ImportPath)
	msg := fmt.Sprintf("Total packages found: %d", external.Count())
	fmt.Println(msg)
	fmt.Println(strings.Repeat("-", len(msg)))
	for _, pkg := range external {
		pkg.DisplayResult()
	}
}
