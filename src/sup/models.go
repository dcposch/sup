package sup

import (
    "time"
)

// We ask you what you're up to every hour or so. Each answer is a status.
type Status struct {
    IP               string
	User             string
	Tags             string
	Description      string
	CreateTime       time.Time
}
