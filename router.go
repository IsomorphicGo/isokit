// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import (
	"context"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/dom"
)

type Router struct {
	routes   []*Route
	listener func(*js.Object)
}

func NewRouter() *Router {

	initializeHistoryInteractions()

	return &Router{
		routes: []*Route{},
	}
}

func (r *Router) Handle(path string, handler Handler) *Route {
	return r.HandleFunc(path, handler.(HandlerFunc))
}

func (r *Router) HandleFunc(path string, handler HandlerFunc) *Route {
	route := NewRoute(path, handler)
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) Listen() {
	r.RegisterLinks("body a")
}

func (r *Router) RegisterLinks(querySelector string) {
	document := dom.GetWindow().Document().(dom.HTMLDocument)
	links := document.QuerySelectorAll(querySelector)

	for _, link := range links {

		href := link.GetAttribute("href")
		switch {

		case strings.HasPrefix(href, "/") && !strings.HasPrefix(href, "//"):

			if r.listener != nil {
				link.RemoveEventListener("click", false, r.listener)
			}

			r.listener = link.AddEventListener("click", false, r.linkHandler)
		}
	}

}

func (r *Router) linkHandler(event dom.Event) {

	uri := event.CurrentTarget().GetAttribute("href")
	path := strings.Split(uri, "?")[0]
	//	leastParams := -1
	var matchedRoute *Route
	var parts []string
	var lowestMatchCountSet bool = false
	var lowestMatchCount int = -1

	for _, route := range r.routes {

		matches := route.regex.FindStringSubmatch(path)
		matchesExist := len(matches) > 0 && matches != nil
		isLowestMatchCount := (lowestMatchCountSet == false) || (len(matches) < lowestMatchCount)

		if matchesExist && isLowestMatchCount {
			matchedRoute = route
			parts = matches[1:]
			lowestMatchCount = len(matches)
			lowestMatchCountSet = true
		}
	}

	if matchedRoute != nil {
		event.PreventDefault()
		js.Global.Get("history").Call("pushState", nil, "", uri)
		routeVars := make(map[string]string)

		for i, part := range parts {
			routeVars[matchedRoute.varNames[i]+`}`] = part
		}

		var ctx context.Context
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		k := RouteVarsKey("Vars")
		ctx = context.WithValue(ctx, k, routeVars)
		go matchedRoute.handler.ServeRoute(ctx)
	}
}

func initializeHistoryInteractions() {
	// Handler for back/forward button interactions
	dom.GetWindow().AddEventListener("popstate", false, func(event dom.Event) {
		var location = js.Global.Get("location")
		js.Global.Set("location", location)
	})
}
