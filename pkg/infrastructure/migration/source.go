// @AI_GENERATED
package migration

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// MigrationSource abstracts the origin of migration files.
type MigrationSource interface {
	GetSourceURL() string
	Validate() error
}

// FilesystemSource reads migrations from a directory on disk.
type FilesystemSource struct {
	path string
}

// NewFilesystemSource creates a new FilesystemSource.
func NewFilesystemSource(path string) *FilesystemSource {
	return &FilesystemSource{path: path}
}

// GetSourceURL returns the file:// URL for golang-migrate.
func (s *FilesystemSource) GetSourceURL() string {
	absPath, err := filepath.Abs(s.path)
	if err != nil {
		return "file://" + s.path
	}
	return "file://" + absPath
}

// Validate checks that the path exists and is a directory.
func (s *FilesystemSource) Validate() error {
	info, err := os.Stat(s.path)
	if err != nil {
		return fmt.Errorf("migration source path %q does not exist: %w", s.path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("migration source path %q is not a directory", s.path)
	}
	return nil
}

// EmbeddedSource reads migrations from an embedded filesystem.
type EmbeddedSource struct {
	fsys   fs.FS
	subdir string
	driver *iofs.PartialDriver
}

// NewEmbeddedSource creates a new EmbeddedSource.
func NewEmbeddedSource(fsys fs.FS, subdir string) *EmbeddedSource {
	return &EmbeddedSource{fsys: fsys, subdir: subdir}
}

// GetSourceURL returns the iofs:// URL for golang-migrate.
func (s *EmbeddedSource) GetSourceURL() string {
	return "iofs://" + s.subdir
}

// Validate checks that the embedded filesystem is accessible.
func (s *EmbeddedSource) Validate() error {
	if s.fsys == nil {
		return fmt.Errorf("embedded migration filesystem is nil")
	}
	entries, err := fs.ReadDir(s.fsys, s.subdir)
	if err != nil {
		return fmt.Errorf("failed to read embedded migration directory %q: %w", s.subdir, err)
	}
	if len(entries) == 0 {
		return fmt.Errorf("embedded migration directory %q contains no files", s.subdir)
	}
	return nil
}

// @AI_GENERATED: end
