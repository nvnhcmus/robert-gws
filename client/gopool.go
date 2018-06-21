package main

import (
	"github.com/ivpusic/grpool"
	"fmt"
	"runtime"
)

func main() {

	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)

	// number of workers, and size of job queue
	pool := grpool.NewPool(10, 50)

	// release resources used by pool
	defer pool.Release()

	// submit one or more jobs to pool
	// how many jobs we should wait
	pool.WaitCount(100000)

	for i := 0; i < 100000; i++ {
		count := i

		pool.JobQueue <- func() {
			defer pool.JobDone()

			fmt.Printf("I am worker! Number %d\n", count)
			//x := 0
			//y := 1
			//try.This(func() {
			//	t := y / x
			//	log.Println(t)
			//}).Catch(func(e try.E) {
			//	log.Println("Broadcast message fail", e)
			//})
			//pool.JobDone()
		}
	}

	// dummy wait until jobs are finished
	//time.Sleep(1 * time.Second)

	//pool.WaitAll()
	// wait until we call JobDone for all jobs
	pool.WaitAll()

	//for i:=0; i< 10; i++ {
	//	x := 0
	//	y := 1
	//	try.This(func() {
	//		t := y / x
	//		log.Println(t)
	//	}).Catch(func(e try.E) {
	//		log.Println("Broadcast message fail", e)
	//	})
	//}
}