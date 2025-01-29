package schemas

type CategoryInput struct {
	Name string `json:"name"`
}

type CategoryUpdateInput struct {
	Name string `json:"name"`
}
type CategoryResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}