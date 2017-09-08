package main

import (
	"time"
	"common"
	"net/http"
)

type ViewData struct {
	Title string
	Theme string
	RealmSlug string
	Realm *RealmView
	Year int
	Scripts []string
	Timestamp time.Time
	Regions []common.Region
	Locales []common.Locale
	Types []string
	Realms []RealmView
}

func NewViewData(dataClient *DataClient, req *http.Request) (*ViewData, error) {
	vd := ViewData{}

	vd.Year = time.Now().Year()

	timestamp, err := dataClient.GetLastUpdateTimestamp()
	if err != nil {
		return &vd, err
	}
	vd.Timestamp = time.Unix(0, timestamp)

	vd.Regions, err = dataClient.GetAllRegions()
	if err != nil {
		return &vd, err
	}

	vd.Locales, err = dataClient.GetAllLocales()
	if err != nil {
		return &vd, err
	}

	vd.Realms, err = dataClient.GetAllRealms()
	if err != nil {
		return &vd, err
	}

	theme, err := req.Cookie("realmstatus_net_theme")
	if err != nil {
		vd.Theme = "dark"
	} else {
		vd.Theme = theme.Value
	}

	return &vd, nil
}

func (this *ViewData) HasRealm() bool {
	if len(this.RealmSlug) > 0 && this.Realm != nil {
		return true
	}
	return false
}

func (this *ViewData) SetRealm(slug string) {
	this.RealmSlug = slug
	this.Realm = nil

	if len(this.RealmSlug) > 0 {
		for i, v := range this.Realms {
			if this.RealmSlug == v.GetRealSlug() {
				this.Realm = &this.Realms[i]
			}
		}
	}
}

func (this *ViewData) BattlegroupRealms() []*RealmView {
	results := []*RealmView{}

	for i, v := range this.Realms {
		if v.Battlegroup == this.Realm.Battlegroup {
			results = append(results, &this.Realms[i])
		}
	}

	return results
}

func (this *ViewData) ConnectedRealms() []*RealmView {
	results := []*RealmView{}
	slugs := make([]string, len(this.Realm.Connected_Realms) - 1)

	i := 0
	for _, v := range this.Realm.Connected_Realms {
		if v == this.Realm.Slug {
			continue;
		}
		slugs[i] = v + "-" + this.Realm.Region
		i++
	}

	for i, v := range this.Realms {
		for _, x := range slugs {
			if x == v.GetRealSlug() {
				results = append(results, &this.Realms[i])
			}
		}
	}

	return results
}