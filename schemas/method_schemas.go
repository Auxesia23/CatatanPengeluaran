package schemas

type MethodInput struct {
	Name string `json:"name"`
}

type MethodUpdateInput struct {
	Name string `json:"name"`
}	

type MethodResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
