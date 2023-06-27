package dto

type Query struct {
	Search   string `form:"s"`
	SortBy   string `form:"sortBy,default=order_date"`
	Sort     string `form:"sort,default=desc"`
	Limit    int    `form:"limit,default=10"`
	Page     int    `form:"page,default=1"`
	Category string `form:"cat,default="`
	Days     string `form:"days,default=0"`
}
