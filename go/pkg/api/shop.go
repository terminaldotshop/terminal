package api

type WidgetResponse struct {
	Count     int      `json:"count"`
	PageCount int      `json:"page_count"`
	Page      int      `json:"page"`
	Results   []Widget `json:"results"`
}

type Widget struct {
	Art  string
	Name string `json:"name"`
	// Sku             string `json:"sku"`
	Active bool `json:"active"`
	// Images          any    `json:"images"`
	// PurchaseOptions struct {
	// 	Standard struct {
	// 		Active    bool  `json:"active"`
	// 		Price     int   `json:"price"`
	// 		Sale      bool  `json:"sale"`
	// 		SalePrice any   `json:"sale_price"`
	// 		Prices    []any `json:"prices"`
	// 	} `json:"standard"`
	// } `json:"purchase_options"`
	// Variable        bool      `json:"variable"`
	Description string `json:"description"`
	// Tags            []any     `json:"tags"`
	// MetaTitle       any       `json:"meta_title"`
	// MetaDescription any       `json:"meta_description"`
	// Slug            string    `json:"slug"`
	// Attributes      struct{}  `json:"attributes"`
	// Delivery        string    `json:"delivery"`
	// Bundle          any       `json:"bundle"`
	Price float64 `json:"price"`
	// StockTracking   bool      `json:"stock_tracking"`
	// Options         []any     `json:"options"`
	// Currency        string    `json:"currency"`
	// Sale            bool      `json:"sale"`
	SalePrice any `json:"sale_price"`
	// Prices          []any     `json:"prices"`
	// Type            string    `json:"type"`
	// DateCreated     time.Time `json:"date_created"`
	// StockStatus     any       `json:"stock_status"`
	// DateUpdated     time.Time `json:"date_updated"`
	ID string `json:"id"`
}

func GetWidgets() []Widget {
	return []Widget{
		{
			Art: `
                  ██████
                ██      ██
      ██████    ██      ██    ██████
    ██      ████          ████      ██
    ██                              ██
    ██                              ██
      ██                          ██
      ██                          ██
  ████            ██████            ████
██              ██      ██              ██
██              ██      ██              ██
██              ██      ██              ██
  ████            ██████            ████
      ██                          ██
      ██                          ██
    ██                              ██
    ██                              ██
    ██      ████          ████      ██
      ██████    ██      ██    ██████
                ██      ██
                  ██████`,
			Name:        "Foo Widget",
			Price:       5.0,
			Active:      true,
			Description: "Test Description",
		},
	}
}
