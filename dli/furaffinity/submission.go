package fa

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path"
	"strings"

	"github.com/kycklingar/FurLoaderGO/dli"
	"golang.org/x/net/html"
)

type submission struct {
	id      int
	fileURL string

	scraps bool

	fa   *furaffinity
	user user
}

func (s *submission) SiteName() string {
	return "furaffinity.net"
}

func (s *submission) Folder() string {
	if s.scraps {
		return "scraps"
	}
	return "gallery"
}

func (s *submission) ID() string {
	return fmt.Sprint(s.id)
}

func (s *submission) Filename() string {
	if s.fileURL == "" {
		return ""
	}

	return path.Base(s.fileURL)
}

func (s *submission) FileURL() string {
	return s.fileURL
}

func (s *submission) Download() (io.ReadCloser, error) {
	if s.fa == nil {
		return nil, errors.New("s.fa is nil")
	}

	res, err := s.fa.client.Get(s.fileURL)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if err = httpError(res); err != nil {
		log.Println(err)
		res.Body.Close()
		return nil, err
	}

	return res.Body, nil
}

func (s *submission) GetDetails() ([]dli.Submission, error) {
	if s.fa == nil {
		return nil, errors.New("s.fa is nil")
	}

	res, err := s.fa.client.Get(s.subURL())
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	if err = httpError(res); err != nil {
		log.Println(err)
		return nil, err
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	imgNode := getNodeID(node, "submissionImg")
	if imgNode == nil {
		return nil, errors.New("could not find the image html node")
	}

	for _, a := range imgNode.Attr {
		if a.Key == "data-fullview-src" {
			s.fileURL = "https:" + a.Val
			break
		}
	}

	if s.fileURL == "" {
		return nil, errors.New("file url could not be located")
	}

	pchan := make(chan *html.Node)

	go getNodeClasses(node, "minigallery-title", pchan)

	for {
		gNode := <-pchan

		if gNode == nil {
			break
		}

		titleNode := getNodeFollowingPattern(gNode, []string{"u", "s", "a"})

		s.scraps = strings.TrimSpace(titleNode.FirstChild.Data) == "Scraps"
	}

	go getNodeClasses(node, "classic-submission-title", pchan)

	for {
		gNode := <-pchan

		if gNode == nil {
			break
		}

		var t bool
		for _, a := range gNode.Attr {
			if a.Key == "class" {
				for _, class := range strings.Split(a.Val, " ") {
					if class == "information" {
						t = true
						break
					}
				}
				break
			}
		}

		if !t {
			continue
		}

		var child *html.Node

		for c := gNode.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "a" {
				child = c
			}
		}

		if child == nil || child.Type != html.ElementNode {
			return nil, errors.New("child node is nil")
		}

		for _, a := range child.Attr {
			if a.Key == "href" {
				s.user.id = a.Val[6 : len(a.Val)-1]
			}
		}

		s.user.name = child.FirstChild.Data
	}

	return nil, nil
}

func (s *submission) subURL() string {
	if s.id <= 0 {
		return ""
	}

	return fmt.Sprintf("%s/%d", faSubmission, s.id)
}

func (s *submission) User() dli.User {
	return &s.user
}

type user struct {
	id   string
	name string
}

func (u *user) ID() string {
	return u.id
}

func (u *user) Name() string {
	return u.name
}
