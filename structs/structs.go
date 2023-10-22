package structs

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
