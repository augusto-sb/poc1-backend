package main

import(
	"encoding/json"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
)

func handleError(err error){
	if(err!=nil){
		panic(err.Error());
	}
}

type Entity struct{
	Id uint64
	Name string
	Description string
}

var contextPath string;
var dataBase []Entity = []Entity{
	Entity{
		Id: 0,
		Name: "First",
		Description: "the first element.",
	},
	Entity{
		Id: 1,
		Name: "Second",
		Description: "the second element.",
	},
	Entity{
		Id: 2,
		Name: "Third",
		Description: "the third element.",
	},
};
var mu sync.Mutex = sync.Mutex{};

func init()(){
	var tmp string = os.Getenv("CONTEXT_PATH");
	if(tmp == ""){
		tmp = "/backend";
	}
	contextPath = tmp;
}

func getEntities(rw http.ResponseWriter, req *http.Request)(){
	mu.Lock();
	jsonByteArr, err := json.Marshal(dataBase)
	mu.Unlock();
	if(err != nil){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	_, err = rw.Write(jsonByteArr);
	if(err != nil){
		rw.WriteHeader(http.StatusInternalServerError);
	}
}

func getEntity(rw http.ResponseWriter, req *http.Request)(){
	var idParam string = req.PathValue("id");
	var id uint64;
	var err error;
	var jsonByteArr []byte = []byte{};

	id, err = strconv.ParseUint(idParam, 10, 64);
	if(err != nil){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	mu.Lock();
	for _, v := range dataBase {
		if(v.Id == id){
			jsonByteArr, err = json.Marshal(v);
			break;
		}
	}
	mu.Unlock();
	if(err != nil){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	if(len(jsonByteArr) == 0){
		rw.WriteHeader(http.StatusNotFound);
		return;
	}
	_, err = rw.Write(jsonByteArr);
	if(err != nil){
		rw.WriteHeader(http.StatusInternalServerError);
	}
}

func removeEntity(rw http.ResponseWriter, req *http.Request)(){
	var idParam string = req.PathValue("id");
	var id uint64;
	var err error;
	var isDeleted bool = false;

	id, err = strconv.ParseUint(idParam, 10, 64);
	if(err != nil){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	mu.Lock();
	for i, v := range dataBase {
		if(v.Id == id){
			dataBase[i] = dataBase[len(dataBase)-1];
			dataBase = dataBase[:len(dataBase)-1]
			isDeleted = true;
			break;
		}
	}
	mu.Unlock();
	if(isDeleted){
		rw.WriteHeader(http.StatusNoContent);
	}else{
		rw.WriteHeader(http.StatusNotFound);
	}
}

func addEntity(rw http.ResponseWriter, req *http.Request)(){
	var e Entity;
	var tmp map[string]any;
	var id uint64 = 0;

	err := json.NewDecoder(req.Body).Decode(&tmp); //decode to map/any instead of struct to validate not extra keys...
	if(err != nil){
		rw.WriteHeader(http.StatusBadRequest);
		return;
	}
	//validations
	if(len(tmp) != 2){
		rw.WriteHeader(http.StatusBadRequest);
		rw.Write([]byte("expected exactly 2 keys"));
		return;
	}
	valName, okName := tmp["Name"];
	valDescription, okDescription := tmp["Description"];
	if(!okName || !okDescription){
		rw.WriteHeader(http.StatusBadRequest);
		rw.Write([]byte("Name or Description missing!"));
		return;
	}
	switch valName.(type) {
	default:
		rw.WriteHeader(http.StatusBadRequest);
		rw.Write([]byte("Name should be string!"));
		return;
	case string:
		e.Name = valName.(string);
	}
	switch valDescription.(type) {
	default:
		rw.WriteHeader(http.StatusBadRequest);
		rw.Write([]byte("Description should be string!"));
		return;
	case string:
		e.Description = valDescription.(string);
	}
	//generate id and append
	mu.Lock();
	for _, v := range dataBase {
		if(v.Id > id){
			id = v.Id;
		}
	}
	e.Id = id+1;
	dataBase = append(dataBase, e);
	mu.Unlock();
	rw.WriteHeader(http.StatusCreated);
}

func updateEntity(rw http.ResponseWriter, req *http.Request)(){
	var idParam string = req.PathValue("id");
	var id uint64;
	var err error;
	var isUpdated bool = false;
	id, err = strconv.ParseUint(idParam, 10, 64);
	if(err != nil){
		rw.WriteHeader(http.StatusInternalServerError);
		return;
	}
	var tmp map[string]any;

	err = json.NewDecoder(req.Body).Decode(&tmp); //decode to map/any instead of struct to validate not extra keys...
	if(err != nil){
		rw.WriteHeader(http.StatusBadRequest);
		return;
	}
	//validations
	if(len(tmp) != 2){
		rw.WriteHeader(http.StatusBadRequest);
		rw.Write([]byte("expected exactly 2 keys"));
		return;
	}
	valName, okName := tmp["Name"];
	valDescription, okDescription := tmp["Description"];
	if(!okName || !okDescription){
		rw.WriteHeader(http.StatusBadRequest);
		rw.Write([]byte("Name or Description missing!"));
		return;
	}
	switch valName.(type) {
	default:
		rw.WriteHeader(http.StatusBadRequest);
		rw.Write([]byte("Name should be string!"));
		return;
	case string:
		//is ok!
	}
	switch valDescription.(type) {
	default:
		rw.WriteHeader(http.StatusBadRequest);
		rw.Write([]byte("Description should be string!"));
		return;
	case string:
		//is ok!
	}
	//find and update
	mu.Lock();
	for i, v := range dataBase {
		if(v.Id == id){
			dataBase[i].Name = valName.(string);
			dataBase[i].Description = valDescription.(string);
			isUpdated = true;
			break;
		}
	}
	mu.Unlock();
	if(isUpdated){
		rw.WriteHeader(http.StatusNoContent);
	}else{
		rw.WriteHeader(http.StatusNotFound);
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
	handlerMux.HandleFunc("/", http.NotFound);
	handlerMux.HandleFunc("GET "+contextPath+"/entities", getEntities);
	handlerMux.HandleFunc("GET "+contextPath+"/entities/{id}", getEntity);
	handlerMux.HandleFunc("POST "+contextPath+"/entities", addEntity);
	handlerMux.HandleFunc("DELETE "+contextPath+"/entities/{id}", removeEntity);
	handlerMux.HandleFunc("PUT "+contextPath+"/entities/{id}", updateEntity);
	//server
	err = http.Serve(listener, handlerMux);
	handleError(err);
}