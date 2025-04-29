// entities/report.go
package entities

type Report struct {
	WebID      string            `json:"web_id"`
	Properties map[string]string `json:"properties"`
	Elements   []string          `json:"elements_to_modify"`
}
