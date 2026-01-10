package http

type BoardDTO struct {
	Id          string `json:"id"`
	OwnerId     string `json:"owner_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type ColumnDTO struct {
	Id       string `json:"id"`
	BoardId  string `json:"board_id"`
	Position int32  `json:"position"`
}

type TaskDTO struct {
	Id          string `json:"id"`
	ColumnId    string `json:"column_id"`
	Position    int32  `json:"position"`
	Title       string `json:"title"`
	Description string `json:"description"`
	AssigneeId  string `json:"assignee_id"`
}
