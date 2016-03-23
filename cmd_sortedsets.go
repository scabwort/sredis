package sredis

const (
	CMD_zadd             = "ZADD"
	CMD_zcard            = "ZCARD"
	CMD_zcount           = "ZCOUNT"
	CMD_zincrby          = "ZINCRBY"
	CMD_zinterstore      = "ZINTERSTORE"
	CMD_zlexcount        = "ZLEXCOUNT"
	CMD_zrange           = "ZRANGE"
	CMD_zrangebylex      = "ZRANGEBYLEX"
	CMD_zrangebyscore    = "ZRANGEBYSCORE"
	CMD_zrank            = "ZRANK"
	CMD_zrem             = "ZREM"
	CMD_zremrangebylex   = "ZREMRANGEBYLEX"
	CMD_zremrangebyrank  = "ZREMRANGEBYRANK"
	CMD_zremrangebyscore = "ZREMRANGEBYSCORE"
	CMD_zrevrange        = "ZREVRANGE"
	CMD_zrevrangebylex   = "ZREVRANGEBYLEX"
	CMD_zrevrangebyscore = "ZREVRANGEBYSCORE"
	CMD_zrevrank         = "ZREVRANK"
	CMD_zscan            = "ZSCAN"
	CMD_zscore           = "ZSCORE"
	CMD_zunionstore      = "ZUNIONSTORE"
)

var (
	Arg_Weights    = "WEIGHTS"
	Arg_Withscores = "WITHSCORES"
)

func (cmd *Commond) ZAdd(key string, name string, val interface{}) *Result {
	cmd.dostr3arg(CMD_zadd, key, name, val)
	return &cmd.result
}

func (cmd *Commond) ZCard(key string) *Result {
	cmd.dostr1(CMD_zcard, key)
	return &cmd.result
}

func (cmd *Commond) ZCount(key, start, end string) *Result {
	cmd.dostr3(CMD_zcount, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZIncrby(key, name string, val int64) *Result {
	cmd.dostr3int(CMD_zincrby, key, name, val)
	return &cmd.result
}

// Computes the intersection of numkeys sorted sets given by the specified keys
func (cmd *Commond) ZInterstore(key, key1, key2 string, rate1, rate2 int64) *Result {
	cmd.buf.WriteCmd(CMD_zinterstore, 8)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteInt64(2)
	cmd.buf.WriteString(&key1)
	cmd.buf.WriteString(&key2)
	cmd.buf.WriteString(&Arg_Weights)
	cmd.buf.WriteInt64(rate1)
	cmd.buf.WriteInt64(rate2)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) ZLexcount(key, start, end string) *Result {
	cmd.dostr3(CMD_zlexcount, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRange(key string, start, end int64) *Result {
	cmd.doint2(CMD_zrange, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRangeWith(key string, start, end int64) *Result {
	cmd.buf.WriteCmd(CMD_zrange, 3)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteInt64(start)
	cmd.buf.WriteInt64(end)
	cmd.buf.WriteString(&Arg_Withscores)
	return &cmd.result
}

func (cmd *Commond) ZRangebylex(key, start, end string) *Result {
	cmd.dostr3(CMD_zrangebylex, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRangebyscore(key, start, end string) *Result {
	cmd.dostr3(CMD_zrangebyscore, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRank(key string, name string) *Result {
	cmd.dostr2(CMD_zrank, key, name)
	return &cmd.result
}

func (cmd *Commond) ZRem(key string, name string) *Result {
	cmd.dostr2(CMD_zrank, key, name)
	return &cmd.result
}

func (cmd *Commond) ZRemrangebylex(key, start, end string) *Result {
	cmd.dostr3(CMD_zremrangebylex, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRemrangebyrank(key, start, end string) *Result {
	cmd.dostr3(CMD_zremrangebyrank, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRemrangebyscore(key, start, end string) *Result {
	cmd.dostr3(CMD_zremrangebyscore, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRevrange(key string, start, end int64) *Result {
	cmd.doint2(CMD_zrevrange, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRevrangebylex(key string, start, end string) *Result {
	cmd.dostr3(CMD_zrevrangebylex, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRevrangebyscore(key string, start, end string) *Result {
	cmd.dostr3(CMD_zrevrangebyscore, key, start, end)
	return &cmd.result
}

func (cmd *Commond) ZRevrank(key, name string, val interface{}) *Result {
	cmd.dostr3arg(CMD_zrevrank, key, name, val)
	return &cmd.result
}

func (cmd *Commond) ZScan(key string) *Result {
	return &cmd.result
}

func (cmd *Commond) ZScore(key string, name string) *Result {
	cmd.dostr2(CMD_zscore, key, name)
	return &cmd.result
}

func (cmd *Commond) ZUnionstore(key, key1, key2 string, rate1, rate2 int64) *Result {
	cmd.buf.WriteCmd(CMD_zunionstore, 8)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteInt64(2)
	cmd.buf.WriteString(&key1)
	cmd.buf.WriteString(&key2)
	cmd.buf.WriteString(&Arg_Weights)
	cmd.buf.WriteInt64(rate1)
	cmd.buf.WriteInt64(rate2)
	cmd.waitConn()
	return &cmd.result
}
