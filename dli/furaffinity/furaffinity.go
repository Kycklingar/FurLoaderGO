package fa

import (
	"errors"
	"net/http"
	"net/http/cookiejar"

	"github.com/kycklingar/FurLoaderGO/dli"
)

type furaffinity struct {
	client http.Client

	nextPage map[string]string
}

const (
	faBase       = "https://www.furaffinity.net/"
	faLogin      = faBase + "login/"
	faGallery    = faBase + "gallery/"
	faSubmission = faBase + "view/"
)

func init() {
	var fa furaffinity
	fa.client.Jar, _ = cookiejar.New(nil)
	dli.Logins["furaffinity"] = &fa
	dli.Galleries["furaffinity"] = &fa
}

func httpError(res *http.Response) error {
	if res.StatusCode != 200 {
		return errors.New(res.Status)
	}

	return nil
}
