package filesys

import (
	"html/template"
	"path/filepath"

	"github.com/pkg/errors"
)

type masks []string

func (m masks) match(filename string) (bool, error) {
	for i := range m {
		ok, err := filepath.Match(m[i], filename)
		if err != nil {
			return false, errors.Wrapf(err, "file-mask %s is invalid", m[i])
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

type dirsFiles map[string]map[string]struct{}

func (df dirsFiles) addDir(dirKey string) {
	if _, ok := df[dirKey]; ok {
		return
	}
	required[dirKey] = make(map[string]struct{})
}

func ParseTemplates(templatesPath string) (*template.Template, error) {
	if templatesPath == "" {
		return nil, errors.New("cannot parse templates: empty path string")
	}
	templates, err := template.ParseGlob(templatesPath + "/*.html")
	if err != nil {
		return nil, err
	}
	return templates, nil
}
