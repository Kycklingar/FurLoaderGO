package inkbunny

import (
	"fmt"
	"io"
	"log"

	"github.com/kycklingar/FurLoaderGO/dli"
)

type ibSub struct {
	id       string
	fileName string
	url      string
	fileURL  string

	user user
}

func (s *ibSub) SiteName() string {
	return "inkbunny.net"
}

func (s *ibSub) Folder() string {
	//TODO:
	return ""
}

func (s *ibSub) ID() string {
	return fmt.Sprint(s.id)
}

func (s *ibSub) Filename() string {
	return s.fileName
}

func (s *ibSub) FileURL() string {
	return s.fileURL
}

func (s *ibSub) Download() (io.ReadCloser, error) {
	res, err := client.Get(s.fileURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = httpError(res); err != nil {
		res.Body.Close()
		log.Println(err)
		return nil, err
	}

	return res.Body, nil
}

func (s *ibSub) User() dli.User {
	return s.user
}
