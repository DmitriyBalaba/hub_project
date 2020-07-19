package filesys

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Manager struct {
	config *Config
}

var required = make(dirsFiles)

func SetRequiredDirs(dirKeys ...string) {
	for _, k := range dirKeys {
		required.addDir(k)
	}
}

func SetRequiredStaticFiles(dirKey string, fileKeys ...string) {
	required.addDir(dirKey)
	for _, fileKey := range fileKeys {
		if _, ok := required[dirKey][fileKey]; ok {
			continue
		}
		required[dirKey][fileKey] = struct{}{}
	}
}

func NewManager(c *Config) (m *Manager, err error) {
	if c == nil {
		panic("cannot create file manager: nil config")
	}

	if err := c.checkDirectories(required); err != nil {
		return nil, err
	}

	m = &Manager{
		config: c,
	}

	return
}

func (m *Manager) getDir(key string) *Directory {
	dir, ok := m.config.Directories[key]
	if !ok {
		panic("key not found in directories: " + key)
	}
	return dir
}

func (m *Manager) GetDirPath(key string) string {
	dir, ok := m.config.Directories[key]
	if !ok {
		return ""
	}
	return dir.Path
}

func (m *Manager) GetFilePath(dirKey, fileKey string) string {
	var dir = m.getDir(dirKey)
	return filepath.Join(dir.Path, dir.getFileName(fileKey))
}

func (m *Manager) NewFileIterator(dirKey string) (*FileIterator, error) {
	dir := m.getDir(dirKey)
	return NewFileIterator(dir.Path, dir.FileMasks)
}

func (m *Manager) CreateDirectories() error {
	for _, d := range m.config.Directories {
		err := os.MkdirAll(d.Path, 0755)
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
