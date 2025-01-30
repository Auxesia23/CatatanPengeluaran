package schemas

type InputTransaction struct {
	CategoryID  uint      `json:"category_id"`
	MethodID    uint      `json:"method_id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
}
