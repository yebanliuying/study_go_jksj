package main

import (
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"sync"
	"time"
)

func main() {

	client := goredislib.NewClient(&goredislib.Options{
		Addr: "139.198.21.42:6379",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)
	rs := redsync.New(pool)



	gNum := 2
	mutexname := "421"
	var wg sync.WaitGroup
	wg.Add(gNum)
	for i := 0; i< gNum; i++ {
		go func(){
			defer wg.Done()

			mutex := rs.NewMutex(mutexname)

			fmt.Println("开始获取锁")
			if err := mutex.Lock(); err != nil {
				panic(err)
			}

			fmt.Println("获取锁成功")
			time.Sleep(time.Second * 5)
			fmt.Println("开始释放锁")
			if ok, err := mutex.Unlock(); !ok || err != nil {
				panic("unlock failed")
			}
			fmt.Println("释放锁成功")

		}()
	}
	wg.Wait()


}