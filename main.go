package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// EDIT THIS
var (
	BaseURL      = "https://apigolang7.herokuapp.com/"
	ClientID     = "BhWus13WIhI4HI7loycM42"
	ClientSecret = "HJHgKSrqUuhIepCzNkCx7E82RSTN1m47dqPoS1Lf6VA"
)

type numsResponseData struct {
	UserID string  `json:"UserID"`
	Point  float64 `json:"Point"`
	// Jsondata string `json:"Jsondata"`
}

func Authorize(w http.ResponseWriter, req *http.Request) {
	// c, err := auth.New(ClientID, BaseURL+"/callback")
	// if err != nil {
	// 	fmt.Fprintf(w, "error:%v", err)
	// 	return
	// }
	// http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(60 * time.Second)})

	// c.Redirect(w, req)
	fmt.Fprintf(w, "token:authorize")

}

func Callback(w http.ResponseWriter, req *http.Request) {
	// resp, err := auth.ParseRequest(req)
	// if err != nil {
	// 	fmt.Fprintf(w, "error:%v", err)
	// 	return
	// }

	// state, err := req.Cookie("state")
	// if err != nil {
	// 	fmt.Fprintf(w, "error:%v", err)
	// 	return
	// }
	// if resp.State != state.Value {
	// 	fmt.Fprintf(w, "error:%v", err)
	// 	return
	// }

	// c := token.NewClient(BaseURL+"/callback", ClientID, ClientSecret)
	// accessToken, err := c.GetAccessToken(context.Background(), resp.Code)
	// if err != nil {
	// 	fmt.Fprintf(w, "error:%v", err)
	// 	return
	// }
	// adddata(accessToken)

	// // notify(accesstoken)

	// fmt.Fprintf(w, "token:%v", accessToken)
	fmt.Fprintf(w, "token:callback")

}

func main() {

	godotenv.Load()
	port := os.Getenv("PORT")

	router := mux.NewRouter()
	headers := handlers.AllowedOrigins([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedOrigins([]string{"GET", "POST", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/auth", Authorize)
	router.HandleFunc("/callback", Callback)
	// router.HandleFunc("/notify", notificationtoline)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))

}
