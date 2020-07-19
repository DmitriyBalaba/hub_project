package filesys

type Directory struct {
	Path        string            `yaml:"path" validate:"required"`
	FileMasks   masks             `yaml:"file-masks"`
	StaticFiles map[string]string `yaml:"static-files"`
}

func (d *Directory) getFileName(key string) string {
	name, ok := d.StaticFiles[key]
	if !ok {
		panic("key not found in static files: " + key)
	}
	return name
}
