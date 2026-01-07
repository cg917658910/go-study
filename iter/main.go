package main

import (
	"cmp"
	"iter"
	"maps"
	"slices"
)

type Iterator[V any] struct {
	iter iter.Seq[V]
}

func From[V any](slice []V) *Iterator[V] {
	return &Iterator[V]{
		iter: func(yield func(V) bool) {
			for _, v := range slice {
				if !yield(v) {
					return
				}
			}
		},
	}
}

func (i *Iterator[V]) Each(f func(V)) {
	i.iter(func(v V) bool {
		f(v)
		return true
	})
}

func FromSlice[T any](s []T) iter.Seq[T] {
	return slices.Values(s)
}

type User struct {
	Name string
}

func main() {
	// Example usage
	slice := []int{1, 2, 3, 4, 5}
	From(slice).Each(func(v int) {
		println(v)
	})
	users := []User{
		{Name: "cg"},
		{Name: "fc"},
	}
	From(users).Each(func(u User) { println(u.Name) })
	/* 	iter.Each(func(v int) {
		println(v)
	}) */
	slices.Values(slice)(func(v int) bool {
		println(v)
		return true
	})
	// Convert slice to iterator and back
	/* newSlice := FromSlice(slice)
	newSlice(func(v int) bool {
		println(v)
		return true
	}) */
	 cmp.Compare()
	 maps.Values()
}
