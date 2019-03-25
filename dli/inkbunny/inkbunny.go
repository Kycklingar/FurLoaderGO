package inkbunny

import (
	"log"
	"net/http"
	"net/url"
)

type InkBunny struct {
	sid string

	ridMap map[string]string
}

var client http.Client

const (
	baseURL        = "https://inkbunny.net/"
	apiLogin       = baseURL + "api_login.php"
	apiSubmissions = baseURL + "api_submissions.php"
	apiSearch      = baseURL + "api_search.php"
	apiWatchlist   = baseURL + "api_watchlist.php"
)

func (ib *InkBunny) sidURLValues() url.Values {
	var v = url.Values{}
	v.Set("sid", ib.sid)
	return v
}

func (ib *InkBunny) Login(username, password string) error {
	v := url.Values{}
	v.Set("username", username)
	v.Set("password", password)

	res, err := client.PostForm(apiLogin, v)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	if err = httpError(res); err != nil {
		log.Println(err)
		return err
	}

	var m ibJsonLogin
	if err = m.decode(res.Body); err != nil {
		log.Println(err)
		return err
	}

	ib.sid = m.Sid

	return nil
}

func (ib *InkBunny) SetCookies(sid string) error {
	ib.sid = sid

	return ib.checkLogin()
}

func (ib *InkBunny) GetCookies() string {
	return ib.sid
}

func (ib *InkBunny) checkLogin() error {
	v := ib.sidURLValues()
	v.Set("no_submissions", "yes")

	res, err := client.PostForm(apiSubmissions, v)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	if err = httpError(res); err != nil {
		log.Println(err)
		return err
	}

	var m ibJsonSearch
	if err = m.decode(res.Body); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
