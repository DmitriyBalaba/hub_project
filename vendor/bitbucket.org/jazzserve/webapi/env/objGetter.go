package env

import (
	"bitbucket.org/jazzserve/webapi/sessions/redis/store"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

type ObjGetter = store.ObjGetter

type DefaultObjGetter struct{}

func (g *DefaultObjGetter) GetObjId(sessionId string) (interface{}, error) {
	return GetUserIdFromSessionId(sessionId)
}

func GetUserIdFromSessionId(sessionId string) (int64, error) {
	splitSessionId := strings.Split(sessionId, ":")
	if len(splitSessionId) != 2 {
		return 0, errors.New("invalid form of session id")
	}
	userId, err := strconv.ParseInt(splitSessionId[0], 10, 64)
	if err != nil {
		return 0, errors.Wrap(err, "account id is invalid")
	}
	return userId, nil
}
