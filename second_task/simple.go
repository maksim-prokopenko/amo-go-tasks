package main

import (
	"sync"
)

func merge2(ch1, ch2 <-chan int) <-chan int {
	//return merge(ch1, ch2)
	out := make(chan int)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		wg.Wait()
		close(out)
	}()

	go func() {
		defer wg.Done()
		for v := range ch1 {
			out <- v
		}
	}()

	go func() {
		defer wg.Done()
		for v := range ch2 {
			out <- v
		}
	}()

	return out
}

func merge(chs ...<-chan int) <-chan int {
	out := make(chan int)

	var wg sync.WaitGroup
	wg.Add(len(chs))

	go func() {
		wg.Wait()
		close(out)
	}()

	for _, ch := range chs {
		go func() {
			defer wg.Done()
			for v := range ch {
				out <- v
			}
		}()
	}

	return out
}
