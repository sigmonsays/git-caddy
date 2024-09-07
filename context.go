package gitcaddy

// state we track per repo, between each step
type Context struct {
	RepoExists bool

	// resolved working directory (absolute path)
	ResolvedWorkingDir string

	// the remote branch name, as if you just cloned the repo
	UpstreamBranchName string

	LocalHash  string
	RemoteHash string
}
