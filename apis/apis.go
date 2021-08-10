package apis

// Book struct (Model)
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author struct
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// Author struct
type Healthcheck struct {
	Status   int    `json:"status"`
	Message  string `json: message`
	Datetime string `json: datetime`
}
