// ---------------------------------------------------------------------------------------------- //
// -- Copyright (c) 2024 Braden Hitchcock - MIT License (https://opensource.org/licenses/MIT)  -- //
// ---------------------------------------------------------------------------------------------- //
package server

import "net/http"

type Middleware interface {
	Middleware(http.HandlerFunc) http.HandlerFunc
}

type MiddlewareChain struct {
	links []Middleware
}

func NewMiddlewareChain(hs ...Middleware) MiddlewareChain {
	return MiddlewareChain{
		links: append(make([]Middleware, 0, len(hs)), hs...),
	}
}

func (c *MiddlewareChain) Then(n Middleware) *MiddlewareChain {
	c.links = append(c.links, n)
	return c
}

func (c *MiddlewareChain) Finish(f http.HandlerFunc) http.HandlerFunc {
	return buildMiddlewareHandler(f, c.links...)
}

func buildMiddlewareHandler(f http.HandlerFunc, hs ...Middleware) http.HandlerFunc {
	if len(hs) == 0 {
		return f
	}
	return hs[0].Middleware(buildMiddlewareHandler(f, hs[1:]...))
}
