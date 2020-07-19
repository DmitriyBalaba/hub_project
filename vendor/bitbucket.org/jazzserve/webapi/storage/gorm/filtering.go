package gorm

import (
	"bitbucket.org/jazzserve/webapi/utils"
	"fmt"
	"strings"
)

func (s *Storage) Match(field string, value interface{}) *Storage {
	if utils.IsZeroAny(field) {
		return s
	}
	return Wrap(s.Where(field+" = ?", value))
}

func (s *Storage) MatchOptional(field string, value interface{}) *Storage {
	if utils.IsZeroAny(value) {
		return s
	}
	return s.Match(field, value)
}

func (s *Storage) MatchNot(field string, value interface{}) *Storage {
	if utils.IsZeroAny(field) {
		return s
	}
	return Wrap(s.Where(field+" <> ?", value))
}

func (s *Storage) MatchNotOptional(field string, value interface{}) *Storage {
	if utils.IsZeroAny(value) {
		return s
	}
	return s.MatchNot(field, value)
}

func (s *Storage) MatchInArr(field string, array interface{}) *Storage {
	if !utils.IsArray(array) {
		return s
	}
	return Wrap(s.Where(field+" IN (?)", array))
}

func (s *Storage) MatchStrNoCase(field, value string) *Storage {
	if utils.IsEmptyStrAny(field, value) {
		return s
	}

	query := fmt.Sprintf("LOWER(%s) = ?", field)
	value = strings.ToLower(value)

	return Wrap(s.Where(query, value))
}

func (s *Storage) StrContains(field, value string) *Storage {
	if utils.IsEmptyStrAny(field, value) {
		return s
	}

	query := fmt.Sprintf("LOWER(%s) LIKE ?", field)
	value = fmt.Sprintf("%%%s%%", strings.ToLower(value))

	return Wrap(s.Where(query, value))
}

func (s *Storage) OffsetAndLimit(offset, limit int64) *Storage {
	defOffset, defLimit := s.GetDefaultOffsetAndLimit()
	if offset == 0 {
		offset = defOffset
	}
	if limit == 0 {
		limit = defLimit
	}

	return Wrap(s.Offset(offset).Limit(limit))
}

func (s *Storage) Exists(field string) *Storage {
	if utils.IsEmptyStrAny(field) {
		return s
	}

	query := fmt.Sprintf("%s NOTNULL", field)
	return Wrap(s.Where(query))
}

func (s *Storage) MatchEmpty(field string) *Storage {
	if utils.IsEmptyStrAny(field) {
		return s
	}

	query := fmt.Sprintf("%s ISNULL", field)
	return Wrap(s.Where(query))
}
