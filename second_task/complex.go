package main

import (
	"reflect"
)

// одна проблема – надо как-то обрабатыать https://pkg.go.dev/reflect#Select panic len(cases) > 65536
// можно выдумать какую-нибудь рекурсию

func minimalMerge[T any](chs ...<-chan T) <-chan T {
	out := make(chan T)

	go func() {
		defer close(out)
		cases := make([]reflect.SelectCase, len(chs))
		for i, c := range chs {
			cases[i].Dir = reflect.SelectRecv
			cases[i].Chan = reflect.ValueOf(c)
		}
		for len(cases) > 0 {
			i, v, ok := reflect.Select(cases)
			if !ok {
				cases = remove(cases, i)
				continue
			}
			out <- v.Interface().(T)
		}
	}()

	return out
}

func remove[S ~[]T, T any](sl S, i int) S {
	sl[i] = sl[len(sl)-1]
	return sl[:len(sl)-1]
}
