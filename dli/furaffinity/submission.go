package fa

import (
	"errors"
	"fmt"
	"io"
	"log"
	"path"

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

	//imgNode := getNodeID(node, "submissionImg")
	//if imgNode == nil {
	//	return nil, errors.New("could not find the image html node")
	//}

	//for _, a := range imgNode.Attr {
	//	if a.Key == "data-fullview-src" {
	//		s.fileURL = "https:" + a.Val
	//		break
	//	}
	//}

	// Get file link
	pchan := make(chan *html.Node)
	go getNodeClasses(node, "download", pchan)

	for {
		pnode := <-pchan
		if pnode == nil {
			break
		}

		if pnode.FirstChild == nil {
			break
		}

		for _, attr := range pnode.FirstChild.Attr {
			if attr.Key == "href" {
				s.fileURL = "https:" + attr.Val
			}
		}
	}

	if s.fileURL == "" {
		return nil, errors.New("file url could not be located")
	}

	// It's impossible to tell if a submission is a scrap or not in the new FA design
	s.scraps = false
	//pchan = make(chan *html.Node)
	//go getNodeClasses(node, "minigallery-title", pchan)

	//for {
	//	gNode := <-pchan

	//	if gNode == nil {
	//		break
	//	}

	//	titleNode := getNodeFollowingPattern(gNode, []string{"u", "s", "a"})

	//	s.scraps = strings.TrimSpace(titleNode.FirstChild.Data) == "Scraps"
	//}

	// Get username and user id
	pchan = make(chan *html.Node)
	go getNodeClasses(node, "submission-id-sub-container", pchan)

	for {
		gNode := <-pchan

		if gNode == nil {
			break
		}

		var child *html.Node

		for c := gNode.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode && c.Data == "a" {
				child = c
				break
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

		s.user.name = child.FirstChild.FirstChild.Data
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
