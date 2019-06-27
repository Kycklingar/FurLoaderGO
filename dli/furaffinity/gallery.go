package fa

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/kycklingar/FurLoaderGO/dli"
	"golang.org/x/net/html"
)

func (fa *furaffinity) Posts(userID string, offset int) ([]dli.Submission, error) {
	nextPage, ok := fa.nextPage[userID]
	if !ok {
		nextPage = &page{
			page:     offset + 1,
			location: faGallery + userID + "/",
		}

		fa.nextPage[userID] = nextPage
	}

	res, err := fa.client.Get(nextPage.location + fmt.Sprint(nextPage.page))
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

	subsNode := getNodeID(node, "gallery-gallery")
	if subsNode == nil {
		return nil, errors.New("no node named 'gallery-gallery'!")
	}

	pchan := make(chan *html.Node)
	go getNodeElements(subsNode, "figure", pchan)

	subs := fa.getSubsFromGalleryPage(pchan)

	if len(subs) <= 0 {
		if strings.Contains(fa.nextPage[userID].location, "gallery") {
			fa.nextPage[userID] = &page{
				page:     nextPage.page - offset,
				location: faScraps + userID + "/",
			}
			return fa.Posts(userID, offset)
		}
	} else {
		nextPage.page += 1
	}

	return subs, nil
}

func (fa *furaffinity) getSubsFromGalleryPage(pchan chan *html.Node) []dli.Submission {
	var subs []dli.Submission
	for {
		subNode := <-pchan
		if subNode == nil {
			break
		}

		var s submission
		s.fa = fa

		for _, a := range subNode.Attr {
			if a.Key == "id" {
				// Clean '/view/12345/'
				var err error
				s.id, err = strconv.Atoi(a.Val[4:])
				if err != nil {
					log.Println("could not convert href to id:", a.Val)
				}
			}
		}

		if s.id <= 0 {
			log.Println("s.id == 0")
			continue
		}

		subs = append(subs, &s)

	}

	return subs
}

func getNodeElements(node *html.Node, element string, ch chan *html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == element {
			ch <- n
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	ch <- nil
}

func getNodeClasses(node *html.Node, class string, ch chan *html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "class" {
					for _, val := range strings.Split(a.Val, " ") {
						if val == class {
							ch <- n
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
	ch <- nil
}

func getNodeID(n *html.Node, id string) *html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "id" && a.Val == id {
				return n
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		node := getNodeID(c, id)
		if node != nil {
			return node
		}
	}
	return nil
}

func getNodeFollowingPattern(node *html.Node, tags []string) *html.Node {
	var f func(*html.Node, int) *html.Node
	f = func(n *html.Node, seq int) *html.Node {
		if n.Type == html.ElementNode {
			if n.Data == tags[seq] {
				seq++
				if seq >= len(tags) {
					return n
				}
			} else {
				seq = 0
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			r := f(c, seq)
			if r != nil {
				return r
			}
		}
		return nil
	}

	return f(node, 0)
}
