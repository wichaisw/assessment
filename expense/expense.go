package expense

type Expense struct {
	Id    int      `json:"id`
	Title string   `json:"title"`
	Note  string   `json:"note"`
	Tags  []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}
