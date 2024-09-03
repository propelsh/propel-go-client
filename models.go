package propel

import "time"

type User struct {
	ID        string
	Email     string
	CreatedAt time.Time
}
