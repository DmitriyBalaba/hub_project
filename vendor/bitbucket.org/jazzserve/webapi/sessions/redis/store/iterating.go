package store

import (
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type getArgsFunc func(sessionId string) (redis.Args, error)

func (s *SessionStore) updateAllSessions(pattern string, getArgs getArgsFunc) error {
	var (
		cursor     int
		sessionIds []string
	)

	conn := s.pool.Get()
	defer conn.Close()

	for {
		values, err := redis.Values(conn.Do(Scan, cursor, Match, pattern))
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = redis.Scan(values, &cursor, &sessionIds)
		if err != nil {
			return errors.WithStack(err)
		}

		err = doHmsetInPipeline(conn, sessionIds, getArgs)
		if err != nil {
			return errors.WithStack(err)
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}

func doHmsetInPipeline(conn redis.Conn, sessionIds []string, getArgs getArgsFunc) error {
	if len(sessionIds) == 0 {
		return nil
	}

	var i int

	// adding request to the pipeline
	for i = 0; i < len(sessionIds)-1; i++ {
		args, err := getArgs(sessionIds[i])
		if err != nil {
			return err
		}

		err = conn.Send(HMSet, args...)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	args, err := getArgs(sessionIds[i])
	if err != nil {
		return err
	}

	// actually sending requests to the server by doing the last one
	_, err = conn.Do(HMSet, args...)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

//delete all sessionIds (keys)
func doDelInPipeline(conn redis.Conn, sessionIds []string) error {
	if len(sessionIds) == 0 {
		return nil
	}

	var i int
	for i = 0; i < len(sessionIds)-1; i++ {
		err := conn.Send(Del, sessionIds[i])
		if err != nil {
			return errors.WithStack(err)
		}
	}
	_, err := conn.Do(Del, sessionIds[i])
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

//delete all sessions which ids match the given pattern
func (s *SessionStore) DeleteAllSessions(pattern string) error {
	var (
		cursor     int
		sessionIds []string
	)

	conn := s.pool.Get()
	defer conn.Close()

	for {
		values, err := redis.Values(conn.Do(Scan, cursor, Match, pattern))
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = redis.Scan(values, &cursor, &sessionIds)
		if err != nil {
			return errors.WithStack(err)
		}

		if err = doDelInPipeline(conn, sessionIds); err != nil {
			return errors.WithStack(err)
		}

		if cursor == 0 {
			break
		}
	}

	return nil
}
