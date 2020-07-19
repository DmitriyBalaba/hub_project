package env

import (
	"hub_project/models"
	"regexp"
	"strconv"

	"bitbucket.org/jazzserve/webapi/resources/filesys"
	"bitbucket.org/jazzserve/webapi/utils/pointers/ptr"
	"github.com/pkg/errors"
)

var (
	migrationRegex         = regexp.MustCompile("^(.+)_(.+)$")
	migrationNotMatchedErr = errors.Errorf("script doesn't matches: \"%s\"", migrationRegex.String())
	emptyMigrationQueryErr = errors.New("migration query is empty")
)

func composeMigration(f *filesys.File) (*models.Migration, error) {
	if f == nil {
		return nil, errors.New("nil file")
	}

	if len(f.Content) == 0 {
		return nil, emptyMigrationQueryErr
	}

	parsedName := migrationRegex.FindStringSubmatch(f.Name())
	if len(parsedName) == 0 {
		return nil, migrationNotMatchedErr
	}

	version, err := strconv.Atoi(parsedName[1])
	if err != nil {
		return nil, errors.Wrap(err, "invalid version")
	}

	m := &models.Migration{
		ID:    ptr.Int64(int64(version)),
		Name:  &parsedName[2],
		Query: ptr.Str(string(f.Content)),
	}

	return m, nil
}

func (e *env) ReadMigrationScripts() (ms []*models.Migration, err error) {
	if e.fileManager == nil {
		return nil, NotBuiltProperlyErr
	}

	iterator, err := e.fileManager.NewFileIterator(MigrationDir)
	if err != nil {
		return
	}

	for iterator.HasNext() {
		file, err := iterator.ReadNext()
		if err != nil {
			return nil, err
		}

		migration, err := composeMigration(file)
		if err != nil {
			return nil, err
		}

		ms = append(ms, migration)
	}

	return
}
