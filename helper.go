package sredis

import (
	"crypto/sha1"
	"io"
	"strings"
)

type RedisOp struct {
	*Commond
}

func NewRedisOp(key int) *RedisOp {
	return &RedisOp{Commond: newCommond(Get(key), 0)}
}

func (this *RedisOp) Set(key string, val interface{}) error {
	return this.Commond.Set(key, val).Err
}

func (this *RedisOp) SetEX(key string, seconds int64, val interface{}) error {
	return this.Commond.Setex(key, seconds, val).Err
}

func (this *RedisOp) SetObject(key string, obj interface{}) error {
	return this.Commond.Do("HMSET", Args{}.Add(key).AddFlat(obj)...).Err
}

func (this *RedisOp) SetField(key string, fieldKey string, obj interface{}) error {
	return this.Commond.HSet(key, fieldKey, obj).Err
}

func (this *RedisOp) Get(key string) *Result {
	return this.Commond.Get(key)
}

func (this *RedisOp) GetObject(key string, obj interface{}) error {
	v, err := this.Commond.HGetAll(key).Values()
	if err != nil {
		return err
	}
	return ScanStruct(v, obj)
}

func (this *RedisOp) Incrby(key string, num int64) error {
	return this.Commond.Incrby(key, num).Err
}

func (this *RedisOp) FieldIncrby(key, fieldKey string, num int64) error {
	return this.Commond.HIncrby(key, fieldKey, num).Err
}

func (this *RedisOp) GetField(objKey, fieldKey string) *Result {
	return this.Commond.HGet(objKey, fieldKey)
}

func (this *RedisOp) Exist(key string) bool {
	v, err := this.Commond.Exist(key).Bool()
	if err != nil {
		return false
	}
	return v
}

func (this *RedisOp) ExistField(objKey, fieldKey string) bool {
	v, err := this.Commond.HExists(objKey, fieldKey).Bool()
	if err != nil {
		return false
	}
	return v
}

func (this *RedisOp) DelField(objKey, fieldKey string) bool {
	v, err := this.Commond.HDel(objKey, fieldKey).Bool()
	if err != nil {
		return false
	}
	return v
}

func (this *RedisOp) Del(key string) bool {
	v, err := this.Commond.Del(key).Bool()
	if err != nil {
		return false
	}
	return v
}

func (this *RedisOp) SetExpire(key string, second int) error {
	return this.Commond.Expire(key, second).Err
}

func (this *RedisOp) SetExpireAt(key string, timestamp int64) error {
	return this.Commond.ExpireAt(key, timestamp).Err
}

func (this *RedisOp) Do(commandName string, args ...interface{}) *Result {
	return this.Commond.Do(commandName, args...)
}

func (this *RedisOp) Send(commandName string, args ...interface{}) error {
	return this.Commond.Do(commandName, args...).Err
}

func (this *RedisOp) GetFields(objKey string, args ...string) *Result {
	return this.Commond.Do("HMGet", Args{}.Add(objKey).AddFlat(args)...)
}

func (this *RedisOp) SetFields(objKey string, args ...interface{}) error {
	return this.Commond.HMSet("HMSet", Args{}.Add(objKey).AddFlat(args)...).Err
}

func (this *RedisOp) DoScript(script string, key string, val interface{}) *Result {
	h := sha1.New()
	io.WriteString(h, script)
	hkey := h.Sum(nil)
	result := this.EvalSha1(string(hkey), key, val)
	if result.Err != nil && strings.HasPrefix(result.Err.Error(), "NOSCRIPT ") {
		result = this.Eval1(script, key, val)
	}
	return result
}
