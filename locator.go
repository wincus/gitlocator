package gitlocator

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
)

type OriginType int

const (
	GitDirName = ".git"
	GitRemote  = "origin"
)

const (
	OriginTypeUnknown OriginType = iota
	OriginTypeGitHub
	OriginTypeGitLab
)

type GitLocation interface {
	GetURL() (string, error)
	IsClean() (bool, error)
}

type GitLocal struct {
	repo   *git.Repository
	subdir string
}

func NewGitLocator(wd string) (GitLocation, error) {

	localRepoPath, err := getRepoPath(wd)

	if err != nil {
		return nil, err
	}

	repo, err := git.PlainOpen(localRepoPath)

	if err != nil {
		return nil, err
	}

	return &GitLocal{
		repo:   repo,
		subdir: getSubDirectory(localRepoPath, wd),
	}, nil

}

func (g *GitLocal) GetURL() (string, error) {

	remoteURI, err := g.getRemoteURI()

	if err != nil {
		return "", err
	}

	branch, err := g.getBranch()

	if err != nil {
		return "", err
	}

	repoType := getRepoType(remoteURI)

	switch repoType {
	case OriginTypeGitHub:
		host, org, repo := parseGithubURI(remoteURI)
		return getGitHubURL(host, org, repo, branch, g.subdir), nil
	case OriginTypeGitLab:
		host, group, project := parseGitlabURI(remoteURI)
		return getGitLabURL(host, group, project, branch, g.subdir), nil
	default:
		return "", fmt.Errorf("unsupported remote repository")
	}

}

func (g *GitLocal) IsClean() (bool, error) {

	wt, err := g.repo.Worktree()

	if err != nil {
		return false, err
	}

	status, err := wt.Status()

	if err != nil {
		return false, err
	}

	return status.IsClean(), nil
}

func (g *GitLocal) getBranch() (string, error) {

	head, err := g.repo.Head()

	if err != nil {
		return "", err
	}

	return head.Name().Short(), nil
}

func (g *GitLocal) getRemoteURI() (string, error) {

	remote, err := g.repo.Remote(GitRemote)

	if err != nil {
		return "", err
	}

	if remote == nil || len(remote.Config().URLs) == 0 {
		return "", fmt.Errorf("remote repository not found")
	}

	return remote.Config().URLs[0], nil
}

func getRepoType(uri string) OriginType {

	if strings.Contains(uri, "github.com") {
		return OriginTypeGitHub
	}

	if strings.Contains(uri, "gitlab.com") {
		return OriginTypeGitLab
	}

	if strings.Contains(uri, "gitlab-ssh") {
		return OriginTypeGitLab
	}

	return OriginTypeUnknown
}

// getRepoPath looks for a git repository in the current working directory
// or any of its parent directories.
// It returns the path to the repository root if found, or an error if not.
func getRepoPath(cwd string) (string, error) {

	path := cwd

	for {

		info, err := os.Stat(fmt.Sprintf("%s/%s", path, GitDirName))

		if err == nil && info.IsDir() {
			return path, nil
		}

		path = filepath.Dir(path)

		if path == "/" {
			return "", fmt.Errorf("git repository not found")
		}
	}
}

// getSubDirectory returns the subdirectory path relative to the root directory.
func getSubDirectory(root, path string) string {
	return strings.TrimPrefix(path[len(root):], "/")
}

// parseGithubURI extracts the host, organization and repository name from a github
// repository URI.
func parseGithubURI(uri string) (string, string, string) {
	regex := regexp.MustCompile(`(?:git@|https://)(?P<host>[^:/]+)[:/](?P<org>[^/]+)/(?P<repo>[^.]+)(?:\.git)?`)
	m := regex.FindStringSubmatch(uri)
	return m[1], m[2], m[3]
}

// getGitHubURL returns the URL to a file in a GitHub repository.
func getGitHubURL(host, org, repo, branch, path string) string {
	return fmt.Sprintf("https://%s/%s/%s/tree/%s/%s", host, org, curateRepo(repo), branch, path)
}

// parseGitlabURI extracts the host, group and project name from a gitlab repository URI.
func parseGitlabURI(uri string) (string, string, string) {
	regex := regexp.MustCompile(`(?:git@|https://)(?P<host>[^:/]+)[:/](?P<org>[^/]+)/(?P<repo>[^.]+)(?:\.git)?`)
	m := regex.FindStringSubmatch(uri)
	return m[1], m[2], m[3]
}

// getGitLabURL returns the URL to a file in a GitLab repository.
func getGitLabURL(host, group, project, branch, path string) string {
	// replace -ssh with . in host
	host = strings.Replace(host, "-ssh.", ".", 1)
	return fmt.Sprintf("https://%s/%s/%s/-/tree/%s/%s", host, group, project, branch, path)
}

// curate repo
func curateRepo(repo string) string {
	return strings.TrimSuffix(repo, ".git")
}
