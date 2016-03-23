package sredis

import (
	"fmt"
)

var (
	conn *Conn
)

func init() {
	var err error
	conn, err = NewConn("192.168.124.130:6379", "ztgame123654", "")
	if err != nil {
		fmt.Println("link error:", err)
	}
}

//func TestChan(t *testing.T) {
//	b := make(chan ICommond, 10)
//	go func() {
//		for {
//			cmd := NewBenchCommond(nil)
//			b <- cmd
//			cmd.Set("pse", "srere")
//			time.Sleep(time.Second)
//		}
//	}()
//	time.Sleep(time.Second)

//	for {
//		cmd := <-b
//		fmt.Println(len(cmd.GetBytes()))
//		time.Sleep(time.Second)
//	}
//}
