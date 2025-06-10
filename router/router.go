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
	var allowedMethods []string = []string{}
	for key := range rm {
		allowedMethods = append(allowedMethods, key)
	}
	allowedMethodsStr := strings.Join(allowedMethods, ",")
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		valor, existe := rm[req.Method]
		if existe {
			if valor.authRole == "" /*nil*/ {
				corsMiddleware(printNextFuncName(valor.hndlr)).ServeHTTP(rw, req)
			} else {
				corsMiddleware(auth.Middleware(printNextFuncName(valor.hndlr), valor.authRole)).ServeHTTP(rw, req)
				/*for c := range chain {
					...
				}*/
			}
		} else {
			rw.Header().Set("Allow", allowedMethodsStr)
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
	contextPath = os.Getenv("CONTEXT_PATH")
	if contextPath == "" {
		contextPath = "/backend" // default
	} else {
		// standarizar
		parts := strings.Split(contextPath, "/")
		contextPath = ""
		for _, p := range parts {
			if p != "" {
				contextPath += "/" + p
			}
		}
		fmt.Println("standarized contextPath: '" + contextPath + "'")
	}
	Mux = http.NewServeMux()
	innerMux := http.NewServeMux()
	//Mux.HandleFunc("OPTIONS /", corsMiddleware(printNextFuncName(corsHandler)))
	innerMux.HandleFunc("/entities", GlobalMiddleware(RouteMap{
		http.MethodGet: Route{
			hndlr:    entity.GetEntities,
			authRole: "entity-read",
		},
		http.MethodPost: Route{
			hndlr:    entity.AddEntity,
			authRole: "entity-create",
		},
	}))
	innerMux.HandleFunc("/entities/{id}", GlobalMiddleware(RouteMap{
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
	Mux.HandleFunc("/", corsMiddleware(printNextFuncName(http.NotFound)))
	Mux.Handle(contextPath+"/", http.StripPrefix(contextPath, innerMux))
}
