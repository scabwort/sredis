package sredis

const (
	CMD_hdel         = "HDEL"
	CMD_hexists      = "HEXISTS"
	CMD_hget         = "HGET"
	CMD_hgetall      = "HGETALL"
	CMD_hincrby      = "HINCRBY"
	CMD_hincrbyFloat = "HINCRBYFLOAT"
	CMD_hkeys        = "HKEYS"
	CMD_hlen         = "HLEN"
	CMD_hmget        = "HMGET"
	CMD_hmset        = "HMSET"
	CMD_hscan        = "HSCAN"
	CMD_hset         = "HSET"
	CMD_hsetnx       = "HSETNX"
	CMD_hstrlen      = "HSTRLEN"
	CMD_hvals        = "HVALS"
)

//////////////////////////////////////////////////////////
//////////////////        hashs       ////////////////////
//////////////////////////////////////////////////////////

func (cmd *Commond) HDel(key, field string) *Result {
	cmd.dostr2(CMD_hdel, key, field)
	return &cmd.result
}

func (cmd *Commond) HExists(key, field string) *Result {
	cmd.dostr2(CMD_hexists, key, field)
	return &cmd.result
}

func (cmd *Commond) HGet(key, field string) *Result {
	cmd.dostr2(CMD_hget, key, field)
	return &cmd.result
}

func (cmd *Commond) HSet(key string, field string, val interface{}) *Result {
	cmd.dostr3arg(CMD_hset, key, field, val)
	return &cmd.result
}

func (cmd *Commond) HSetStr(key string, field string, val string) *Result {
	cmd.dostr3(CMD_hset, key, field, val)
	return &cmd.result
}

func (cmd *Commond) HSetInt(key string, field string, val int64) *Result {
	cmd.dostr3int(CMD_hset, key, field, val)
	return &cmd.result
}

func (cmd *Commond) HSetnx(key, field string, val interface{}) *Result {
	cmd.dostr3arg(CMD_hsetnx, key, field, val)
	return &cmd.result
}

func (cmd *Commond) HGetAll(key string) *Result {
	cmd.dostr1(CMD_hgetall, key)
	return &cmd.result
}

func (cmd *Commond) HIncrby(key, field string, incr int64) *Result {
	cmd.dostr3int(CMD_hincrby, key, field, incr)
	return &cmd.result
}

func (cmd *Commond) HKeys(key string) *Result {
	cmd.dostr1(CMD_hkeys, key)
	return &cmd.result
}

func (cmd *Commond) HLen(key string) *Result {
	cmd.dostr1(CMD_hlen, key)
	return &cmd.result
}

func (cmd *Commond) HMGet(key string, fields ...string) *Result {
	cmd.Cmd = CMD_hmget
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_hmget, len(fields)+2)
	cmd.buf.WriteString(&key)
	for idx, _ := range fields {
		cmd.buf.WriteString(&fields[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) HMSet(key string, values ...interface{}) *Result {
	if len(values) == 0 || len(values)%2 != 0 {
		cmd.result.Err = ErrCommondArg
		return &cmd.result
	}
	cmd.Cmd = CMD_hmget
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_hmget, len(values)+2)
	cmd.buf.WriteString(&key)
	for idx, _ := range values {
		cmd.buf.WriteArg(values[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) HScan(key string, count int64) *Result {
	cmd.dostr2int(CMD_hscan, key, count)
	return &cmd.result
}

func (cmd *Commond) HStrlen(key, field string) *Result {
	cmd.dostr2(CMD_hstrlen, key, field)
	return &cmd.result
}
func (cmd *Commond) HVals(key string) *Result {
	cmd.dostr1(CMD_hvals, key)
	return &cmd.result
}
