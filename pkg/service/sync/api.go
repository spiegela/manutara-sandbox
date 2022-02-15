package secret

type Synchronizer interface {
	// Name is the name of the synchronizer
	Name() string

	// Start starts a watch for a synchronized resource
	Start()

	// Refresh forces synchronous refresh of local object
	// with object got from kubernetes.
	Refresh()
}
