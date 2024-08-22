package gitcaddy

// state we track per repo, between each step
type Context struct {
	RepoExists bool
	LocalHash  string
	RemoteHash string
}
