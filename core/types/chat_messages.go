package types

type Chat struct {
	UserName  string `xml:"username,attr" json:"user_name"`
	Message   string `xml:"message,attr" json:"message"`
	Timestamp int    `xml:"timestamp,attr" json:"timestamp"`
}

type ChatMessages struct {
	ChatMessage []Chat `xml:"chatMessage" json:"chatMessage"`
}
