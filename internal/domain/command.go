package domain

import "time"

type Command struct {
	ID          int64
	Alias       string
	Description string
	Command     string
	Spotlighted bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c Command) IsValid() bool { return c.Alias != "" }
