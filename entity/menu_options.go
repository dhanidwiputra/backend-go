package entity

type MenuOption struct {
	Title           string           `json:"title" binding:"required"`
	Type            string           `json:"type" binding:"required"`
	Max             int              `json:"max,string" binding:"required"`
	MenuOptionLists []MenuOptionList `json:"menu_option_lists" binding:"required"`
}
