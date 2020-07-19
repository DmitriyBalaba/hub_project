package payload

import (
	"bitbucket.org/jazzserve/webapi/storage"
	"bitbucket.org/jazzserve/webapi/utils"
	"context"
	"github.com/justinas/alice"
	"net/http"
	"strings"
)

type OffsetLimit struct {
	Offset int64 `schema:"offset"`
	Limit  int64 `schema:"limit"`
}

func (qp *OffsetLimit) GetOffsetLimit() (int64, int64) {
	return qp.Offset, qp.Limit
}

func (qp *OffsetLimit) IsEmpty() bool {
	if qp.Offset == 0 && qp.Limit == 0 {
		return true
	}
	return false
}

func GetOffsetLimit(ctx context.Context) *OffsetLimit {
	return GetQueryParams(ctx).(*OffsetLimit)
}

func GetQueryParams(ctx context.Context) interface{} {
	return getOrPanic(ctx, QueryParamsCtxKey)
}

func QueryParams(obj interface{}) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			obj := getCopyByType(obj)

			if err := r.ParseForm(); err != nil {
				ErrorResponse(w, WrapBadRequestErr(err))
				return
			}

			err := schemaDecoder.Decode(obj, r.Form)
			if err != nil {
				ErrorResponse(w, WrapBadRequestErr(err))
				return
			}

			if err := validate.Struct(obj); err != nil {
				ErrorResponse(w, WrapBadRequestErr(err))
				return
			}

			ctx := context.WithValue(r.Context(), QueryParamsCtxKey, obj)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

type SortItem struct {
	field string
	order string
}

func (item *SortItem) Field() string { return item.field }

func (item *SortItem) SetField(val string) { item.field = val }

func (item *SortItem) Order() string { return item.order }

const (
	multiSortParamSeparator      = ","
	multiSortFieldOrderSeparator = ":"
	orderAsc, orderDesc          = "ASC", "DESC"
)

func ParseMultiSort(rawParameters string, fieldNames ...string) ([]storage.SortItem, error) {
	uniqueNames := make(map[string]struct{})
	for i := range fieldNames {
		uniqueNames[fieldNames[i]] = struct{}{}
	}

	parameters := strings.Split(rawParameters, multiSortParamSeparator)
	sortItems := make([]storage.SortItem, len(parameters))

	for i, rawP := range parameters {
		p := strings.Split(rawP, multiSortFieldOrderSeparator)
		if len(p) != 2 {
			return nil, NewBadRequestErr("%s is invalid for multi sorting", rawP)
		}

		item := &SortItem{
			field: strings.ToLower(strings.TrimSpace(p[0])),
			order: strings.ToUpper(strings.TrimSpace(p[1])),
		}

		if !utils.IsStrIn(item.order, orderAsc, orderDesc) {
			return nil, NewBadRequestErr("invalid order: %s", item.order)
		}

		if _, ok := uniqueNames[item.field]; !ok {
			return nil, NewBadRequestErr("can't sort by field %s: not found", item.field)
		}

		sortItems[i] = item
	}

	return sortItems, nil
}
