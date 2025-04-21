package main

import(
	"net"
	"net/http"
)

func handleError(err error){
	if(err!=nil){
		panic(err.Error());
	}
}

func main()(){
	//vars
	var err error;
	var listener net.Listener;
	var handlerMux *http.ServeMux;
	//listener
	listener, err = net.Listen("tcp4", "0.0.0.0:8080");
	handleError(err);
	//handler
	handlerMux = http.NewServeMux();
	//routes
	handlerMux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		var err2 error;
		_, err2 = rw.Write([]byte("Hello!\n"));
		handleError(err2);
	});
	//server
	err = http.Serve(listener, handlerMux);
	handleError(err);
}