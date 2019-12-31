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
	fileID   string
	fileName string
	url      string
	fileURL  string
	folder   string

	lastFileUpdate string

	pageCount int

	user user

	ib *InkBunny
}

func (s *ibSub) SiteName() string {
	return "inkbunny.net"
}

func (s *ibSub) Folder() string {
	return s.folder
}

func (s *ibSub) ID() string {
	return s.id + ":" + s.lastFileUpdate
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

func (s *ibSub) GetDetails() ([]dli.Submission, error) {
	if s.pageCount > 1 {
		subs, err := s.ib.getFileUrls(s.id)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		jsubs := make([]dli.Submission, len(subs))
		for i, _ := range subs {
			jsubs[i] = &subs[i]
		}

		return jsubs, nil
	}
	return nil, nil
}

func (s *ibSub) User() dli.User {
	return s.user
}

func (ib *InkBunny) fromJson(j ibJsonSub) (ibSub, error) {
	var s ibSub

	s.ib = ib

	s.id = j.SubID
	s.user.name = j.Username
	s.fileName = j.FileName
	s.fileURL = j.FileURL
	s.lastFileUpdate = j.LastFileUpdate
	s.folder = "gallery"
	if j.Scraps == "t" {
		s.folder = "scraps"
	}

	var err error
	s.pageCount, err = strconv.Atoi(j.PageCount)
	if err != nil {
		log.Println(err)
		return s, err
	}

	return s, nil
}

func (ib *InkBunny) getFileUrls(id string) ([]ibSub, error) {
	if id == "" {
		return nil, fmt.Errorf("no id specified in search")
	}
	v := ib.sidURLValues()
	v.Set("submission_ids", id)

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

	var a ibJsonSubmissions
	if err = a.decode(res.Body); err != nil {
		log.Println(err)
		return nil, err
	}

	var subs []ibSub

	if len(a.Submissions) <= 0 {
		erstr := "no submissions returned from submissions"
		log.Println(erstr)
		return nil, fmt.Errorf(erstr)
	}

	for _, file := range a.Submissions[0].Files {
		var s ibSub
		s.ib = ib
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
