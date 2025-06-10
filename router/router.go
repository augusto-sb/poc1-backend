package router

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"

	"example.com/auth"
	"example.com/entity"
)

type Route struct {
	hndlr    http.HandlerFunc
	authRole string // *string
}

type RouteMap map[string] /*http.Method*/ Route

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

func GlobalMiddleware(rm RouteMap) http.HandlerFunc {
	fmt.Println(rm)
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		valor, existe := rm[req.Method]
		if existe {
			if valor.authRole == "" /*nil*/ {
				corsMiddleware(printNextFuncName(valor.hndlr)).ServeHTTP(rw, req)
			} else {
				corsMiddleware(auth.Middleware(printNextFuncName(valor.hndlr), valor.authRole)).ServeHTTP(rw, req)
			}
		} else {
			var allowedMethods []string = []string{}
			for key := range rm {
				allowedMethods = append(allowedMethods, key)
			}
			rw.Header().Set("Allow", strings.Join(allowedMethods, ","))
			rw.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
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
	contextPath = os.Getenv("CONTEXT_PATH") // standarizar
	if contextPath == "" {
		contextPath = "/backend"
	}
	Mux = http.NewServeMux()
	//Mux.HandleFunc("OPTIONS /", corsMiddleware(printNextFuncName(corsHandler)))
	Mux.HandleFunc("/", corsMiddleware(printNextFuncName(http.NotFound)))
	Mux.HandleFunc(contextPath+"/entities", GlobalMiddleware(RouteMap{
		http.MethodGet: Route{
			hndlr:    entity.GetEntities,
			authRole: "entity-read",
		},
		http.MethodPost: Route{
			hndlr:    entity.AddEntity,
			authRole: "entity-create",
		},
	}))
	Mux.HandleFunc(contextPath+"/entities/{id}", GlobalMiddleware(RouteMap{
		http.MethodGet: Route{
			hndlr:    entity.GetEntity,
			authRole: "entity-read",
		},
		http.MethodDelete: Route{
			hndlr:    entity.RemoveEntity,
			authRole: "entity-delete",
		},
		http.MethodPut: Route{
			hndlr:    entity.UpdateEntity,
			authRole: "entity-update",
		},
	}))
}
