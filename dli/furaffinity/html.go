package fa

import (
	"strings"

	"golang.org/x/net/html"
)

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

func getNodeClassesFull(node *html.Node, class string, ch chan *html.Node) {
	defer close(ch)

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Key == "class" {
					if a.Val == class {
						ch <- n
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(node)
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
