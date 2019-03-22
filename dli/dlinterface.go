package dli

type Login interface {
	Login(username, password string) error
	SetCookies(cookies string) error
	GetCookies() string
}

type Watcher interface {
	Watchlist() error
	Feed(int) ([]Submission, error)
}

type Gallery interface {
	Posts(offset int) error
}

type Submission interface {
	ID() string
	FileURL() string
}
