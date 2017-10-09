// The Isomorphic Go Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package isokit

import (
	"context"
	"regexp"
	"strings"
)

const (
	RouteWithParamsPattern = `/([^/]*)`
	RouteOnlyPrefixPattern = `/`
	RouteSuffixPattern     = `/?$`
)

type Handler interface {
	ServeRoute(context.Context)
}

type HandlerFunc func(context.Context)

func (f HandlerFunc) ServeRoute(ctx context.Context) {
	f(ctx)
}

type RouteVarsKey string

type Route struct {
	handler  Handler
	regex    *regexp.Regexp
	varNames []string
}

func NewRoute(path string, handler HandlerFunc) *Route {
	r := &Route{
		handler: handler,
	}

	path = strings.TrimPrefix(path, `/`)
	if strings.HasSuffix(path, `/`) {
		path = strings.TrimSuffix(path, `/`)
	}

	routeParts := strings.Split(path, "/")

	var routePattern string = `^`
	for _, routePart := range routeParts {
		if strings.HasPrefix(routePart, `{`) && strings.HasSuffix(routePart, `}`) {
			routePattern += RouteWithParamsPattern
			routePart = strings.TrimPrefix(path, `{`)
			routePart = strings.TrimSuffix(path, `}`)
			r.varNames = append(r.varNames, routePart)
		} else {
			routePattern += RouteOnlyPrefixPattern + routePart
		}
	}
	routePattern += RouteSuffixPattern
	r.regex = regexp.MustCompile(routePattern)
	return r
}
