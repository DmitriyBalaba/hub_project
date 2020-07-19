package gorm

import (
	"bitbucket.org/jazzserve/webapi/storage"
	"bitbucket.org/jazzserve/webapi/utils"
	"fmt"
	"strings"
)

const (
	OrderDesc = "DESC"
	OrderAsc  = "ASC"
)

func (s *Storage) SortAsc(field string) *Storage {
	return s.Sort(field, OrderAsc)
}

func (s *Storage) SortDesc(field string) *Storage {
	return s.Sort(field, OrderDesc)
}

func (s *Storage) Sort(field, order string) *Storage {
	field, order = unifyStr(field), strings.ToUpper(strings.TrimSpace(order))
	if !utils.IsStrIn(order, OrderAsc, OrderDesc) {
		return s
	}

	query := fmt.Sprintf("%s %s NULLS LAST", field, order)

	return Wrap(s.Order(query))
}

func (s *Storage) MultiSorting(items []storage.SortItem) *Storage {
	for _, i := range items {
		s = s.Sort(i.Field(), i.Order())
	}
	return s
}
