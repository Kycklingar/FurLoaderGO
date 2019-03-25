package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/kycklingar/FurLoaderGO/dli"
	_ "github.com/kycklingar/FurLoaderGO/dli/inkbunny"
)

func main() {
	log.SetFlags(log.Llongfile)
	username := flag.String("username", "", "Your username")
	password := flag.String("password", "", "Your password")
	cookies := flag.String("cookies", "", "Use instead of logging in")
	user := flag.String("user", "", "Gallery of user you want to download from")

	flag.Parse()

	var ibl = dli.Logins[0]
	var ibg = dli.Galleries[0]

	if *cookies != "" {
		err := ibl.SetCookies(*cookies)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		err := ibl.Login(*username, *password)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(ibl.GetCookies())
	if len(ibl.GetCookies()) <= 0 {
		log.Fatal("no cookies gotten")
	}

	if len(*user) <= 0 {
		log.Fatal("No user gallery specified")
	}

	posts, err := ibg.Posts(*user, 0)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d posts in %s gallery", len(posts), *user)

	q := Queue(posts)
	q.start()

}
