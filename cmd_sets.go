package sredis

const (
	CMD_sadd        = "SADD"
	CMD_scard       = "SCARD"
	CMD_sdiff       = "SDIFF"
	CMD_sdiffstore  = "SDIFFSTORE"
	CMD_sinter      = "SINTER"
	CMD_sinterstore = "SINTERSTORE"
	CMD_sismember   = "SISMEMBER"
	CMD_smembers    = "SMEMBERS"
	CMD_smove       = "SMOVE"
	CMD_spop        = "SPOP"
	CMD_srandmember = "SRANDMEMBER"
	CMD_srem        = "SREM"
	CMD_sscan       = "SSCAN"
	CMD_sunion      = "SUNION"
	CMD_sunionstore = "SUNIONSTORE"
)

func (cmd *Commond) SAdd(key string, val interface{}) *Result {
	cmd.dostr2arg(CMD_sadd, key, val)
	return &cmd.result
}

// Returns the set cardinality (number of elements) of the set stored at key
func (cmd *Commond) SCard(key string) *Result {
	cmd.dostr1(CMD_scard, key)
	return &cmd.result
}

func (cmd *Commond) SDiff(key1, key2 string) *Result {
	cmd.dostr2(CMD_sdiff, key1, key2)
	return &cmd.result
}

func (cmd *Commond) SDiffstore(key, key1, key2 string) *Result {
	cmd.dostr3(CMD_sdiffstore, key, key1, key2)
	return &cmd.result
}

func (cmd *Commond) SInter(key1, key2 string) *Result {
	cmd.dostr2(CMD_sinter, key1, key2)
	return &cmd.result
}

func (cmd *Commond) SInterstore(key, key1, key2 string) *Result {
	cmd.dostr3(CMD_sinterstore, key, key1, key2)
	return &cmd.result
}

func (cmd *Commond) SIsmember(key string, val interface{}) *Result {
	cmd.dostr2arg(CMD_sismember, key, val)
	return &cmd.result
}

// O(N) Returns all the members of the set value stored at key
func (cmd *Commond) SMembers(key string) *Result {
	cmd.dostr1(CMD_smembers, key)
	return &cmd.result
}

// O(1) Removes and returns one or more random elements from the set value store at key
func (cmd *Commond) SMove(from, to string, val interface{}) *Result {
	cmd.dostr3arg(CMD_smove, from, to, val)
	return &cmd.result
}

// O(1) Removes and returns one or more random elements from the set value store at key
func (cmd *Commond) SPop(key string) *Result {
	cmd.dostr1(CMD_spop, key)
	return &cmd.result
}

func (cmd *Commond) SPopNum(key string, count int) *Result {
	cmd.dostr2int(CMD_spop, key, int64(count))
	return &cmd.result
}

// When called with just the key argument, return a random element from the set value stored at key
func (cmd *Commond) SRandmember(key string) *Result {
	cmd.dostr1(CMD_srandmember, key)
	return &cmd.result
}

// When called with just the key argument, return a random element from the set value stored at key
func (cmd *Commond) SRandmemberNum(key string, count int) *Result {
	cmd.dostr2int(CMD_srandmember, key, int64(count))
	return &cmd.result
}

// O(N) Remove the specified members from the set stored at key
func (cmd *Commond) SRem(key string, val interface{}) *Result {
	cmd.dostr2arg(CMD_srem, key, val)
	return &cmd.result
}

func (cmd *Commond) SScan(index, count int64, match string) *Result {
	arglen := 2
	if count > 0 {
		arglen += 2
	}
	if match != "" {
		arglen += 2
	}
	cmd.buf.WriteCmd(CMD_sscan, arglen)
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

func (cmd *Commond) SUnion(key1, key2 string) *Result {
	cmd.dostr2(CMD_sunion, key1, key2)
	return &cmd.result
}

func (cmd *Commond) SUnionstore(key, key1, key2 string) *Result {
	cmd.dostr3(CMD_sunionstore, key, key1, key2)
	return &cmd.result
}
