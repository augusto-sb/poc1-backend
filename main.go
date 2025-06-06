package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"

	"example.com/router"
)

func handleError(err error) {
	if err != nil {
		//fmt.Println(err)
		panic(err.Error())
	}
}

var version string = "1.0.0"

func main() {
	fmt.Println("Version: ", version)
	var err error
	var listener net.Listener
	var port string = os.Getenv("PORT")
	var srv *http.Server = &http.Server{
		Handler: router.Mux,
	}
	if port == "" {
		port = "8080"
	}
	listener, err = net.Listen("tcp4", "0.0.0.0:"+port)
	handleError(err)
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1) // buffered channel
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		//fmt.Println("closing gracefully")
		err := srv.Shutdown(nil) // context.WithTimeout(context.Background(), time.Second*3)
		handleError(err)
		//no acepta mas conexiones, terminar las activas
		//fmt.Println("requests finished")
		close(sigint)
		close(idleConnsClosed)
	}()
	err = srv.Serve(listener)
	if err != http.ErrServerClosed {
		handleError(err)
	}
	<-idleConnsClosed
}
