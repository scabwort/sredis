package sredis

import (
	"errors"
	"fmt"
	"strconv"
)

type Error string

func (err Error) Error() string { return string(err) }

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("[redis] nil returned")
var errNegativeInt = errors.New("redis: unexpected value for Uint64")

type Result struct {
	Data interface{}
	Err  error
}

func (r *Result) Interface() (interface{}, error) {
	return r.Data, r.Err
}

func (r *Result) Int() (int, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return 0, fmt.Errorf("redis: unexpected type for Int, got type %T", r.Data)
	}
	return parseToInt(replay)
}

func Int(r *Result) (int, error) {
	return r.Int()
}

func (r *Result) Int32() (int32, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return 0, fmt.Errorf("redis: unexpected type for Int, got type %T", r.Data)
	}
	n, err := parseToInt(replay)
	return int32(n), err
}

func (r *Result) UInt32() (uint32, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return 0, fmt.Errorf("redis: unexpected type for UInt32, got type %T", r.Data)
	}
	n, err := parseToUInt64(replay)
	return uint32(n), err
}

func UInt32(r *Result) (uint32, error) {
	return r.UInt32()
}

func (r *Result) Int64() (int64, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return 0, fmt.Errorf("redis: unexpected type for Int, got type %T", r.Data)
	}
	return parseToInt64(replay)
}

func Int64(r *Result) (int64, error) {
	return r.Int64()
}

func (r *Result) Uint64() (uint64, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return 0, fmt.Errorf("redis: unexpected type for Int, got type %T", r.Data)
	}
	return parseToUInt64(replay)
}

func Uint64(r *Result) (uint64, error) {
	return r.Uint64()
}

func (r *Result) Float64() (float64, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return 0, fmt.Errorf("redis: unexpected type for Int, got type %T", r.Data)
	}
	return strconv.ParseFloat(string(replay), 64)
}

func Float64(r *Result) (float64, error) {
	return r.Float64()
}

func (r *Result) String() (n string, err error) {
	if r.Err != nil {
		return "", r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return "", fmt.Errorf("redigo: unexpected type for Int, got type %T", r.Data)
	}
	n = string(replay)
	return
}

func String(r *Result) (string, error) {
	return r.String()
}

func (r *Result) Bytes() (b []byte, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for []byte, got type %T", r.Data)
	}
	b = replay
	return
}

func Bytes(r *Result) ([]byte, error) {
	return r.Bytes()
}

func (r *Result) Bool() (bool, error) {
	if r.Err != nil {
		return false, r.Err
	}
	replay, ok := r.Data.([]byte)
	if !ok {
		return false, fmt.Errorf("redis: unexpected type for Int, got type %T", r.Data)
	}
	return strconv.ParseBool(string(replay))
}

func Bool(r *Result) (bool, error) {
	return r.Bool()
}

func (r *Result) Values() ([]interface{}, error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for []interface{}, got type %T", r.Data)
	}
	return replay, nil
}

func Values(r *Result) ([]interface{}, error) {
	return r.Values()
}

func (r *Result) Strings() (result []string, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for Int, got []interface{} %T", r.Data)
	}
	result = make([]string, len(replay))
	for i, _ := range replay {
		if replay[i] == nil {
			continue
		}
		p, ok := replay[i].([]byte)
		if !ok {
			return nil, fmt.Errorf("redis: unexpected element type for string, got type %T", replay[i])
		}
		result[i] = string(p)
	}
	return
}

func Strings(r *Result) ([]string, error) {
	return r.Strings()
}

func (r *Result) Ints() (result []int, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for Ints, got type %T", r.Data)
	}
	n := 0
	result = make([]int, len(replay))
	for i := range replay {
		if replay[i] == nil {
			continue
		}
		p, ok := replay[i].([]byte)
		if !ok {
			return nil, fmt.Errorf("redis: unexpected element type for Ints, got type %T", replay[i])
		}
		n, err = parseToInt(p)
		if err != nil {
			return nil, err
		}
		result[i] = n
	}
	return
}

func Ints(r *Result) ([]int, error) {
	return r.Ints()
}

func (r *Result) UInt32s() (result []uint32, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for Ints, got type %T", r.Data)
	}
	var n uint64
	result = make([]uint32, len(replay))
	for i := range replay {
		if replay[i] == nil {
			continue
		}
		p, ok := replay[i].([]byte)
		if !ok {
			return nil, fmt.Errorf("redis: unexpected element type for Ints, got type %T", replay[i])
		}
		n, err = parseToUInt64(p)
		if err != nil {
			return nil, err
		}
		result[i] = uint32(n)
	}
	return
}

func UInt32s(r *Result) ([]uint32, error) {
	return r.UInt32s()
}

func (r *Result) Int64s() (result []int64, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for Ints, got type %T", r.Data)
	}
	var n int64
	result = make([]int64, len(replay))
	for i := range replay {
		if replay[i] == nil {
			continue
		}
		p, ok := replay[i].([]byte)
		if !ok {
			return nil, fmt.Errorf("redis: unexpected element type for Ints, got type %T", replay[i])
		}
		n, err = parseToInt64(p)
		if err != nil {
			return nil, err
		}
		result[i] = n
	}
	return
}

func Int64s(r *Result) ([]int64, error) {
	return r.Int64s()
}

func (r *Result) UInt64s() (result []uint64, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for Ints, got type %T", r.Data)
	}
	var n uint64
	result = make([]uint64, len(replay))
	for i := range replay {
		if replay[i] == nil {
			continue
		}
		p, ok := replay[i].([]byte)
		if !ok {
			return nil, fmt.Errorf("redis: unexpected element type for Ints, got type %T", replay[i])
		}
		n, err = parseToUInt64(p)
		if err != nil {
			return nil, err
		}
		result[i] = n
	}
	return
}

func UInt64s(r *Result) ([]uint64, error) {
	return r.UInt64s()
}

func (r *Result) StringMap() (result map[string]string, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for Int, got type %T", r.Data)
	}
	var (
		key, value     []byte
		okKey, okValue bool
	)
	result = make(map[string]string, len(replay)/2)
	for i := 0; i < len(replay); i += 2 {
		key, okKey = replay[i].([]byte)
		value, okValue = replay[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redis: ScanMap key not a bulk string value")
		}
		result[string(key)] = string(value)
	}
	return
}

func StringMap(r *Result) (map[string]string, error) {
	return r.StringMap()
}

func (r *Result) IntMap() (result map[string]int, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for []interface{}, got type %T", r.Data)
	}
	var (
		key, value     []byte
		okKey, okValue bool
		v              int
	)
	result = make(map[string]int, len(replay)/2)
	for i := 0; i < len(replay); i += 2 {
		key, okKey = replay[i].([]byte)
		value, okValue = replay[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redis: ScanMap key not a bulk string value")
		}
		v, err = parseToInt(value)
		if err != nil {
			return nil, err
		}
		result[string(key)] = v
	}
	return
}

func IntMap(r *Result) (map[string]int, error) {
	return r.IntMap()
}

func (r *Result) Int64Map() (result map[string]int64, err error) {
	if r.Err != nil {
		return nil, r.Err
	}
	replay, ok := r.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("redis: unexpected type for []interface{}, got type %T", r.Data)
	}
	var (
		key, value     []byte
		okKey, okValue bool
		v              int64
	)
	result = make(map[string]int64, len(replay)/2)
	for i := 0; i < len(replay); i += 2 {
		key, okKey = replay[i].([]byte)
		value, okValue = replay[i+1].([]byte)
		if !okKey || !okValue {
			return nil, errors.New("redis: ScanMap key not a bulk string value")
		}
		v, err = parseToInt64(value)
		if err != nil {
			return nil, err
		}
		result[string(key)] = v
	}
	return
}

func Int64Map(r *Result) (map[string]int64, error) {
	return r.Int64Map()
}

// parseInt parses an integer reply.
func parseToInt(p []byte) (n int, err error) {
	if len(p) == 0 {
		return 0, Error("malformed integer")
	}
	var negate bool
	if p[0] == '-' {
		negate = true
		p = p[1:]
		if len(p) == 0 {
			return 0, Error("malformed integer")
		}
	}
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return 0, Error("illegal bytes in length")
		}
		n += int(b - '0')
	}
	if negate {
		n = -n
	}
	return
}

// parseInt parses an integer reply.
func parseToInt64(p []byte) (n int64, err error) {
	if len(p) == 0 {
		return 0, Error("malformed integer")
	}
	var negate bool
	if p[0] == '-' {
		negate = true
		p = p[1:]
		if len(p) == 0 {
			return 0, Error("malformed integer")
		}
	}
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return 0, Error("illegal bytes in length")
		}
		n += int64(b - '0')
	}
	if negate {
		n = -n
	}
	return
}

func parseToUInt64(p []byte) (n uint64, err error) {
	if len(p) == 0 {
		return 0, Error("malformed integer")
	}
	if p[0] == '-' {
		p = p[1:]
		if len(p) == 0 {
			return 0, Error("malformed integer")
		}
	}
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return 0, Error("illegal bytes in length")
		}
		n += uint64(b - '0')
	}
	return
}

func ToUInt64(p interface{}) (n uint64) {
	if p == nil {
		return 0
	}
	replay, ok := p.([]byte)
	if !ok {
		return 0
	}
	n, _ = parseToUInt64(replay)
	return
}

func ToString(p interface{}) string {
	if p == nil {
		return ""
	}
	replay, ok := p.([]byte)
	if !ok {
		return ""
	}
	return string(replay)
}
