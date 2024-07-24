package gitlocator

import "testing"

func TestParseGithubURI(t *testing.T) {

	type output struct {
		host string
		org  string
		repo string
	}

	tests := []struct {
		uri  string
		want output
	}{
		{
			uri: "git@github.com:org-dev/infra.git",
			want: output{
				host: "github.com",
				org:  "org-dev",
				repo: "infra",
			},
		},
		{
			uri: "https://github.com/redpanda-data/benthos.git",
			want: output{
				host: "github.com",
				org:  "redpanda-data",
				repo: "benthos",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			host, org, repo := parseGithubURI(tt.uri)
			if host != tt.want.host || org != tt.want.org || repo != tt.want.repo {
				t.Errorf("got: %s, %s, %s, want: %s, %s, %s", host, org, repo, tt.want.host, tt.want.org, tt.want.repo)
			}
		})
	}
}

func TestParseGitlabURI(t *testing.T) {

	type output struct {
		host    string
		group   string
		project string
	}

	tests := []struct {
		uri  string
		want output
	}{
		{
			uri: "git@gitlab-ssh.tools.devrtb.com:devops/iac.git",
			want: output{
				host:    "gitlab.tools.devrtb.com",
				group:   "devops",
				project: "iac",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.uri, func(t *testing.T) {
			host, group, project := parseGithubURI(tt.uri)
			if host != tt.want.host || group != tt.want.group || project != tt.want.project {
				t.Errorf("got: %s, %s, %s, want: %s, %s, %s", host, group, project, tt.want.host, tt.want.group, tt.want.project)
			}
		})
	}
}
