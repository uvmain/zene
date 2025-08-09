package types

import "encoding/xml"

type Chat struct {
	UserName  string `xml:"username,attr" json:"user_name"`
	Message   string `xml:"message,attr" json:"message"`
	Timestamp int64  `xml:"timestamp,attr" json:"timestamp"`
}

type ChatMessages struct {
	ChatMessage []Chat `xml:"chatMessage" json:"chatMessage"`
}

type SubsonicChatMessages struct {
	XMLName       xml.Name       `xml:"subsonic-response" json:"-"`
	Xmlns         string         `xml:"xmlns,attr" json:"-"`
	Status        string         `xml:"status,attr" json:"status"`
	Version       string         `xml:"version,attr" json:"version"`
	Type          string         `xml:"type,attr" json:"type"`
	ServerVersion string         `xml:"serverVersion,attr" json:"serverVersion"`
	OpenSubsonic  bool           `xml:"openSubsonic,attr" json:"openSubsonic"`
	Error         *SubsonicError `xml:"error,omitempty" json:"error,omitempty"`
	ChatMessages  *ChatMessages  `xml:"chatMessages" json:"chatMessages"`
}

type SubsonicChatMessagesResponse struct {
	SubsonicResponse SubsonicChatMessages `json:"subsonic-response"`
}
