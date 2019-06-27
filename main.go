package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kycklingar/FurLoaderGO/data"
	"github.com/kycklingar/FurLoaderGO/dli"
	_ "github.com/kycklingar/FurLoaderGO/dli/furaffinity"
	_ "github.com/kycklingar/FurLoaderGO/dli/inkbunny"
)

var db *data.DB

func main() {
	log.SetFlags(log.Llongfile)
	site := flag.String("site", "", "Website")
	username := flag.String("username", "", "Your username")
	password := flag.String("password", "", "Your password")
	//cookies := flag.String("cookies", "", "Use instead of logging in")
	page := flag.Int("page", 0, "Start the search from this page")
	user := flag.String("user", "", "Gallery of user you want to download from")
	feed := flag.Bool("feed", false, "Download your feed")
	ignoreLatestFeedPost := flag.Bool("ignore-lfp", false, "Ignore the latest feed post in storage")
	feedMaxPages := flag.Int("feed-max", 5, "Max pages to download from your feed")
	login := flag.Bool("login", false, "Perform a login")

	flag.Parse()

	if *site == "" {
		fmt.Println("no site specified")
		fmt.Println("Sites available:")
		for key, _ := range dli.Galleries {
			fmt.Println(key)
		}
		os.Exit(0)
	}

	db = data.OpenDB()
	defer db.CloseDB()

	if *login {
		err := loginUser(*site, *username, *password)
		if err != nil {
			fmt.Println(err)
			return
		}
		return
	}

	if *feed {
		fmt.Println("Using stored cookies")
		if err := useStoredCookies(*site); err != nil {
			log.Println(err)
			return
		}

		dlw, ok := dli.Watchers[*site]
		if !ok {
			fmt.Println(&implementError{*site, "Watcher"})
			return
		}

		feed := dlw.Feed()

		queue := Queue()
		go queue.startThread()

		var latestFeedPost string

		if !*ignoreLatestFeedPost {
			latestFeedPost = db.Get(fmt.Sprintf("LatestFeedPost:%s", *site))
			fmt.Println("LFP: ", latestFeedPost)
		} else {
			fmt.Println("Ignoring LFP")
		}

		var lfpFound bool
		var newLfp string

		queue.addIncDL(func(i int) []dli.Submission {
			if lfpFound || (*feedMaxPages > 0 && i > *feedMaxPages) {
				return nil
			}

			fmt.Println("getting subs")

			subs, err := feed.NextPage()
			if err != nil {
				log.Println(err)
				return nil
			}
			fmt.Println(subs)

			for _, sub := range subs {
				if sub.ID() == latestFeedPost {
					fmt.Println("LFP found, stopping next batch")
					lfpFound = true
				}
			}

			if i == 0 && len(subs) > 0 {
				newLfp = subs[0].ID()
			}

			return subs
		})

		if newLfp != "" {
			db.Store(fmt.Sprintf("LatestFeedPost:%s", *site), newLfp)
		}

		return
	}

	if len(*user) >= 0 {
		if err := useStoredCookies(*site); err != nil {
			log.Println(err)
			return
		}

		ibg, ok := dli.Galleries[*site]
		if !ok {
			fmt.Println(&implementError{*site, "Gallery"})
			return
		}

		queue := Queue()
		go queue.startThread()

		queue.addIncDL(func(i int) []dli.Submission {
			posts, err := ibg.Posts(*user, i+*page-1)
			if err != nil {
				log.Println(err)
				return nil
			}

			return posts
		})
	}

}

type implementError struct {
	site   string
	interf string
}

func (ie *implementError) Error() string {
	return fmt.Sprintf("%s does not implement the %s interface", ie.site, ie.interf)
}

func useStoredCookies(siteName string) error {
	site, ok := dli.Logins[siteName]
	if !ok {
		return &implementError{siteName, "Login"}
	}

	cookies := db.Get(fmt.Sprintf("cookies:%s", siteName))

	if err := site.SetCookies(cookies); err != nil {
		return err
	}

	return nil
}

func loginUser(siteName, username, password string) error {
	site, ok := dli.Logins[siteName]
	if !ok {
		return &implementError{siteName, "Login"}
	}

	if username == "" {
		return errors.New("Empty username in login")
	}

	err := site.Login(username, password)
	if err != nil {
		return err
	}

	cookies := site.GetCookies()
	db.Store(fmt.Sprintf("cookies:%s", siteName), cookies)

	//}else {
	//	// Use stored cookkies if exists
	//	cookies := db.Get(fmt.Sprintf("cookies:%s", *site))
	//	if cookies == "" {
	//		log.Fatal("no cookies, pleases login first")
	//	}
	//	ibl.SetCookies(cookies)
	//}

	//fmt.Println(ibl.GetCookies())
	//if len(ibl.GetCookies()) <= 0 {
	//	log.Fatal("no cookies gotten")
	//}

	return nil
}

func startFeedDownload() {

}
