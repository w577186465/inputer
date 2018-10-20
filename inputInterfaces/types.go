package inputInterfaces

type Article struct {
	Id      int
	Content string
	Title   string
	Keyword string
}

type Message struct {
	Status  string
	Code    int
	Message string
}
