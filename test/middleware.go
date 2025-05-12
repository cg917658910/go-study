package test

import (
	"context"
	"fmt"
)

type HandlerFunc func(ctx context.Context) error

type MiddlewareFunc func(HandlerFunc) HandlerFunc

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}

func applyHandler(handler HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	return func(ctx context.Context) error {
		h := applyMiddleware(handler, middleware...)
		return h(ctx)
	}
}
func handler() HandlerFunc {
	return func(ctx context.Context) error {
		fmt.Println("I am handler")
		return nil
	}
}

var m1 = func(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context) error {
		fmt.Println("M1 is runing...")
		err := next(ctx)
		if err != nil {
			fmt.Println("M1 error:", err)
		}
		fmt.Println("M1 is end...")
		return err
	}
}

var m2 = func(next HandlerFunc) HandlerFunc {
	return func(ctx context.Context) error {
		fmt.Println("M2 is running...")
		err := next(ctx)
		if err != nil {
			fmt.Println("M2 error:", err)
		}
		fmt.Println("M2 is end...")
		return err
	}
}

func testMiddleware() {
	applyHandler(handler(), m1, m2)(context.Background())
}
