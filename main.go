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
	http.HandleFunc("/publish_product", publishProduct)

	err := http.ListenAndServe(":5000", nil)
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

	var aReq AuthReq

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

	var aResp AuthResp

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

	var Orders Orders

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

func publishProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusBadRequest)
	}

	// var Orders structs.Orders

	// apiUrl := fmt.Sprintf("https://supplier.rstore.com/rest/product_feed.api?supplier_id=%d", supplier_id)

	// request, err := http.NewRequest("GET", apiUrl, bytes.NewBuffer([]byte(`{"mode":"Live"}`)))
	// if err != nil {
	// 	errMsg := fmt.Sprintf("Could not create new API request : %s", err)
	// 	panic(errMsg)
	// }

	// request.Header.Add("Cookie", fmt.Sprintf("%s=%s", cookie.Name, cookie.Value))

	// request.Header.Set("Content-Type", "application/json; charset=utf-8")
	// client := &http.Client{}
	// response, err := client.Do(request)
	// if err != nil {
	// 	errMsg := fmt.Sprintf("Could not send API request : %s", err)
	// 	panic(errMsg)
	// }

	// if response.StatusCode != http.StatusOK {
	// 	fmt.Println(response.Status)
	// 	fmt.Println(response.Body)

	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("Publish Product Endpoint"))

	// 	return
	// }

	// err = json.NewDecoder(response.Body).Decode(&Orders)
	// if err != nil {
	// 	errMsg := fmt.Sprintf("Could not decode API response : %s", err)
	// 	panic(errMsg)
	// }

	// for _, order := range Orders.Array {
	// 	fmt.Printf("%+v\n\n", order)
	// }

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Publish Product Endpoint"))
}

type AuthReq struct {
	UserName string
	Password string
}

type AuthResp struct {
	Mode        string `json:"mode"`
	Supplier_id int    `json:"supplier_id"`
}

type Order struct {
	OrderDisp            string `json:"order_disp"`
	CustomerName         string `json:"customer_name"`
	PostCode             string `json:"postcode"`
	BuildingNameNumber   string `json:"building_name_number"`
	Organisation         string `json:"organisation"`
	Street               string `json:"street"`
	City                 string `json:"city"`
	County               string `json:"county"`
	Country              string `json:"country"`
	CustomerPhone        string `json:"customer_telephone"`
	CustomerEmailAddress string `json:"customer_email_address"`
	ProductCode          string `json:"product_code"`
	Title                string `json:"title"`
	Quantity             int    `json:"qty"`
	Status               string `json:"status"`
	SKU                  int    `json:"sku"`
	Price                int    `json:"price"`
	OrderPlacedDate      string `json:"order_placed_date"`
	DespatchDate         string `json:"despatch_date"`
	CourierName          string `json:"courier_name"`
	DeliveryService      string `json:"delivery_service"`
	TrackingReference    string `json:"tracking_reference"`
	Notes                string `json:"notes"`
}

type Orders struct {
	Array []Order `json:"order_Arr"`
}

type PublishResponseItem struct {
	Label   string `json:"label"`
	SKUList string `json:"sku_list"`
}

type PublishResponse struct {
	Array []PublishResponseItem `json:"result"`
}
