package version

// version is a private field and should be set when compiling with --ldflags="-X github.com/innovia/secrets-consumer-env/pkg/version.version=X.Y.Z"
var version = "0.0.0-unset"

// gitCommitID is a private field and should be set when compiling with --ldflags="-X github.com/innovia/secrets-consumer-env/pkg/version.gitCommitID=<commit-id>"
var gitCommitID = ""

// GetVersion returns the current version
func GetVersion() string {
	return version
}

// GetGitCommitID returns the git commit id from which it is being built
func GetGitCommitID() string {
	return gitCommitID
}
