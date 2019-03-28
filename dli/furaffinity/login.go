package fa

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func (fa *furaffinity) Login(username string, password string) error {
	res, err := fa.client.Get(faLogin)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	if err = httpError(res); err != nil {
		log.Println(err)
		return err
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		res.Body.Close()
		log.Println(err)
		return err
	}

	var captchaNode *html.Node
	var f func(*html.Node) bool
	f = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "captcha_img" {
					captchaNode = n
					return true
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if f(c) {
				return true
			}
		}
		return false
	}

	if !f(node) {
		return errors.New("Couldn't find the captcha image")
	}

	var captchaURL string

	for _, a := range captchaNode.Attr {
		if a.Key == "src" {
			captchaURL = a.Val
		}
	}

	captchaRes, err := fa.client.Get(faBase + captchaURL)
	if err != nil {
		log.Println(err)
		return err
	}
	defer captchaRes.Body.Close()

	if err = httpError(captchaRes); err != nil {
		log.Println(err)
		return err
	}

	file, err := os.OpenFile("./fa-captcha.jpg", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	//defer os.Remove(file.Name())
	defer file.Close()

	_, err = io.Copy(file, captchaRes.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	var solve string
	fmt.Println("Etner the captcha")
	fmt.Scanln(&solve)

	var v = url.Values{}
	v.Set("name", username)
	v.Set("pass", password)
	v.Set("use_old_captcha", "1")
	v.Set("captcha", solve)
	v.Set("action", "login")

	lres, err := fa.client.PostForm(faLogin, v)
	if err != nil {
		log.Println(err)
		return err
	}
	lres.Body.Close()

	if err = httpError(lres); err != nil {
		log.Println(err)
		return err
	}

	return fa.isLoggedIn()
}

func (fa *furaffinity) SetCookies(cookies string) error {
	var c []*http.Cookie
	for _, cookie := range strings.Split(cookies, "\n") {
		spl := strings.Split(cookie, "=")
		if len(spl) != 2 {
			continue
		}
		c = append(c, &http.Cookie{Name: spl[0], Value: spl[1], Domain: "furaffinity.net"})
	}

	url, _ := url.Parse("https://furaffinity.net")
	fa.client.Jar.SetCookies(url, c)

	return fa.isLoggedIn()
}

func (fa *furaffinity) GetCookies() string {
	var cookies string
	url, err := url.Parse("https://furaffinity.net")
	if err != nil {
		log.Println(err)
		return ""
	}
	if fa.client.Jar == nil {
		log.Println("nil cookiejar")
		return ""
	}

	for _, cookie := range fa.client.Jar.Cookies(url) {
		cookies += cookie.String() + "\n"
	}
	return cookies
}

func (fa *furaffinity) isLoggedIn() error {
	res, err := fa.client.Get(faBase)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	if err = httpError(res); err != nil {
		log.Println(err)
		return err
	}

	var f func(*html.Node) bool
	f = func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "id" && a.Val == "my-username" {
					return true
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if f(c) {
				return true
			}
		}
		return false
	}

	node, err := html.Parse(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}

	if !f(node) {
		return errors.New("could not find username in html")
	}

	return nil
}
