package database

import (
	"fmt"
	"time"
)

type Author struct {
	name  string
	email string
	time  time.Time
}

func (a *Author) New(name string, email string, time time.Time) error {
	a.name = name
	a.email = email
	a.time = time
	return nil
}

func (a *Author) ToString() string {
	return fmt.Sprintf("%s <%s> %s", a.name, a.email, a.time.Format(time.RFC3339))
}
