package handlers

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"error"`
	Help    string `json:"help,omitempty"`
}
