package router

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"

	"example.com/auth"
	"example.com/entity"
)

var Mux *http.ServeMux

func printNextFuncName(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println("invoked: ", runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name())
		next.ServeHTTP(rw, req)
	})
}

var corsMiddleware func(http.HandlerFunc) http.HandlerFunc = func(next http.HandlerFunc) http.HandlerFunc {
	return next
}

func corsHandler(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(200)
}

func init() {
	var corsOrigin = os.Getenv("CORS_ORIGIN")
	if corsOrigin != "" {
		corsMiddleware = func(next http.HandlerFunc) http.HandlerFunc {
			return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				rw.Header().Set("Access-Control-Allow-Origin", corsOrigin)
				rw.Header().Set("Access-Control-Allow-Headers", "authorization,content-type")
				rw.Header().Set("Access-Control-Allow-Methods", "GET,DELETE,POST,PUT")
				next.ServeHTTP(rw, req)
			})
		}
	}
	var contextPath string
	contextPath = os.Getenv("CONTEXT_PATH")
	if contextPath == "" {
		contextPath = "/backend"
	}
	Mux = http.NewServeMux()
	Mux.HandleFunc("/", corsMiddleware(printNextFuncName(http.NotFound)))
	Mux.HandleFunc("OPTIONS /", corsMiddleware(printNextFuncName(corsHandler)))
	Mux.HandleFunc("GET "+contextPath+"/entities", corsMiddleware(auth.Middleware(printNextFuncName(entity.GetEntities), "entity-read")))
	Mux.HandleFunc("GET "+contextPath+"/entities/{id}", corsMiddleware(auth.Middleware(printNextFuncName(entity.GetEntity), "entity-read")))
	Mux.HandleFunc("POST "+contextPath+"/entities", corsMiddleware(auth.Middleware(printNextFuncName(entity.AddEntity), "entity-create")))
	Mux.HandleFunc("DELETE "+contextPath+"/entities/{id}", corsMiddleware(auth.Middleware(printNextFuncName(entity.RemoveEntity), "entity-delete")))
	Mux.HandleFunc("PUT "+contextPath+"/entities/{id}", corsMiddleware(auth.Middleware(printNextFuncName(entity.UpdateEntity), "entity-update")))
}
