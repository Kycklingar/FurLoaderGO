package inkbunny

import (
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/kycklingar/FurLoaderGO/dli"
)

type ibSub struct {
	id       string
	fileName string
	url      string
	fileURL  string
	folder   string

	user user
}

func (s *ibSub) SiteName() string {
	return "inkbunny.net"
}

func (s *ibSub) Folder() string {
	return s.folder
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

func (s *ibSub) GetDetails() error {
	return nil
}

func (s *ibSub) User() dli.User {
	return s.user
}

func (ib *InkBunny) fromJson(j ibJsonSub) ([]ibSub, error) {
	var subs []ibSub

	v, err := strconv.Atoi(j.PageCount)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if v > 1 {
		s, err := ib.getFileUrls(j.SubID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		subs = append(subs, s...)
	} else {
		var s ibSub
		s.user.name = j.Username
		s.fileName = j.FileName
		s.fileURL = j.FileURL
		s.folder = "gallery"
		if j.Scraps == "t" {
			s.folder = "scraps"
		}
		subs = append(subs, s)
	}
	return subs, nil
}

func (ib *InkBunny) getFileUrls(id string) ([]ibSub, error) {
	v := ib.sidURLValues()
	v.Set("submission_ids", fmt.Sprint(id))

	res, err := client.PostForm(apiSubmissions, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	if err = httpError(res); err != nil {
		log.Println(err)
		return nil, err
	}

	var a ibJsonSearch
	if err = a.decode(res.Body); err != nil {
		log.Println(err)
		return nil, err
	}

	var subs []ibSub

	for _, file := range a.Submissions[0].Files {
		var s ibSub
		s.user.name = a.Submissions[0].Username
		s.id = fmt.Sprintf("s%sf%s", id, file.FileID)
		s.fileURL = file.FileURL
		s.fileName = file.FileName
		s.folder = "gallery"
		if a.Submissions[0].Scraps == "t" {
			s.folder = "scraps"
		}
		subs = append(subs, s)
	}

	return subs, nil
}
