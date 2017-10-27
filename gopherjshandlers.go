// +build !clientonly

package isokit

import (
	"net/http"
)

func GopherjsScriptHandler(webAppRoot string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webAppRoot+"/client/client.js")
	})
}

func GopherjsScriptMapHandler(webAppRoot string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, webAppRoot+"/client/client.js.map")
	})
}
