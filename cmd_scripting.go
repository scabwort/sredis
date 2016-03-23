package sredis

const (
	CMD_eval          = "EVAL"
	CMD_evalsha       = "EVALSHA"
	CMD_script_exists = "SCRIPT EXISTS"
	CMD_script_flush  = "SCRIPT FLUSH"
	CMD_script_kill   = "SCRIPT KILL"
	CMD_script_load   = "SCRIPT LOAD"
)

func (cmd *Commond) Eval(script string, argCount int, vals ...interface{}) *Result {
	cmd.Cmd = CMD_eval
	cmd.buf.WriteCmd(CMD_eval, 3+len(vals))
	cmd.buf.WriteString(&script)
	cmd.buf.WriteUInt64(uint64(argCount))
	for idx, _ := range vals {
		cmd.buf.WriteArg(vals[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) Eval1(script string, key string, val interface{}) *Result {
	cmd.Cmd = CMD_eval
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_eval, 5)
	cmd.buf.WriteString(&script)
	cmd.buf.WriteUInt64(1)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteArg(val)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) EvalSha(hash string, argCount int, vals ...interface{}) *Result {
	cmd.Cmd = CMD_evalsha
	cmd.buf.WriteCmd(CMD_evalsha, 3+len(vals))
	cmd.buf.WriteString(&hash)
	cmd.buf.WriteUInt64(uint64(argCount))
	for idx, _ := range vals {
		cmd.buf.WriteArg(vals[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) EvalSha1(hash string, key string, val interface{}) *Result {
	cmd.Cmd = CMD_evalsha
	cmd.Key = key
	cmd.buf.WriteCmd(CMD_evalsha, 5)
	cmd.buf.WriteString(&hash)
	cmd.buf.WriteUInt64(1)
	cmd.buf.WriteString(&key)
	cmd.buf.WriteArg(val)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) ScriptExists(hash string) *Result {
	cmd.dostr1(CMD_script_exists, hash)
	return &cmd.result
}

func (cmd *Commond) ScriptKill() *Result {
	cmd.Cmd = CMD_script_kill
	cmd.buf.WriteCmd(CMD_script_kill, 1)
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) ScriptLoad(script string) *Result {
	cmd.dostr1(CMD_script_load, script)
	return &cmd.result
}
