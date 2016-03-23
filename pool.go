package sredis

import (
	"strings"
	"sync"

	glog "github.com/scabwort/glog"
)

type RedisPool struct {
	connPools        map[int]*Conn
	cmdpool          []*sync.Pool
	cmdpoolSize      int
	cmdIdx           int
	benchCmdpool     []*sync.Pool
	benchCmdpoolSize int
	benchCmdIdx      int
}

var rpool *RedisPool

func init() {
	rpool = NewPool()
}

func CloseAll() {
	if rpool == nil || len(rpool.connPools) == 0 {
		return
	}
	for _, p := range rpool.connPools {
		p.Close(false)
	}
}

func NewPool() (p *RedisPool) {
	p = new(RedisPool)
	p.connPools = make(map[int]*Conn)
	p.cmdpoolSize = 1024
	p.benchCmdpoolSize = 1024
	p.cmdpool = make([]*sync.Pool, p.cmdpoolSize)
	p.benchCmdpool = make([]*sync.Pool, p.benchCmdpoolSize)
	for i := 0; i < p.cmdpoolSize; i++ {
		p.cmdpool[i] = new(sync.Pool)
	}
	for i := 0; i < p.benchCmdpoolSize; i++ {
		p.benchCmdpool[i] = new(sync.Pool)
	}
	return
}

func RegRedisKey(key int, address string) error {
	server, pwd, dbKey := ParseRedisKey(address)
	if server == "" {
		return ErrRedisConClose
	}
	conn, err := NewConn(server, pwd, dbKey)
	if err != nil {
		return err
	}
	rpool.connPools[key] = conn
	glog.Info("[Redis] 连接并注册成功:", address, ",key:", key)
	return nil
}

// ztgame123654@tcp(192.168.124.130:6379)/1
func ParseRedisKey(address string) (server, pwd, dbKey string) {
	if strings.LastIndex(address, "(") == -1 || strings.LastIndex(address, ")") == -1 || strings.LastIndex(address, "@") == -1 || strings.LastIndex(address, "/") == -1 {
		if strings.LastIndex(address, ":") == -1 {
			glog.Error("[Redis] read config error!", address)
			return
		}
		server = address
		return
	}
	server = address[strings.LastIndex(address, "(")+1 : strings.LastIndex(address, ")")]
	pwd = address[0:strings.LastIndex(address, "@")]
	dbKey = address[strings.LastIndex(address, "/")+1:]
	return
}

func Get(key int) *Conn {
	return rpool.connPools[key]
}

// 从对象池中建立命令
func NewCmd(key int) *Commond {
	conn := Get(key)
	rpool.cmdIdx++
	idx := rpool.cmdIdx % rpool.cmdpoolSize
	p := rpool.cmdpool[idx]
	if v := p.Get(); v != nil {
		cmd := v.(*Commond)
		cmd.Reset(conn, idx)
		return cmd
	}
	return newCommond(conn, idx)
}

// 回收到对象池中
func RecycleCmd(c *Commond) {
	p := rpool.cmdpool[c.idx]
	p.Put(c)
}

// 从对象池中建立批量命令
func NewBenchCmd(key int) *BenchCommond {
	conn := Get(key)
	//	rpool.benchCmdIdx++
	//	idx := rpool.benchCmdIdx % rpool.benchCmdpoolSize
	//	p := rpool.benchCmdpool[idx]
	//	if v := p.Get(); v != nil {
	//		cmd := v.(*BenchCommond)
	//		cmd.Reset(conn, idx)
	//		return cmd
	//	}
	return NewBenchCommond(conn, 0)
}

// 回收到对象池中
func RecycleBenchCmd(c *BenchCommond) {
	p := rpool.benchCmdpool[c.idx]
	p.Put(c)
}
