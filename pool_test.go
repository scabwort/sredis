package sredis

import (
	"fmt"
	"strconv"
	"testing"
)

type SnsDynamicMsg struct {
	LoveNum      uint32
	FollowNum    uint32
	LeaveMsgNum  uint32
	StarPhoneNum uint32
}

var i int

func init() {
	err := RegRedisKey(1, "@tcp(192.168.124.130:6379)/0")
	if err != nil {
		fmt.Println(err)
	}
}

func TestPool(t *testing.T) {

	op := NewRedisOp(1)
	t.Log(op.SetObject("testaaa", SnsDynamicMsg{}))
	//	t.Log(time.Now().String())
	//	p := sync.WaitGroup{}

	//	for i := 0; i < 100; i++ {
	//		p.Add(1)
	//		go func() {
	//			bcmd := NewBenchCmd(1)
	//			var key string
	//			for j := 0; j < 10; j++ {
	//				for k := 0; k < 100; k++ {
	//					key = "test" + strconv.Itoa(k)
	//					bcmd.Set(key, key)
	//				}
	//				bcmd.Flush()
	//			}
	//			bcmd.Recycle()
	//			p.Done()
	//		}()
	//	}
	//	p.Wait()

	//	t.Log(time.Now().String())
	//	cmd := NewCmd(1)
	//	//	t.Log(cmd.conn)
	//	s, err := cmd.Get("test").String()
	//	if err != nil {
	//		t.Log(err)
	//	}
	//	t.Log(s)
	//	cmd.Recycle()
}

func BenchmarkPoolCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		cmd := NewCmd(1)
		cmd.Set("test1", "test1")
		cmd.Recycle()
	}
}

func BenchmarkPoolBenchCmd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bcmd := NewBenchCmd(1)
		var key string
		for i := 0; i < 100; i++ {
			key = "test" + strconv.Itoa(i)
			bcmd.Set(key, key)
		}
		bcmd.Flush()
		bcmd.Recycle()
	}
}
