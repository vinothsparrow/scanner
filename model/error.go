package model

type Errors struct {
	Errors []*Error `json:"errors"`
}

func (e *Errors) StatusCode() int {
	if len(e.Errors) > 0 {
		return e.Errors[0].Status
	}
	return 500
}

func (e *Errors) HasError() bool {
	return len(e.Errors) > 0
}

type Error struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (e *Error) IsNil() bool {
	return e.Status == 0 && e.Detail == "" && e.Title == ""
}
