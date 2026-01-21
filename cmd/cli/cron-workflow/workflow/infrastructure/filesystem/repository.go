package filesystem

// Repository abstracts filesystem operations required by workflows so that
// alternative implementations (e.g. in-memory for tests) can be injected.
type Repository interface {
	Write(path string, overwrites bool, content string) error
	EnsureDir(path string) error
	WorkingDir() (string, error)
}
