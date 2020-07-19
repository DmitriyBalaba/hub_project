package filesys

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

type File struct {
	os.FileInfo
	Content []byte
}

func (m *Manager) ReadStaticFile(dirKey string, fileKey string) ([]byte, error) {
	fileName := m.getDir(dirKey).getFileName(fileKey)
	return m.ReadFile(dirKey, fileName)
}

func (m *Manager) ReadFile(dirKey, fileName string) ([]byte, error) {
	dir := m.getDir(dirKey)
	content, err := ioutil.ReadFile(filepath.Join(dir.Path, fileName))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return content, nil
}

func (m *Manager) DeleteFile(dirKey, fileName string) error {
	return errors.WithStack(os.Remove(path.Join(m.GetDirPath(dirKey), fileName)))
}

func (m *Manager) CreateFile(file []byte, dirKey, filename string) error {
	f, err := os.Create(filepath.Join(m.GetDirPath(dirKey), filename))
	if err != nil {
		return errors.WithStack(err)
	}

	if _, err := f.Write(file); err != nil {
		return errors.WithStack(err)
	}

	return nil
}
