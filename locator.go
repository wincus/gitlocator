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

	return OriginTypeUnknown
}

// func NewGitLocator(wd string) (GitLocation, error) {

// 	localRepoPath, err := getRepoPath(wd)

// 	if err != nil {
// 		return nil, err
// 	}

// 	repo, err := git.PlainOpen(localRepoPath)

// 	if err != nil {
// 		return nil, err
// 	}

// 	remote, err := repo.Remote(GitRemote)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if remote == nil || len(remote.Config().URLs) == 0 {
// 		return nil, fmt.Errorf("remote repository not found")
// 	}

// 	remoteURL := remote.Config().URLs[0]

// 	// if !matchesScheme(remoteURL) {
// 	// 	return nil, fmt.Errorf("unsupported remote URL scheme")
// 	// }

// 	if !matchesScpLike(remoteURL) {
// 		return nil, fmt.Errorf("unsupported remote URL format")
// 	}

// 	_, host, port, remotePath := findScpLikeComponents(remote.Config().URLs[0])

// 	head, err := repo.Head()

// 	if err != nil {
// 		return nil, err
// 	}

// 	branch := head.Name().Short()

// 	// check repo status in case is dirty
// 	wt, err := repo.Worktree()

// 	if err != nil {
// 		return nil, err
// 	}

// 	status, err := wt.Status()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &GitLocation{
// 		isClean:          status.IsClean(),
// 		remote:           remote,
// 		scheme:           "https",
// 		branch:           branch,
// 		localRepoSubPath: getSubDirectory(localRepoPath, wd),
// 		localRepoPath:    localRepoPath,
// 		host:             host,
// 		port:             port,
// 		remotePath:       remotePath,
// 	}, nil
// }

// func NewGitLocal(wd string) (*GitLocal, error) {

// 	localRepoPath, err := getRepoPath(wd)

// 	if err != nil {
// 		return nil, err
// 	}

// 	repo, err := git.PlainOpen(localRepoPath)

// 	if err != nil {
// 		return nil, err
// 	}

// 	remote, err := repo.Remote(GitRemote)

// 	if err != nil {
// 		return nil, err
// 	}

// 	if remote == nil || len(remote.Config().URLs) == 0 {
// 		return nil, fmt.Errorf("remote repository not found")
// 	}

// 	remoteURL := remote.Config().URLs[0]

// 	// if !matchesScheme(remoteURL) {
// 	// 	return nil, fmt.Errorf("unsupported remote URL scheme")
// 	// }

// 	if !matchesScpLike(remoteURL) {
// 		return nil, fmt.Errorf("unsupported remote URL format")
// 	}

// 	_, _, _, remotePath := findScpLikeComponents(remote.Config().URLs[0])

// 	head, err := repo.Head()

// 	if err != nil {
// 		return nil, err
// 	}

// 	branch := head.Name().Short()

// 	// check repo status in case is dirty
// 	wt, err := repo.Worktree()

// 	if err != nil {
// 		return nil, err
// 	}

// 	status, err := wt.Status()

// 	if err != nil {
// 		return nil, err
// 	}

// 	return &GitLocal{
// 		origin:  OriginTypeUnknown,
// 		wd:      wd,

// func (g *GitLocation) String() string {
// 	return fmt.Sprintf("GitLocation{scheme: %s, host: %s, port: %s, org: %s, branch: %s, repoPath: %s, path: %s, isClean: %t}", g.scheme, g.host, g.port, g.org, g.branch, g.localRepoPath, g.localRepoSubPath, g.isClean)
// }

// func (g *GitLocation) GetURL() string {
// 	return fmt.Sprintf("%s://%s/%s/tree%s/%s", g.scheme, g.host, g.org, g.branch, g.localRepoSubPath)
// }

// func (g *GitLocation) IsClean() bool {
// 	return g.isClean
// }

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

func parseGithubURI(uri string) (string, string, string) {

	// needs to be fixed
	scpLikeUrlRegExp := regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5}):)?(?P<path>[^\\].*)$`)
	m := scpLikeUrlRegExp.FindStringSubmatch(uri)
	return m[1], m[2], m[3]
}

func getGitHubURL(host, org, repo, branch, path string) string {
	return fmt.Sprintf("https://%s/%s/%s/tree/%s/%s", host, org, repo, branch, path)
}

func parseGitlabURI(uri string) (string, string, string) {

	// needs to be fixed
	scpLikeUrlRegExp := regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5}):)?(?P<path>[^\\].*)$`)
	m := scpLikeUrlRegExp.FindStringSubmatch(uri)

	return m[1], m[2], m[3]
}

func getGitLabURL(host, group, project, branch, path string) string {
	return fmt.Sprintf("https://%s/%s/%s/-/tree/%s/%s", host, group, project, branch, path)
}

// // copied from go-git internal implementation

// var (
// 	isSchemeRegExp = regexp.MustCompile(`^[^:]+://`)
// 	// Ref: https://github.com/git/git/blob/master/Documentation/urls.txt#L37
// 	scpLikeUrlRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5}):)?(?P<path>[^\\].*)$`)
// )

// // matchesScheme returns true if the given string matches a URL-like
// // format scheme.
// func matchesScheme(url string) bool {
// 	return isSchemeRegExp.MatchString(url)
// }

// // matchesScpLike returns true if the given string matches an SCP-like
// // format scheme.
// func matchesScpLike(url string) bool {
// 	return scpLikeUrlRegExp.MatchString(url)
// }

// // findScpLikeComponents returns the user, host, port and path of the
// // given SCP-like URL.
// func findScpLikeComponents(url string) (user, host, port, path string) {
// 	m := scpLikeUrlRegExp.FindStringSubmatch(url)
// 	return m[1], m[2], m[3], m[4]
// }

// func getOriginType(g *git.Repository) OriginType {

// 	remote, err := g.Remote(GitRemote)

// 	if err != nil {
// 		return OriginTypeUnknown
// 	}

// 	if len(remote.Config().URLs) == 0 {
// 		return OriginTypeUnknown
// 	}

// 	remoteURL := remote.Config().URLs[0]

// 	if strings.Contains(remoteURL, "github.com") {
// 		return OriginTypeGitHub
// 	}

// 	if strings.Contains(remoteURL, "gitlab.com") {
// 		return OriginTypeGitLab
// 	}

// 	return OriginTypeUnknown
// }
