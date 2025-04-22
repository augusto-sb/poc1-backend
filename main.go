package main

import(
	"net"
	"net/http"
	"os"
)

func handleError(err error){
	if(err!=nil){
		panic(err.Error());
	}
}

var contextPath string;

func init()(){
	var tmp string = os.Getenv("CONTEXT_PATH");
	if(tmp == ""){
		tmp = "/backend";
	}
	contextPath = tmp;
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
	handlerMux.HandleFunc("/", http.NotFound);
	//server
	err = http.Serve(listener, handlerMux);
	handleError(err);
}