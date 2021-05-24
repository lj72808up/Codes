package test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)



var runner bool = false
var wg sync.WaitGroup

func workFunction(x int32) {
	fmt.Println("deal with", x)
	time.Sleep(5*time.Second)
}

func worker(wg *sync.WaitGroup, ch chan int32) {
	for d := range ch {
		workFunction(d)
	}
}

func createWorkerPool(noOfWorkers int, ch chan int32) {
	for i := 0; i < noOfWorkers; i++ {
		wg.Add(1)
		go worker(&wg, ch)
	}
	wg.Done()
}

func produce(c chan int32) {
	for runner == false {
		num := rand.Int31n(3)
		time.Sleep(time.Duration(num) * time.Second)
		c <- num
		fmt.Println(fmt.Sprintf("write channel c %d", num))
	}
}

func TestChan2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	c := make(chan int32, 5)  // channel
	createWorkerPool(2, c)   // 消费
	produce(c)    // 生产
}
