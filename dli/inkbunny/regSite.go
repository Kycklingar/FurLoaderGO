package inkbunny

import "github.com/kycklingar/FurLoaderGO/dli"

func init() {
	var i InkBunny
	dli.Logins = append(dli.Logins, &i)
	//dli.Watchers = append(dli.Watchers, i)
	//dli.Galleries = append(dli.Galleries, i)
}
