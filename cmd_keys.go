package sredis

const (
	CMD_del       = "DEL"
	CMD_dump      = "DUMP"
	CMD_exists    = "EXISTS"
	CMD_expire    = "EXPIRE"
	CMD_expireat  = "EXPIREAT"
	CMD_keys      = "KEYS"
	CMD_migrate   = "MIGRATE"
	CMD_move      = "MOVE"
	CMD_object    = "OBJECT"
	CMD_persist   = "PERSIST"
	CMD_pexpire   = "PEXPIRE"
	CMD_pexpireat = "PEXPIREAT"
	CMD_pttl      = "PTTL"
	CMD_randomkey = "RANDOMKEY"
	CMD_rename    = "RENAME"
	CMD_renamenx  = "RENAMENX"
	CMD_restore   = "RESTORE"
	CMD_scan      = "SCAN"
	CMD_sort      = "SORT"
	CMD_ttl       = "TTL"
	CMD_type      = "TYPE"
	CMD_wait      = "WAIT"
)

var (
	Arg_Match = "MATCH"
	Arg_Count = "COUNT"
)

//////////////////////////////////////////////////////////
//////////////////        keys       /////////////////////
//////////////////////////////////////////////////////////
func (cmd *Commond) Del(key string) *Result {
	cmd.dostr1(CMD_del, key)
	return &cmd.result
}

func (cmd *Commond) Dump(key string) *Result {
	cmd.dostr1(CMD_dump, key)
	return &cmd.result
}

func (cmd *Commond) Exist(key string) *Result {
	cmd.dostr1(CMD_exists, key)
	return &cmd.result
}

func (cmd *Commond) Expire(key string, val int) *Result {
	cmd.doint1(CMD_expire, key, int64(val))
	return &cmd.result
}

func (cmd *Commond) ExpireAt(key string, val int64) *Result {
	cmd.doint1(CMD_expireat, key, val)
	return &cmd.result
}

func (cmd *Commond) Keys(key string) *Result {
	cmd.dostr1(CMD_keys, key)
	return &cmd.result
}

// Remove the existing timeout on key, turning the key from volatile (a key with an expire set) to persistent
func (cmd *Commond) Persist(key string) *Result {
	cmd.dostr1(CMD_persist, key)
	return &cmd.result
}

// milliseconds ttl time
func (cmd *Commond) PExpire(key string, val int64) *Result {
	cmd.doint1(CMD_pexpire, key, val)
	return &cmd.result
}

// milliseconds Unix time
func (cmd *Commond) PExpireAt(key string, val int64) *Result {
	cmd.doint1(CMD_pexpireat, key, val)
	return &cmd.result
}

// TTL in milliseconds
func (cmd *Commond) PTTL(key string) *Result {
	cmd.dostr1(CMD_pttl, key)
	return &cmd.result
}

func (cmd *Commond) Rename(key string, name string) *Result {
	cmd.dostr2(CMD_rename, key, name)
	return &cmd.result
}

func (cmd *Commond) Renamenx(key string, name string) *Result {
	cmd.dostr2(CMD_renamenx, key, name)
	return &cmd.result
}

func (cmd *Commond) Scan(index, count int64, match string) *Result {
	arglen := 2
	if count > 0 {
		arglen += 2
	}
	if match != "" {
		arglen += 2
	}
	cmd.buf.WriteCmd(CMD_scan, arglen)
	cmd.buf.WriteInt64(index)
	if match != "" {
		cmd.buf.WriteString(&Arg_Match)
		cmd.buf.WriteString(&match)
	}
	if count > 0 {
		cmd.buf.WriteString(&Arg_Count)
		cmd.buf.WriteInt64(count)
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) Sort(key string) *Result {
	cmd.dostr1(CMD_sort, key)
	return &cmd.result
}

func (cmd *Commond) Ttl(key string) *Result {
	cmd.dostr1(CMD_ttl, key)
	return &cmd.result
}

func (cmd *Commond) Type(key string) *Result {
	cmd.dostr1(CMD_type, key)
	return &cmd.result
}
