package main

import "sync"

func merge(ch1, ch2 <-chan int) <-chan int {
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

