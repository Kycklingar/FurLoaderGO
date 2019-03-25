package inkbunny

import (
	"fmt"
	"log"

	"github.com/kycklingar/FurLoaderGO/dli"
)

func (ib *InkBunny) Posts(userID string, offset int) ([]dli.Submission, error) {
	v := ib.sidURLValues()

	rid, ok := ib.ridMap["gallery "+userID]
	if ok {
		v.Set("rid", rid)
	}
	v.Set("username", userID)
	v.Set("page", fmt.Sprint(offset+1))

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

	var j ibJsonSearch
	if err = j.decode(res.Body); err != nil {
		log.Println(err)
		return nil, err
	}

	var subs []dli.Submission

	for _, jsub := range j.Submissions {
		s, err := ib.fromJson(jsub)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		dsubs := make([]dli.Submission, len(s))
		for i, sub := range s {
			dsubs[i] = &sub
		}

		subs = append(subs, dsubs...)
	}

	return subs, nil
}
