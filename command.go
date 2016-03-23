package sredis

import (
	"errors"
)

const (
	WriteCmd       = 0
	ReadCmd        = 1
	defaultCmdSize = 128
)

var (
	ErrCommondArg = errors.New("redis:arg is error!")
)

var (
	flagAuth   = "AUTH"
	flagPing   = "PING"
	flagSelect = "SELECT"
)

type ICommond interface {
	GetBytes() []byte
	SetData(t interface{}, err error) bool
	Done()
	String() string
	//	Recycle()
}

type Commond struct {
	Cmd        string
	Key        string
	buf        *Packet
	conn       *Conn
	c          chan struct{}
	result     Result
	isAutoSend bool
	max        int
	idx        int
}

func newCommond(conn *Conn, idx int) (cmd *Commond) {
	cmd = &Commond{}
	cmd.buf = NewPacketSize(defaultCmdSize)
	cmd.c = make(chan struct{}, 1)
	cmd.conn = conn
	cmd.isAutoSend = true
	cmd.idx = idx
	return
}

func (cmd *Commond) ResizeBuf(size int) {
	cmd.buf.grow(size)
}

func (cmd *Commond) Do(c string, args ...interface{}) *Result {
	cmd.buf.WriteCmd(c, len(args)+1)
	for idx, _ := range args {
		cmd.buf.WriteArg(args[idx])
	}
	cmd.waitConn()
	return &cmd.result
}

func (cmd *Commond) Reset(conn *Conn, idx int) {
	cmd.conn = conn
	cmd.idx = idx
	cmd.buf.Reset()
	cmd.result.Data = nil
	cmd.result.Err = nil
}

func (cmd *Commond) Recycle() {
	RecycleCmd(cmd)
}

func (cmd *Commond) GetBytes() (b []byte) {
	b = cmd.buf.buf[:cmd.buf.w]
	return
}

func (cmd *Commond) SetData(data interface{}, err error) bool {
	cmd.result.Data = data
	cmd.result.Err = err
	return true
}

func (cmd *Commond) Done() {
	select {
	case cmd.c <- struct{}{}:
	default:
	}
}

func (cmd *Commond) waitConn() {
	if !cmd.isAutoSend {
		cmd.max++
		return
	}
	cmd.conn.Send(cmd)
	<-cmd.c
	cmd.buf.Reset()
	cmd.max = 0
}

func (cmd *Commond) doint1(c, arg1 string, arg2 int64) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 3)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteInt64(arg2)
	cmd.waitConn()
}

func (cmd *Commond) doint2(c, arg1 string, arg2 int64, arg3 int64) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 4)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteInt64(arg2)
	cmd.buf.WriteInt64(arg3)
	cmd.waitConn()
}

func (cmd *Commond) docmd(c string) {
	cmd.buf.WriteCmd(c, 1)
	cmd.waitConn()
}

func (cmd *Commond) dostr1(c, arg1 string) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 2)
	cmd.buf.WriteString(&arg1)
	cmd.waitConn()
}

func (cmd *Commond) dostr2(c, arg1 string, arg2 string) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 3)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteString(&arg2)
	cmd.waitConn()
}

func (cmd *Commond) dostr2int(c, arg1 string, arg2 int64) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 3)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteInt64(arg2)
	cmd.waitConn()
}

func (cmd *Commond) dostr3(c, arg1, arg2, arg3 string) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 4)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteString(&arg2)
	cmd.buf.WriteString(&arg3)
	cmd.waitConn()
}

func (cmd *Commond) dostr3int(c, arg1, arg2 string, arg3 int64) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 4)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteString(&arg2)
	cmd.buf.WriteInt64(arg3)
	cmd.waitConn()
}

func (cmd *Commond) dostr2arg(c, arg1 string, arg2 interface{}) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 3)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteArg(arg2)
	cmd.waitConn()
}

func (cmd *Commond) dostr3arg(c, arg1, arg2 string, arg3 interface{}) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 4)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteString(&arg2)
	cmd.buf.WriteArg(arg3)
	cmd.waitConn()
}

func (cmd *Commond) dostr4(c, arg1, arg2, arg3, arg4 string) {
	cmd.Cmd = c
	cmd.Key = arg1
	cmd.buf.WriteCmd(c, 5)
	cmd.buf.WriteString(&arg1)
	cmd.buf.WriteString(&arg2)
	cmd.buf.WriteString(&arg3)
	cmd.buf.WriteString(&arg4)
	cmd.waitConn()
}

func (cmd *Commond) dostrmore(c string, args []string) {
	cmd.Cmd = c
	cmd.Key = args[0]
	cmd.buf.WriteCmd(c, 1+len(args))
	for idx, _ := range args {
		cmd.buf.WriteString(&args[idx])
	}
	cmd.waitConn()
}

func (cmd *Commond) String() string {
	return "[cmd] " + cmd.Cmd + ":" + cmd.Key
}

// 批量命令
type BenchCommond struct {
	*Commond
	results []Result
	cur     int
}

func NewBenchCommond(conn *Conn, idx int) (cmd *BenchCommond) {
	cmd = &BenchCommond{Commond: newCommond(conn, idx)}
	cmd.isAutoSend = false
	return
}

func (cmd *BenchCommond) SetData(data interface{}, err error) bool {
	cmd.results = append(cmd.results, Result{data, err})
	cmd.cur++
	return cmd.cur == cmd.max
}

func (cmd *BenchCommond) Flush() (res []Result) {
	if cmd.conn.Send(cmd) == nil {
		<-cmd.c
		res = cmd.results
	}
	cmd.Reset(cmd.conn, cmd.idx)
	return
}

func (cmd *BenchCommond) Reset(conn *Conn, idx int) {
	cmd.Commond.Reset(conn, idx)
	cmd.results = cmd.results[:0]
	cmd.cur, cmd.max = 0, 0
}

func (cmd *BenchCommond) String() string {
	return "[cmd] " + cmd.Cmd + ":" + cmd.Key
}

func (cmd *BenchCommond) Recycle() {
	RecycleBenchCmd(cmd)
}

// 异步消息
type SyncCommond struct {
	*Commond
	results []Result
	cur     int
	c       chan *Result
}

func newSyncCommond(conn *Conn, c chan *Result, idx int) (cmd *SyncCommond) {
	cmd = &SyncCommond{Commond: newCommond(conn, idx)}
	cmd.isAutoSend = false
	cmd.c = c
	return
}

func (cmd *SyncCommond) SetData(data interface{}, err error) bool {
	cmd.c <- &Result{data, err}
	return true
}

func (cmd *SyncCommond) Done() {
}

func (cmd *SyncCommond) Flush() error {
	return cmd.conn.Send(cmd)
}
func (cmd *SyncCommond) String() string {
	return "[cmd] " + cmd.Cmd + ":" + cmd.Key
}
