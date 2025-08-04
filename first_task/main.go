package main

import (
	"fmt"
	"os"
)

// PrintSorted принимает на вход два канала, каждый из которых возвращает конечную монотонно неубывающую
// последовательность целых чисел (т.е. отсортированные по возрастанию). Необходимо объединить значения
// из обоих каналов и вывести их в stdout в виде одной монотонно неубывающей последовательности.
//
// Пример:
// In: a = [0, 0, 3, 4]; b = [1, 3, 4, 6, 8]
// Out: res = [0, 0, 1, 3, 3, 4, 4, 6, 8]
//
// Для проверки решения запустите тесты: go test -v
func PrintSorted(ch1, ch2 <-chan int) {

	if os.Getenv("AMO_GO_TASK_MODE") == "simple" {
		printSorted(ch1, ch2) // самое простое решение
	} else {
		for v := range mergeChs([]<-chan int{ch1, ch2}) {
			fmt.Println(v)
		}
	}
}

func main() {
	a := make(chan int)
	b := make(chan int)

	go func() {
		defer close(a)
		a <- 1
		a <- 4
		a <- 6
	}()

	go func() {
		defer close(b)
		b <- 2
		b <- 3
		b <- 5
		b <- 7
		b <- 8
		b <- 9
	}()

	PrintSorted(a, b)
}
