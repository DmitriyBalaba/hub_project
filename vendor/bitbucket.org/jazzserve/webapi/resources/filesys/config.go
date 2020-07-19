package filesys

import (
	"github.com/pkg/errors"
)

type Config struct {
	Directories map[string]*Directory `yaml:"directories" validate:"dive,required"`
}

var (
	dirsLenLowerErr  = errors.New("len of directories is lower than required")
	filesLenLowerErr = errors.New("len of static files is lower than required")
)

func (c *Config) checkDirectories(requiredKeys dirsFiles) error {
	if len(c.Directories) < len(requiredKeys) {
		return dirsLenLowerErr
	}
	for dirKey, files := range requiredKeys {
		dir, ok := c.Directories[dirKey]
		if !ok {
			return errors.New("dir not found in config: " + dirKey)
		}
		_, err := dir.FileMasks.match("file-mask-test")
		if err != nil {
			return err
		}

		if len(dir.StaticFiles) < len(files) {
			return filesLenLowerErr
		}
		for fileKey := range files {
			if _, ok := dir.StaticFiles[fileKey]; !ok {
				return errors.Errorf("file %d not found in dir %s from config", dirKey, fileKey)
			}
		}
	}
	return nil
}
