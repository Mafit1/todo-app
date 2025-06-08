package dto

type UpdateTodoRequest struct {
	Title     *string `json:"title" validate:"omitempty,min=1"`
	Completed *bool   `json:"completed"`
}
