package inkbunny

import "fmt"

type user struct {
	id   int
	name string
}

func (u user) ID() string {
	return fmt.Sprint(u.id)
}

func (u user) Name() string {
	return u.name
}
