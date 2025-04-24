package auth

import(
	"fmt"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var introspectionEndpoint string = "https://example.com/auth/realms/REALM/protocol/openid-connect/token/introspect";
var privateClientId string = "client";
var privateClientSecret string = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx";
var httpClient http.Client = http.Client{};

func init()(){
	//using rbac for simplicity, could use resource/scope based authorization for large project
	introspectionEndpoint = os.Getenv("AUTH_INTROSPECTION_ENDPOINT");
	privateClientId = os.Getenv("AUTH_CLIENT_ID");
	privateClientSecret = os.Getenv("AUTH_CLIENT_SECRET");
}

func Middleware(next http.HandlerFunc, role string)(http.HandlerFunc){
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		//req.Header = map[string][]string
		authHeader, ok := req.Header["Authorization"];
		if(!ok){
			rw.WriteHeader(http.StatusUnauthorized);
			rw.Write([]byte("Authorization header not present"));
			return;
		}
		if(len(authHeader) != 1){
			rw.WriteHeader(http.StatusBadRequest);
			rw.Write([]byte("Authorization header present multiple times"));
			return;
		}
		bearer := authHeader[0];
		if(strings.HasPrefix(strings.ToLower(bearer), "bearer ")){
			bearer = bearer[7:len(bearer)];
		}
		jwt := strings.Split(bearer, ".");
		if(len(jwt) != 3){
			// header, payload, signature
			rw.WriteHeader(http.StatusBadRequest);
			rw.Write([]byte("invalid JWT"));
			return;
		}
		//introspect
		data := url.Values{};
		data.Set("token", bearer);
		data.Set("client_id", privateClientId);
		data.Set("client_secret", privateClientSecret);
		req, err := http.NewRequest(http.MethodPost, introspectionEndpoint, strings.NewReader(data.Encode()));
		if(err != nil){
			rw.WriteHeader(http.StatusInternalServerError);
			return;
		}
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded");
		resp, err := httpClient.Do(req);
		if(err != nil){
			rw.WriteHeader(http.StatusInternalServerError);
			return;
		}
		if(resp.StatusCode != http.StatusOK){
			rw.WriteHeader(http.StatusInternalServerError);
			return;
		}
		bodyBytes, err := io.ReadAll(resp.Body);
		resp.Body.Close();
		if(err != nil){
			rw.WriteHeader(http.StatusInternalServerError);
			return;
		}
		var respJson struct{active bool};
		err = json.Unmarshal(bodyBytes, &respJson);
		if(err != nil){
			rw.WriteHeader(http.StatusInternalServerError);
			return;
		}
		if(!respJson.active){
			/*rw.WriteHeader(http.StatusUnauthorized);
			rw.Write([]byte("inactive token"));
			return;*/
		}
		//StatusForbidden
		payload := jwt[1];
		payloadDecoded, err := base64.RawStdEncoding.DecodeString(payload);
		if(err != nil){
			rw.WriteHeader(http.StatusInternalServerError);
			rw.Write([]byte("error decoding JWT"));
			return;
		}
		var jwtJson map[string]*any;
		err = json.Unmarshal(payloadDecoded, &jwtJson);
		if(err != nil){
			rw.WriteHeader(http.StatusInternalServerError);
			return;
		}
		fmt.Println(jwtJson);
		next.ServeHTTP(rw, req);
	})
}