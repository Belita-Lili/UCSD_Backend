// entities/notification.go
package entities

type Notification struct {
	Error     string `json:"error"`
	Timestamp string `json:"timestamp"`
	Metadata  string `json:"metadata"`
}
