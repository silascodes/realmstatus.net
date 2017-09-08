package main

import (
    "strings"

    "golang.org/x/text/language"
    "golang.org/x/text/language/display"

    "common"
)

type RealmView struct {
    common.Realm
}

var even bool
func (this RealmView) GetClass() string {
    if !even {
        even = true;
        return "highlight"
    } else {
        even = false;
        return ""
    }
}

func (this RealmView) GetRealSlug() string {
    return this.Slug + "-" + this.Region
}

func (this RealmView) GetRegion() string {
    return strings.ToUpper(this.Region)
}

func (this RealmView) GetPopulation() string {
    return strings.Title(this.Population)
}

func (this RealmView) GetLocale() string {
    en := display.English.Languages()
    return en.Name(language.MustParse(this.Locale))
}

func (this RealmView) GetType() string {
    var str string
    for _, v := range this.Type {
        if v != 'v' {
            str += strings.ToUpper(string(v))
        } else {
            str += string(v)
        }
    }
    return str
}

func (this RealmView) GetQueue() string {
    if this.Queue {
        return "YES"
    } else {
        return "NO"
    }
}

func (this RealmView) GetStatus() string {
    if this.Status {
        return "UP"
    } else {
        return "DOWN"
    }
}

func (this RealmView) GetTimezone() string {
    return this.Timezone
}

func (this RealmView) HasConnectedRealms() bool {
    if len(this.Connected_Realms) > 1 {
        return true
    }
    return false
}