package model

type Menu struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	ParentID *uint  `json:"parent_id"`
	Route    string `json:"route"`
	Parent   *Menu  `gorm:"foreignKey:ParentID"`
	Children []Menu `gorm:"foreignKey:ParentID"`
}
 