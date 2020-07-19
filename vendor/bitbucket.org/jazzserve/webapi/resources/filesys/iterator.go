package filesys

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type FileIterator struct {
	dirPath     string
	unreadFiles []os.FileInfo
}

func NewFileIterator(dirPath string, masks masks) (*FileIterator, error) {
	list, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, errors.Wrap(err, "can't init file iterator")
	}

	iter := &FileIterator{
		dirPath: dirPath,
	}

	for i := range list {
		matchesFileMask, err := masks.match(list[i].Name())
		if err != nil {
			return nil, errors.WithStack(err)
		}
		if !list[i].IsDir() && matchesFileMask {
			iter.unreadFiles = append(iter.unreadFiles, list[i])
		}
	}

	return iter, nil
}

func (iter *FileIterator) HasNext() bool {
	if len(iter.unreadFiles) == 0 {
		return false
	}
	return true
}

func (iter *FileIterator) ReadNext() (*File, error) {
	if len(iter.unreadFiles) == 0 {
		return nil, nil
	}

	current := iter.unreadFiles[0]
	if len(iter.unreadFiles) == 1 {
		// rm ptr to inner array so GC can free its memory
		iter.unreadFiles = nil
	} else {
		// or reslice if there're more elements
		iter.unreadFiles = iter.unreadFiles[1:]
	}

	path := filepath.Join(iter.dirPath, current.Name())
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	file := &File{
		FileInfo: current,
		Content:  content,
	}

	return file, nil
}
