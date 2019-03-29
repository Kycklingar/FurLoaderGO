package fa

import (
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/kycklingar/FurLoaderGO/dli"
)

func NewFurAffinity() *furaffinity {
	var f furaffinity
	f.nextPage = make(map[string]*page)
	f.client.Jar, _ = cookiejar.New(nil)
	return &f
}

type furaffinity struct {
	client http.Client

	nextPage map[string]*page
}

type page struct {
	page     int
	location string
}

const (
	faBase       = "https://www.furaffinity.net/"
	faLogin      = faBase + "login/"
	faGallery    = faBase + "gallery/"
	faScraps     = faBase + "scraps/"
	faSubmission = faBase + "view/"
)

func init() {
	var fa = NewFurAffinity()
	dli.Logins["furaffinity"] = fa
	dli.Galleries["furaffinity"] = fa
}

func httpError(res *http.Response) error {
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}
