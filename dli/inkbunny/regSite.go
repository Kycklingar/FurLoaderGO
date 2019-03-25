package inkbunny

import "github.com/kycklingar/FurLoaderGO/dli"

var inkbunny InkBunny

func init() {
	dli.Logins = append(dli.Logins, &inkbunny)
	//dli.Watchers = append(dli.Watchers, &inkbunny)
	dli.Galleries = append(dli.Galleries, &inkbunny)
}
