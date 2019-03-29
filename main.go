package main

import (
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
	cookies := flag.String("cookies", "", "Use instead of logging in")
	page := flag.Int("page", 0, "Start the search from this page")
	user := flag.String("user", "", "Gallery of user you want to download from")

	flag.Parse()

	if *site == "" {
		log.Fatal("no site specified")
	}

	var ibl = dli.Logins[*site]
	var ibg = dli.Galleries[*site]

	db = data.OpenDB()

	if ibl == nil {
		log.Fatalf("%s does not implement Login", *site)
	}
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
		cookies := ibl.GetCookies()
		//fmt.Println(cookies)
		db.Store(fmt.Sprintf("cookies:%s", *site), cookies)
		os.Exit(1)
	} else {
		// Use stored cookkies if exists
		cookies := db.Get(fmt.Sprintf("cookies:%s", *site))
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

	if ibg == nil {
		log.Fatalf("%s doesn't implement Gallery", *site)
	}

	queue := Queue()
	go queue.startThread()
	//go queue.startThread()

	queue.addIncDL(func(i int) []dli.Submission {
		posts, err := ibg.Posts(*user, i+*page-1)
		if err != nil {
			log.Println(err)
			return nil
		}

		return posts
	})
}
