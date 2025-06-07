package model

// Todo はTodoモデルを表す
type Todo struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

// TodoQuery はTodoの検索クエリを表す
type TodoQuery struct {
	Status string `form:"status"`
}
