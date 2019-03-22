package inkbunny

import (
	"fmt"
	"io"
	"log"
)

type ibSub struct {
	id       string
	username string
	fileName string
	url      string
	fileURL  string
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

func (s *ibSub) ID() string {
	return fmt.Sprint(s.id)
}

func (s *ibSub) FileURL() string {
	return s.fileURL
}
