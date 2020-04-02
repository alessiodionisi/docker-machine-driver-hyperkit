package hyperkit

// The current version of the docker-machine-driver-hyperkit

// version is a private field and should be set when compiling with --ldflags="-X k8s.io/minikube/pkg/drivers/hyperkit.version=vX.Y.Z"
var version = "v0.0.0-unset"

// gitCommitID is a private field and should be set when compiling with --ldflags="-X k8s.io/minikube/pkg/drivers/hyperkit.gitCommitID=<commit-id>"
var gitCommitID = ""

// GetVersion returns the current docker-machine-driver-hyperkit version
func GetVersion() string {
	return version
}

// GetGitCommitID returns the git commit id from which it is being built
func GetGitCommitID() string {
	return gitCommitID
}
