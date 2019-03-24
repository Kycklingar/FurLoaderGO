package inkbunny

import (
	"fmt"
	"log"
	"strconv"

	"github.com/kycklingar/FurLoaderGO/dli"
)

func (ib *InkBunny) Watchlist() ([]dli.User, error) {
	v := ib.sidURLValues()
	res, err := client.PostForm(apiWatchlist, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	if err = httpError(res); err != nil {
		log.Println(err)
		return nil, err
	}

	var wl ibJsonWatchlist
	if err = wl.decode(res.Body); err != nil {
		log.Println(err)
		return nil, err
	}

	var users []dli.User
	for _, w := range wl.Watches {
		var u user
		u.name = w.Username
		u.id, err = strconv.Atoi(w.UserID)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		users = append(users, u)
	}

	return users, nil
}

func (ib *InkBunny) Feed(page int) ([]dli.Submission, error) {
	v := ib.sidURLValues()
	v.Set("page", fmt.Sprint(page))

	if rid := ib.ridMap[fmt.Sprint("FEED", page)]; rid != nil {
		v.Set("rid", *rid)
	} else {
		v.Set("unread_submissions", "yes")
		v.Set("get_rid", "yes")
	}

	res, err := client.PostForm(apiSearch, v)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	if err = httpError(res); err != nil {
		log.Println(err)
		return nil, err
	}

	var se ibJsonSearch
	if err = se.decode(res.Body); err != nil {
		log.Println(err)
		return nil, err
	}

	var subs []dli.Submission

	for _, sub := range se.Submissions {
		v, err := strconv.Atoi(sub.PageCount)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if v > 1 {
			s, err := ib.getFileUrls(sub.SubID)
			if err != nil {
				log.Println(err)
				return nil, err
			}

			r := make([]dli.Submission, len(s))
			for i, v := range s {
				r[i] = &v
			}
			subs = append(subs, r...)
		} else {
			var s ibSub
			s.user.name = sub.Username
			s.fileName = sub.FileName
			s.fileURL = sub.FileURL
			subs = append(subs, &s)
		}
	}

	return subs, nil
}
