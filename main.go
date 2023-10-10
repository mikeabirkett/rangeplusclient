package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/authenticate", authenticate)
	http.HandleFunc("/order_feed", orderFeed)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

var cookie string

func getRoot(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello World"))
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}

	type authReq struct {
		UserName string
		Password string
	}

	type authResp struct {
		mode        string
		supplier_id int32
	}

	var aReq authReq

	err := json.NewDecoder(r.Body).Decode(&aReq)
	if err != nil {
		errMsg := fmt.Sprintf("Could not marshal JSON request : %s", err)
		panic(errMsg)
	}

	apiUrl := "https://uatsupplier.rstore.com/rest/authenticate.api"
	userData := []byte(`{"user":"` + aReq.UserName + `", "pass":"` + aReq.Password + `"}`)

	request, err := http.NewRequest("GET", apiUrl, bytes.NewBuffer(userData))
	if err != nil {
		errMsg := fmt.Sprintf("Could not create new API request : %s", err)
		panic(errMsg)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errMsg := fmt.Sprintf("Could not send API request : %s", err)
		panic(errMsg)
	}

	var aResp authResp

	c := response.Cookies()

	if response.StatusCode != http.StatusOK {
		fmt.Println(response.Status)
		fmt.Println(response.Body)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Authentication Endpoint"))

		return
	}

	err = json.NewDecoder(response.Body).Decode(&aResp)
	if err != nil {
		errMsg := fmt.Sprintf("Could not decode API response : %s", err)
		panic(errMsg)
	}

	fmt.Println(c)
	fmt.Println(&aResp)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authentication Endpoint"))

}

func orderFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order Feed Endpoint"))
}
