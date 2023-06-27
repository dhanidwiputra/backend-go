package entity

type MenuOptionList struct {
	Name        string `json:"name" binding:"required"`
	Price       int    `json:"price,string" binding:"required"`
	Description string `json:"description" binding:"required"`
	Checked     bool   `json:"checked"`
}
