package resources

type AddPoll struct {
	Title   string   `json:"title"`
	Options []string `json:"options"`
}

type Poll struct {
	Id string `json:"id"`
	AddPoll
}
