package main

import(
	"net"
	"net/http"

	"example.com/router"
)

func handleError(err error){
	if(err!=nil){
		panic(err.Error());
	}
}

func main()(){
	var err error;
	var listener net.Listener;
	listener, err = net.Listen("tcp4", "0.0.0.0:8081");
	handleError(err);
	err = http.Serve(listener, router.Mux);
	handleError(err);
}