package model

type Notification struct {
	Id          string `json:"id"`
	Start_data  string `json:"start_data"`
	Text        string `json:"text"`
	Filter      Filter `json:"filter"`
	Finish_data string `json:"finish_data"`
}

type Filter struct {
	Num_kod	string `json:"num_kod"`
	Teg		string `json:"teg"`
}