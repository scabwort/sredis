package sredis

const (
	CMD_append      = "APPEND"
	CMD_bitcount    = "BITCOUNT"
	CMD_bitop       = "BITOP"
	CMD_bitpos      = "BITPOS"
	CMD_decr        = "DECR"
	CMD_decrby      = "DECRBY"
	CMD_get         = "GET"
	CMD_getbit      = "GETBIT"
	CMD_getrange    = "GETRANGE"
	CMD_getset      = "GETSET"
	CMD_incr        = "INCR"
	CMD_incrby      = "INCRBY"
	CMD_incrbyfloat = "INCRBYFLOAT"
	CMD_mget        = "MGET"
	CMD_mset        = "MSET"
	CMD_msetnx      = "MSETNX"
	CMD_psetex      = "PSETEX"
	CMD_set         = "SET"
	CMD_setbit      = "SETBIT"
	CMD_setex       = "SETEX"
	CMD_setnx       = "SETNX"
	CMD_setrange    = "SETRANGE"
	CMD_strlen      = "STRLEN"
)

func (cmd *Commond) Append(key string, val interface{}) *Result {
	cmd.dostr2arg(CMD_append, key, val)
	return &cmd.result
}

func (cmd *Commond) Bitcount(key string) *Result {
	cmd.dostr1(CMD_append, key)
	return &cmd.result
}

func (cmd *Commond) Bitop(op, dest, key1, key2 string) *Result {
	cmd.dostr4(CMD_bitop, op, dest, key1, key2)
	return &cmd.result
}

func (cmd *Commond) Bitpos(key string, start, end int64) *Result {
	cmd.doint2(CMD_bitpos, key, start, end)
	return &cmd.result
}

func (cmd *Commond) Decr(key string) *Result {
	cmd.dostr1(CMD_decr, key)
	return &cmd.result
}

func (cmd *Commond) Decrby(key string, num int64) *Result {
	cmd.dostr2int(CMD_decrby, key, num)
	return &cmd.result
}

func (cmd *Commond) Get(key string) *Result {
	cmd.dostr1(CMD_get, key)
	return &cmd.result
}

func (cmd *Commond) Getbit(key string, num int64) *Result {
	cmd.doint1(CMD_getbit, key, num)
	return &cmd.result
}

func (cmd *Commond) Getrange(key string, start, end int64) *Result {
	cmd.doint2(CMD_getrange, key, start, end)
	return &cmd.result
}

// Atomically sets key to value and returns the old value stored at key
func (cmd *Commond) Getset(key string, val interface{}) *Result {
	cmd.dostr2arg(CMD_getset, key, val)
	return &cmd.result
}

func (cmd *Commond) Incr(key string) *Result {
	cmd.dostr1(CMD_incr, key)
	return &cmd.result
}

func (cmd *Commond) Incrby(key string, num int64) *Result {
	cmd.dostr2int(CMD_incrby, key, num)
	return &cmd.result
}

func (cmd *Commond) Incrbyfloat(key string, num float64) *Result {
	cmd.buf.WriteCmd(CMD_incrbyfloat, 3)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteFloat64(&num)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) Mget(keys ...string) *Result {
	cmd.buf.WriteCmd(CMD_mget, 1+len(keys))
	for idx, _ := range keys {
		cmd.buf.WriteString(&keys[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) Mset(keys ...interface{}) *Result {
	cmd.buf.WriteCmd(CMD_mset, 1+len(keys))
	for idx, _ := range keys {
		cmd.buf.WriteArg(keys[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) Msetnx(keys ...interface{}) *Result {
	cmd.buf.WriteCmd(CMD_msetnx, 1+len(keys))
	for idx, _ := range keys {
		cmd.buf.WriteArg(keys[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

// Set key to hold the string value and set key to timeout after a given number of seconds
func (cmd *Commond) Setex(key string, ttl int64, val interface{}) *Result {
	cmd.buf.WriteCmd(CMD_setex, 4)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteInt64(ttl)
	cmd.buf.WriteArg(val)
	cmd.waitConn()
	return &cmd.result
}

// PSETEX works exactly like SETEX with the sole difference that the expire time is specified in milliseconds instead of seconds
func (cmd *Commond) Psetex(key string, ttl int64, val interface{}) *Result {
	cmd.buf.WriteCmd(CMD_psetex, 4)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteInt64(ttl)
	cmd.buf.WriteArg(val)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) Set(key string, val interface{}) *Result {
	cmd.dostr2arg(CMD_set, key, val)
	return &cmd.result
}

func (cmd *Commond) SetInt(key string, val int64) *Result {
	cmd.doint1(CMD_set, key, val)
	return &cmd.result
}

func (cmd *Commond) SetStr(key string, val string) *Result {
	cmd.dostr2(CMD_set, key, val)
	return &cmd.result
}

// Sets or clears the bit at offset in the string value stored at key
func (cmd *Commond) Setbit(key string, start, end int64) *Result {
	cmd.doint2(CMD_setbit, key, start, end)
	return &cmd.result
}

// Set key to hold string value if key does not exist
func (cmd *Commond) Setnx(key string, val interface{}) *Result {
	cmd.dostr2arg(CMD_setnx, key, val)
	return &cmd.result
}

func (cmd *Commond) setrange(key string, pos int64, val interface{}) *Result {
	cmd.buf.WriteCmd(CMD_setrange, 4)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteInt64(pos)
	cmd.buf.WriteArg(val)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) strlen(key string) *Result {
	cmd.dostr1(CMD_strlen, key)
	return &cmd.result
}
