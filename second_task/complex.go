package main

import (
	"iter"
	"reflect"
	"slices"
)

// В общем, идея была в том, чтобы не запускать отдельную горутину
// под каждый канал, за пускать сильно меньше горутин с очень большим
// select (собираемым в рантайме через рефлект).

// Чтобы обойти ограничение 65_536 кейса в селекте, при заполнении селекта
// до максимума, в последний кейс кладется канал в который будут писаться
// все остальные горутины (функция вызывется рекурсивно).

// Но оказалось, что это работает сильно хуже, чем просто создать много горутин.
// Возможно, select не оптимизирован для работы с большим количесвом кейсов,
// а может проблема в коде.

// https://pkg.go.dev/reflect#Select panic if len(cases) > 65536
const maxCaseInSelect = 65_536

func hugeMerge[T any](chs ...<-chan T) <-chan T {

	if len(chs) <= maxCaseInSelect {
		return bunch(chs...)
	}

	next, stop := iter.Pull(slices.Chunk(chs, maxCaseInSelect))
	defer stop()

	return bunchOfBunch(next)

}

func bunch[T any](chs ...<-chan T) <-chan T {
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

func bunchOfBunch[T any](next func() ([]<-chan T, bool)) <-chan T {

	outs := make([]<-chan T, 0)

	var ok bool
	for range maxCaseInSelect - 1 {
		var v []<-chan T
		if v, ok = next(); !ok {
			return bunch(outs...)
		}
		outs = append(outs, bunch(v...))
	}
	if ok {
		outs = append(outs, bunchOfBunch(next))
	}
	return bunch(outs...)
}
