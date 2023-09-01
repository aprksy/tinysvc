package redisstore

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// REDIS STORAGE
type RedisStore struct {
	pool *redis.Pool
}

func New(host string, port int) *RedisStore {
	return &RedisStore{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial:        func() (redis.Conn, error) { return redis.Dial("tcp", fmt.Sprintf("%s:%d", host, port)) },
		},
	}
}

func (s *RedisStore) AllFields(key string) (result []string, err error) {
	worstCaseResult := []string{}
	conn := s.pool.Get()
	defer conn.Close()

	result, err = redis.Strings(conn.Do("HKEYS", key))
	if err != nil {
		return worstCaseResult, err
	}
	return
}

func (s *RedisStore) AllValues(key string) (result [][]byte, err error) {
	worstCaseResult := [][]byte{}

	conn := s.pool.Get()
	defer conn.Close()
	values, err := conn.Do("HVALS", key)
	if err != nil {
		return worstCaseResult, err
	}

	preresult := worstCaseResult
	for _, v := range values.([]interface{}) {
		if v != nil {
			preresult = append(preresult, v.([]byte))
		} else {
			preresult = append(preresult, []byte{})
		}
	}
	result = preresult
	return
}

func (s *RedisStore) AllFieldValues(key string) (result map[string][]byte, err error) {
	worstCaseResult := map[string][]byte{}
	conn := s.pool.Get()
	defer conn.Close()
	fields, e1 := s.AllFields(key)
	if e1 != nil {
		return worstCaseResult, e1
	}
	vals, e2 := redis.ByteSlices(conn.Do("HVALS", key))
	if e2 != nil {
		return worstCaseResult, e2
	}
	preresult := worstCaseResult
	for i, f := range fields {
		preresult[f] = vals[i]
	}
	result = preresult
	return
}

func (s *RedisStore) Value(key, field string) (result []byte, err error) {
	worstCaseResult := []byte{}
	conn := s.pool.Get()
	defer conn.Close()
	if key != "" {
		result, err = redis.Bytes(conn.Do("HGET", key, field))
	} else {
		result, err = redis.Bytes(conn.Do("GET", field))
	}
	if err != nil {
		return worstCaseResult, err
	}
	return
}

func (s *RedisStore) Values(key string, fields []string) (result [][]byte, err error) {
	worstCaseResult := [][]byte{}
	conn := s.pool.Get()
	defer conn.Close()
	args := append([]interface{}{}, key)
	fieldIntf := []interface{}{}
	for _, e := range fields {
		fieldIntf = append(fieldIntf, e)
	}
	args = append(args, fieldIntf...)
	intfs, err := conn.Do("HMGET", args...)
	if err != nil {
		return worstCaseResult, err
	}
	preresult := worstCaseResult
	for _, v := range intfs.([]interface{}) {
		if v != nil {
			preresult = append(preresult, v.([]byte))
		} else {
			preresult = append(preresult, []byte{})
		}
	}
	result = preresult
	return
}

func (s *RedisStore) SetField(key, field string, value []byte) (err error) {
	conn := s.pool.Get()
	defer conn.Close()
	if key != "" {
		_, err = conn.Do("HSET", key, field, value)
	} else {
		_, err = conn.Do("SET", field, value)
	}
	return
}

func (s *RedisStore) SetExpireAge(key string, value int) (err error) {
	if value <= 0 {
		return fmt.Errorf("value cannot less than 0")
	}
	conn := s.pool.Get()
	defer conn.Close()
	_, err = conn.Do("EXPIRE", key, value)
	return
}

func (s *RedisStore) SetExpireTime(key string, value time.Time) (err error) {
	if time.Now().After(value) {
		return fmt.Errorf("value cannot before than now")
	}
	conn := s.pool.Get()
	defer conn.Close()
	_, err = conn.Do("EXPIREAT", key, value.Unix())
	return
}

func (s *RedisStore) DelField(key, field string) (err error) {
	conn := s.pool.Get()
	defer conn.Close()
	_, err = conn.Do("HDEL", key, field)
	return
}

func (s *RedisStore) ClearFields(key string) (err error) {
	conn := s.pool.Get()
	defer conn.Close()
	fstrs, e := s.AllFields(key)
	if e != nil {
		err = e
		return
	}
	fields := []interface{}{key}
	for _, f := range fstrs {
		fields = append(fields, f)
	}
	_, err = conn.Do("HDEL", fields...)
	return
}
