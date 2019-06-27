package dli

import "io"

type Login interface {
	Login(username, password string) error
	SetCookies(cookies string) error
	GetCookies() string
}

type Watcher interface {
	Watchlist(string) ([]User, error)
	Feed() Feed
}

type Feed interface {
	NextPage() ([]Submission, error)
}

type Gallery interface {
	Posts(userID string, offset int) ([]Submission, error)
}

type User interface {
	ID() string
	Name() string
}

type Submission interface {
	// What site does this submission come from
	SiteName() string
	// What folder is the submission in? i.e scraps or empty
	Folder() string

	// A unique ID for this submission
	ID() string
	// The filename of the submission
	Filename() string

	// The full url of the submission file
	FileURL() string

	// File download
	Download() (io.ReadCloser, error)
	// Download the submission details
	// this method can spawn additional
	// submissions/files which are nested within
	// the original submission
	GetDetails() ([]Submission, error)

	User() User
}
