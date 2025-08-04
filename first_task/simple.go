package main

import "fmt"

func printSorted(ch1, ch2 <-chan int) {
	v1, ok1 := <-ch1
	v2, ok2 := <-ch2
	for ok1 && ok2 {
		if v1 < v2 {
			fmt.Println(v1)
			v1, ok1 = <-ch1
		} else {
			fmt.Println(v2)
			v2, ok2 = <-ch2
		}
	}
	for ok1 {
		fmt.Println(v1)
		v1, ok1 = <-ch1
	}
	for ok2 {
		fmt.Println(v2)
		v2, ok2 = <-ch2
	}
}
