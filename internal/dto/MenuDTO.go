package dto

type MenuResponse struct {
	ID       uint           `json:"id"`
	Name     string         `json:"name"`
	Route    string         `json:"route"`
	Icon     string         `json:"icon"`
	Children []MenuResponse `json:"children,omitempty"`
}
