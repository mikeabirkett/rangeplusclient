package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"rangeplusclient/structs"
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

var (
	cookie      *http.Cookie
	supplier_id int
)

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

	var aReq structs.AuthReq

	err := json.NewDecoder(r.Body).Decode(&aReq)
	if err != nil {
		errMsg := fmt.Sprintf("Could not marshal JSON request : %s", err)
		panic(errMsg)
	}

	apiUrl := "https://supplier.rstore.com/rest/authenticate.api"
	userData := []byte(fmt.Sprintf(`{"user":"%s", "pass":"%s"}`, aReq.UserName, aReq.Password))

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

	var aResp structs.AuthResp

	cookie = response.Cookies()[0]

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

	fmt.Println(cookie)
	fmt.Println(&aResp)

	supplier_id = aResp.Supplier_id

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Authentication Endpoint"))
}

func orderFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}

	var Orders structs.Orders

	apiUrl := fmt.Sprintf("https://supplier.rstore.com/rest/order_feed.api?supplier_id=%d", supplier_id)

	request, err := http.NewRequest("GET", apiUrl, bytes.NewBuffer([]byte(`{"mode":"Live"}`)))
	if err != nil {
		errMsg := fmt.Sprintf("Could not create new API request : %s", err)
		panic(errMsg)
	}

	request.Header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		errMsg := fmt.Sprintf("Could not send API request : %s", err)
		panic(errMsg)
	}

	if response.StatusCode != http.StatusOK {
		fmt.Println(response.Status)
		fmt.Println(response.Body)

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Order Feed Endpoint"))

		return
	}

	err = json.NewDecoder(response.Body).Decode(&Orders)
	if err != nil {
		errMsg := fmt.Sprintf("Could not decode API response : %s", err)
		panic(errMsg)
	}

	for _, order := range Orders.Array {
		fmt.Printf("%+v\n\n", order)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Order Feed Endpoint"))
}
