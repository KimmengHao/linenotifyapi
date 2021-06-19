package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/utahta/go-linenotify"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/utahta/go-linenotify/auth"
	"github.com/utahta/go-linenotify/token"
	"google.golang.org/api/iterator"
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
	c, err := auth.New(ClientID, BaseURL+"/callback")
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "state", Value: c.State, Expires: time.Now().Add(60 * time.Second)})

	c.Redirect(w, req)
}

func Callback(w http.ResponseWriter, req *http.Request) {
	resp, err := auth.ParseRequest(req)
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	state, err := req.Cookie("state")
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	if resp.State != state.Value {
		fmt.Fprintf(w, "error:%v", err)
		return
	}

	c := token.NewClient(BaseURL+"/callback", ClientID, ClientSecret)
	accessToken, err := c.GetAccessToken(context.Background(), resp.Code)
	if err != nil {
		fmt.Fprintf(w, "error:%v", err)
		return
	}
	adddata(accessToken)

	// notify(accesstoken)

	fmt.Fprintf(w, "token:%v", accessToken)
}

func adddata(accessToken string) {
	// Use the application default credentials
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: "imake-flutter-firebase"}
	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln(err)
	}
	// getdocdata = "kim"

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	iter := client.Collection("tm_members_uat").OrderBy("issue_date", firestore.Desc).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Ref.ID, "\n")
		fmt.Println(doc.Data(), "\n")
		// getdocdata = doc.Data()
		// break
		_, err = client.Collection("tm_members_uat").Doc(doc.Ref.ID).Set(ctx, map[string]interface{}{
			"lineuserid": accessToken,
		}, firestore.MergeAll)

		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}
		break
	}

	// fmt.Println(getdocdata)
	// doc := make(map[string]interface{})
	// doc["name"] = "Hello Tokyo!"
	// doc["country"] = "Japan"

	// _, _, err = client.Collection("tm_members_uat").Add(ctx, doc)
	// if err != nil {
	// 	// Handle any errors in an appropriate way, such as returning them.
	// 	log.Printf("An error has occurred: %s", err)
	// }

}

func notificationtoline(response http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)

	var numsData numsResponseData

	decoder.Decode(&numsData)
	fmt.Println("numsData")
	fmt.Println(numsData)
	fmt.Println(numsData.UserID)

	token := numsData.UserID // EDIT THIS
	msgtext := fmt.Sprintf("%s%.2f", "Your current point is ", numsData.Point)

	c := linenotify.NewClient()
	c.Notify(context.Background(), token, msgtext, "", "", nil)
}
func handleRequests() {

	godotenv.Load()
	port := os.Getenv("PORT")

	router := mux.NewRouter()
	headers := handlers.AllowedOrigins([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedOrigins([]string{"GET", "POST", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	router.HandleFunc("/auth", Authorize)
	router.HandleFunc("/callback", Callback)
	router.HandleFunc("/notify", notificationtoline)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router)))

}

func main() {
	// adddata()

	handleRequests()
}
