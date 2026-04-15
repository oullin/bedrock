package filesystem

import "github.com/bedrock/packages/container"

// FilesystemServiceProvider registers the filesystem into the container.
// It mirrors Framework\Filesystem\FilesystemServiceProvider.
type FilesystemServiceProvider struct {
	app *container.Container
}

// NewFilesystemServiceProvider constructs the provider.
func NewFilesystemServiceProvider(app *container.Container) *FilesystemServiceProvider {
	return &FilesystemServiceProvider{app: app}
}

// Register binds the filesystem as a singleton under "files".
func (p *FilesystemServiceProvider) Register() {
	p.app.Singleton("files", func(_ *container.Container) (any, error) {
		return New(), nil
	})
}

// Provides returns the abstract keys registered by this provider.
func (p *FilesystemServiceProvider) Provides() []string {
	return []string{"files"}
}
