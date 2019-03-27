package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kycklingar/FurLoaderGO/data"
	"github.com/kycklingar/FurLoaderGO/dli"
	_ "github.com/kycklingar/FurLoaderGO/dli/inkbunny"
)

var db *data.DB

func main() {
	log.SetFlags(log.Llongfile)
	username := flag.String("username", "", "Your username")
	password := flag.String("password", "", "Your password")
	cookies := flag.String("cookies", "", "Use instead of logging in")
	page := flag.Int("page", 0, "Start the search from this page")
	user := flag.String("user", "", "Gallery of user you want to download from")

	flag.Parse()

	var ibl = dli.Logins["inkbunny"]
	var ibg = dli.Galleries["inkbunny"]

	db = data.OpenDB()

	if *cookies != "" {
		err := ibl.SetCookies(*cookies)
		if err != nil {
			log.Fatal(err)
		}
	} else if *username != "" {
		err := ibl.Login(*username, *password)
		if err != nil {
			log.Fatal(err)
		}
		db.Store("cookies:inkbunny", ibl.GetCookies())
	} else {
		// Use stored cookkies if exists
		cookies := db.Get("cookies:inkbunny")
		if cookies == "" {
			log.Fatal("no cookies, pleases login first")
		}
		ibl.SetCookies(cookies)
	}

	fmt.Println(ibl.GetCookies())
	if len(ibl.GetCookies()) <= 0 {
		log.Fatal("no cookies gotten")
	}

	if len(*user) <= 0 {
		log.Fatal("No user gallery specified")
	}

	queue := Queue()
	go queue.startThread()
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
