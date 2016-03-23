package sredis

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	glog "github.com/scabwort/glog"
)

const (
	defaultLinkTimeout = time.Second * 5

	CMD_auth   = "AUTH"
	CMD_echo   = "ECHO"
	CMD_ping   = "PING"
	CMD_quit   = "QUIT"
	CMD_select = "SELECT"
)

type ProtocolError string

func (e ProtocolError) Error() string {
	return string(e)
}

type NetStatus int

const (
	statusReady NetStatus = iota
	statusLinking
	statusOk
	statusClose
)
const (
	recvPacket = 0
	sendPacket = 0
	closeConn  = 1
)

var (
	ErrRedisConClose  = errors.New("Redis conn is closed")
	ErrRedisConFull   = errors.New("Redis conn query is full")
	ErrRedisNetFatal  = errors.New("The connection with redis is fatal")
	ErrInvalidAddress = errors.New("Redis tcp address is invalid")
)

type Conn struct {
	conn        *net.TCPConn
	address     string
	auth        string
	dbindex     string
	isclose     int32
	signal      chan int
	isreconnect int32
	sendQueue   []ICommond
	sendMutex   sync.Mutex
	recvQueue   []ICommond
	recvMutex   sync.Mutex
	maxqueue    int
}

func NewConn(address, auth, dbidx string) (conn *Conn, err error) {
	conn = &Conn{address: address}
	conn.signal = make(chan int, 1)
	conn.maxqueue = 50000
	conn.auth = auth
	conn.dbindex = dbidx
	err = conn.Connect()
	for err != nil {
		glog.Error("[Redis] connect fail, try reconnect...")
		time.Sleep(time.Second)
		err = conn.Connect()
	}

	return
}

func (c *Conn) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, defaultLinkTimeout)
	if err == nil {
		glog.Info("[Redis] 连接成功,", c.address)
		c.conn = conn.(*net.TCPConn)
		return c.handshake()
	}
	return err
}

func (c *Conn) ReConnect() (err error) {
	if c.IsClosed() {
		if atomic.CompareAndSwapInt32(&c.isreconnect, 0, 1) {
			glog.Error("[Redis] ", c.address, ", 重连中...")
			err = c.Connect()
			for err != nil {
				glog.Error("[Redis] ", c.address, ", 重连中...")
				time.Sleep(time.Second)
				err = c.Connect()
			}
			glog.Info("[Redis] 重连成功,", c.address)
		}
	}
	return
}

// handshake include auth and select db index
func (c *Conn) handshake() error {
	if c.auth != "" {
		buf := NewPacketSize(128)
		buf.WriteCmd(CMD_auth, 2)
		buf.WriteString(&c.auth)
		if err := buf.Flush(c.conn); err != nil {
			glog.Error("[Redis] auth error!", c.address, ", err:", c.auth)
			c.Close(true)
			return err
		}
		_, err, cerr := c.readReply(buf)
		if cerr != nil || err != nil {
			glog.Error("[Redis] auth error!", c.address, ", err:", c.auth)
			c.Close(true)
			if cerr != nil {
				return cerr
			}
			return err
		}
	}
	if c.dbindex != "" && c.dbindex != "0" {
		buf := NewPacketSize(128)
		buf.WriteCmd(CMD_select, 2)
		buf.WriteString(&c.dbindex)
		if err := buf.Flush(c.conn); err != nil {
			glog.Error("[Redis] select db index error!", c.address, ", err:", c.dbindex)
			c.Close(true)
			return err
		}
		_, err, cerr := c.readReply(buf)
		if cerr != nil || err != nil {
			glog.Error("[Redis] select db index error!", c.address, ", err:", c.dbindex)
			c.Close(true)
			if cerr != nil {
				return cerr
			}
			return err
		}
	}
	atomic.StoreInt32(&c.isclose, 0)
	atomic.StoreInt32(&c.isreconnect, 0)
	go c.loop()
	c.signalSend()
	glog.Info("[Redis] 认证通过", c.dbindex, c.auth)
	return nil
}

func (c *Conn) Send(cmd ICommond) (err error) {
	// 如果关闭，并且没有在重连，则返回错误
	if c.IsClosed() {
		cmd.SetData(nil, ErrRedisConClose)
		return ErrRedisConClose
	}
	// 添加到发送队列
	c.sendMutex.Lock()
	if len(c.sendQueue) > c.maxqueue {
		c.sendMutex.Unlock()
		cmd.SetData(nil, ErrRedisConFull)
		return ErrRedisConFull
	}
	c.sendQueue = append(c.sendQueue, cmd)
	c.sendMutex.Unlock()

	c.signalSend()
	return
}

func (c *Conn) signalSend() {
	select {
	case c.signal <- sendPacket:
	default:
	}
}

func (c *Conn) IsClosed() bool {
	return atomic.LoadInt32(&c.isclose) == 1
}

func (c *Conn) IsReconnected() bool {
	return atomic.LoadInt32(&c.isreconnect) == 1
}

func (c *Conn) Close(force bool) {
	if atomic.CompareAndSwapInt32(&c.isclose, 0, 1) || force {
		c.conn.Close()
		glog.Error("[Redis] conn be close")
		select {
		case c.signal <- closeConn:
		default:
		}
		c.recvMutex.Lock()
		c.FreeCmds(c.recvQueue)
		c.recvQueue = c.recvQueue[:0]
		c.recvMutex.Unlock()
	}
}

func (c *Conn) beClosed() {
	c.Close(false)
	go c.ReConnect()
}

func (c *Conn) loop() {
	go c.onRecv()
	var (
		buf    = NewPacketSize(8192)
		err    error
		cmd    ICommond
		idx    int
		cmds   []ICommond
		maxCmd int
		// ping per minute when no commond
		pingPack = newCommond(c, 0)
		pingTick = time.NewTicker(time.Minute)
		isPing   bool
	)
	defer pingTick.Stop()
	pingPack.buf.WriteCmd(CMD_ping, 1)
	for {
		select {
		case r := <-c.signal:
			if r != sendPacket {
				c.Close(false)
				return
			}
			for {
				if c.IsClosed() {
					glog.Info("[Redis] conn is closed")
					return
				}
				maxCmd = 64
				c.sendMutex.Lock()
				if len(c.sendQueue) > 0 {
					if len(c.sendQueue) < maxCmd {
						maxCmd = len(c.sendQueue)
					}
					cmds = c.sendQueue[:maxCmd]
					c.sendQueue = c.sendQueue[maxCmd:]
				}
				c.sendMutex.Unlock()

				if len(cmds) > 0 {
					for idx, cmd = range cmds {
						err = buf.Bytes(cmd.GetBytes())
						if err != nil {
							glog.Error("[Redis] send msg get bytes err:", err)
							cmd.SetData(nil, err)
							cmd.Done()
							cmds = append(cmds[:idx], cmds[idx+1:]...)
						}
					}
					if err = buf.Flush(c.conn); err != nil {
						c.FreeCmds(cmds)
						c.beClosed()
						return
					} else {
						c.recvMutex.Lock()
						if !c.IsClosed() {
							c.recvQueue = append(c.recvQueue, cmds...)
						} else {
							c.FreeCmds(cmds)
						}
						c.recvMutex.Unlock()
					}
					cmds = cmds[:0]
					isPing = true
				} else {
					break
				}
			}
		case <-pingTick.C:
			if !isPing {
				c.Send(pingPack)
			}
			isPing = false
		}
	}
}

func (c *Conn) FreeCmds(cmds []ICommond) {
	if cmds == nil || len(cmds) == 0 {
		return
	}
	for _, cmd := range cmds {
		cmd.SetData(nil, ErrRedisConClose)
		cmd.Done()
	}
}

func (c *Conn) onRecv() {
	var (
		buf  = NewPacketSize(8 * 1024)
		cmd  ICommond
		err  error
		cerr error
		data interface{}
	)
	for {
		data, err, cerr = c.readReply(buf)
		if cerr != nil {
			c.beClosed()
			return
		}
		c.recvMutex.Lock()
		if len(c.recvQueue) > 0 {
			cmd = c.recvQueue[0]
			c.recvQueue = c.recvQueue[1:]
		}
		c.recvMutex.Unlock()
		if cmd == nil {
			continue
		}
		if cmd.SetData(data, err) {
			cmd.Done()
		}
		buf.ReFill()
	}
}

func (c *Conn) readReply(buf *Packet) (data interface{}, err error, closeerror error) {
	var line []byte
	var n int
	if line, err = buf.ReadLine(c.conn); err != nil {
		glog.Error("[Redis] read conn error:", err)
		closeerror = err
		return
	}
	if line == nil || len(line) == 0 {
		return nil, errors.New("commond replay is nil"), nil
	}
	switch line[0] {
	case ':', '+':
		data = line[1:]
	case '-':
		err = errors.New(string(line[1:]))
	case '$':
		n, err = parseToLen(line[1:])
		if n < 0 || err != nil {
			return
		}
		p := make([]byte, n+2)
		if err = buf.ReadBytes(c.conn, p, true); err != nil {
			return
		}
		if p[n] != '\r' || p[n+1] != '\n' {
			return nil, errors.New("replay format is error"), nil
		}
		data = p[:n]
	case '*':
		n, err = parseToLen(line[1:])
		if n < 0 || err != nil {
			return
		}
		r := make([]interface{}, n)
		for i := range r {
			r[i], err, closeerror = c.readReply(buf)
			if err != nil {
				return nil, err, closeerror
			}
		}
		data = r
	}
	return
}

// parseLen parses bulk string and array lengths.
func parseToLen(p []byte) (int, error) {
	if len(p) == 0 {
		return -1, ProtocolError("malformed length")
	}

	if p[0] == '-' && len(p) == 2 && p[1] == '1' {
		// handle $-1 and $-1 null replies.
		return -1, nil
	}

	var n int
	for _, b := range p {
		n *= 10
		if b < '0' || b > '9' {
			return -1, ProtocolError("illegal bytes in length")
		}
		n += int(b - '0')
	}

	return n, nil
}
