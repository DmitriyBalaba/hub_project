package store

import "github.com/gomodule/redigo/redis"

type ObjGetter interface {
	GetObjId(sessionId string) (objId interface{}, err error)
	GetObj(objId interface{}) (obj interface{}, err error)
}

type valuesCache struct {
	// sessionId - redis values
	values      map[interface{}]map[string]string
	getter      ObjGetter
	interpreter Interpreter
}

func (c *valuesCache) get(sessionId string) (map[string]string, error) {
	objId, err := c.getter.GetObjId(sessionId)
	if err != nil {
		return nil, err
	}

	redisObj, ok := c.values[objId]
	if ok {
		return redisObj, nil
	}

	obj, err := c.getter.GetObj(objId)
	if err != nil {
		return nil, err
	}

	redisObj, err = c.interpreter.ConvertObjToRedis(obj)
	if err != nil {
		return nil, err
	}

	c.values[objId] = redisObj
	return redisObj, nil
}

func (c *valuesCache) getArgs(sessionId string) (redis.Args, error) {
	redisValues, err := c.get(sessionId)
	if err != nil {
		return nil, err
	}
	return redis.Args{}.Add(sessionId).AddFlat(redisValues), nil
}
