package main

import (
	"fmt"
	"slices"
	"testing"

	"github.com/nalgeon/be"
)

func TestHugeMerge(t *testing.T) {

	tcs := []struct {
		num, size int // смотри aLotNumGen()
	}{
		{num: 3, size: 3},
		{num: 20, size: 5},
		//{num: 70_000, size: 5}, // этот test case будет выполняться ОЧЕНЬ долго.
	}

	for _, tc := range tcs {
		var numGens []<-chan int = aLotNumGen(tc.num, tc.size)

		merged := hugeMerge(numGens...)
		nums := readAllAndSort(merged)

		be.Equal(t, numSeq(tc.num*tc.size), nums)
	}

}

// возвращает слайс каналов генераторов интов
// aLotNumGen(3, 4) вернет 3 канала, которые могут быть вычитаны как:
// [0, 1, 2, 3] , [4, 5, 6, 7] , [8, 9, 10, 11]
func aLotNumGen(num, size int) []<-chan int {
	out := make([]<-chan int, num)

	for i := range num {
		out[i] = numGen(i*size, (i+1)*size)
	}

	return out
}

// возвращает слайс интов: [0, 1, 2, 3, ... , len-1]
func numSeq(len int) []int {
	out := make([]int, len)
	for i := range len {
		out[i] = i
	}
	return out
}

func readAllAndSort(nums <-chan int) []int {
	var out []int
	for v := range nums {
		out = append(out, v)
	}
	slices.Sort(out)
	return out
}

func numGen(begin, end int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)
		for i := begin; i < end; i++ {
			out <- i
		}
	}()

	return out
}
