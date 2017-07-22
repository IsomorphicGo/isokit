package isokit

import (
	"net/http"
	"os"
)

func GopherjsScriptHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, os.Getenv("ISOGO_APP_ROOT")+"/client/client.js")
}

func GopherjsScriptMapHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, os.Getenv("ISOGO_APP_ROOT")+"/client/client.js.map")
}
