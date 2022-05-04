package model

type Stat struct {
	Notifications stat_notifications
	Messages stat_messages
}

type stat_notifications struct {
	Finished string
	Sent string
	Will_be_sent string
}

type stat_messages struct {
	Delivered string
	Failure string
}