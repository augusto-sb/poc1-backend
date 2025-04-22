package router

import(
	"net/http"
	"os"

	"example.com/entity"
)

var Mux *http.ServeMux;

func init()(){
	var contextPath string;
	contextPath = os.Getenv("CONTEXT_PATH");
	if(contextPath == ""){
		contextPath = "/backend";
	}
	Mux = http.NewServeMux();
	Mux.HandleFunc("/", http.NotFound);
	Mux.HandleFunc("GET "+contextPath+"/entities", entity.GetEntities);
	Mux.HandleFunc("GET "+contextPath+"/entities/{id}", entity.GetEntity);
	Mux.HandleFunc("POST "+contextPath+"/entities", entity.AddEntity);
	Mux.HandleFunc("DELETE "+contextPath+"/entities/{id}", entity.RemoveEntity);
	Mux.HandleFunc("PUT "+contextPath+"/entities/{id}", entity.UpdateEntity);
}