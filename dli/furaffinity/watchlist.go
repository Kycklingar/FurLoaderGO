package fa

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/kycklingar/FurLoaderGO/dli"
	"golang.org/x/net/html"
)

func (fa *furaffinity) Watchlist(username string) ([]dli.User, error) {
	var userlist []dli.User
	var i int
	for {
		i++
		res, err := fa.client.Get(fmt.Sprint(faWatchlist(username), "/", i))
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

		pchan := make(chan *html.Node)

		go getNodeClasses(node, "artist_name", pchan)

		var users []dli.User
		for {
			pnode := <-pchan
			if pnode == nil {
				break
			}

			var user user

			for _, attr := range pnode.Parent.Attr {
				if attr.Key == "href" {
					user.id = attr.Val[6 : len(attr.Val)-1]
					break
				}
			}

			user.name = pnode.FirstChild.Data

			users = append(users, &user)

		}

		if len(users) <= 0 {
			break
		}

		userlist = append(userlist, users...)
	}
	return userlist, nil
}

func (fa *furaffinity) Feed() dli.Feed {
	var feed = &feed{}
	feed.fa = fa
	feed.nextPage = faFeed
	return feed
}

type feed struct {
	fa *furaffinity

	nextPage string
}

func (f *feed) NextPage() ([]dli.Submission, error) {
	res, err := f.fa.client.Get(f.nextPage)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()
	f.nextPage = ""

	if err = httpError(res); err != nil {
		log.Println(err)
		return nil, err
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	getNextPage := func(pch chan *html.Node) string {
		var ret string
		for {
			pnode := <-pch
			if pnode == nil {
				break
			}

			for _, attr := range pnode.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "prev") {
					break
				}

				if attr.Key == "href" {
					ret = faBase + attr.Val
				}
			}
		}
		return ret
	}

	ch := make(chan *html.Node)

	go getNodeClasses(node, "more", ch)
	f.nextPage = getNextPage(ch)

	if f.nextPage == "" {
		go getNodeClasses(node, "more-half", ch)
		f.nextPage = getNextPage(ch)
	}

	subsNode := getNodeID(node, "messagecenter-submissions")
	if subsNode == nil {
		return nil, errors.New("no node named messagecenter-submissions")
	}

	pchan := make(chan *html.Node)
	go getNodeElements(subsNode, "figure", pchan)

	return f.fa.getSubsFromGalleryPage(pchan), nil
}
