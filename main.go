package main

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"sync"
	"time"
)

func main() {
	forTo(9)
	log.Println("------ end ------")
	// 等待超时任务的最终输出
	time.Sleep(2 * time.Second)
}

type Wf func() (data interface{}, err error)

func timeOut(to time.Duration, f Wf) (data interface{}, err error) {
	dc := make(chan interface{})
	go func() {
		defer close(dc)
		var da interface{}
		da, err = f()
		dc <- da
	}()
	select {
	case data = <-dc:
	case <-time.After(to):
		err = fmt.Errorf("timeout %s", to)
	}
	return
}

func forTo(ln int) {
	m := make(map[string]string)
	var mut sync.Mutex
	var emut sync.Mutex
	// 错误汇总
	var glerr error
	to := 2 * time.Second
	var wt sync.WaitGroup
	for i := 0; i < ln; i++ {
		wt.Add(1)
		go func() {
			_, err := timeOut(to, func() (data interface{}, err error) {
				var ss []string
				ss, err = exe()
				log.Println("实际结果", ss, err)
				mut.Lock()
				defer mut.Unlock()
				if err != nil {
					return
				}
				data = ss
				for _, s := range ss {
					m[s] = s
				}
				return
			})
			emut.Lock()
			defer emut.Unlock()
			if err != nil {
				if glerr == nil {
					glerr = err
				} else {
					glerr = errors.Join(glerr, err)
				}
			}
			wt.Done()
		}()
	}
	wt.Wait()
	if glerr != nil {
		glerr = errors.New(strings.ReplaceAll(
			glerr.Error(), "\n", ","))
	}
	log.Println(m, glerr)
}

func exe() ([]string, error) {
	i := rand.Intn(3000)
	var err = fmt.Errorf("模拟错误")
	time.Sleep(time.Duration(i) * time.Millisecond)
	if i < 1000 {
		return []string{"1", "a"}, nil
	} else if i < 1600 {
		return []string{}, err
	} else if i < 2100 {
		return []string{"a", "b", "c"}, nil
	} else {
		return []string{}, err
	}
}
