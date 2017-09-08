package main

import (
    "fmt"
    "time"
    "sort"
    "errors"
    "io/ioutil"
    "net/http"
    "encoding/json"

    "common"

    "github.com/bradfitz/gomemcache/memcache"
    "encoding/binary"
)

func FetchRealmList(region string, apiKey string) (common.RealmList, error) {
    fmt.Printf("FetchRealmList(region: \"" + region + "\", apiKey: *)\n")

    var realms common.RealmList

    resp, err := http.Get("https://" + region + ".api.battle.net/wow/realm/status?&apikey=" + apiKey)
    if err != nil {
        return realms, errors.New("error calling bnet api")
    }

    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return realms, errors.New("error reading response body")
    }

    err = json.Unmarshal(body, &realms)
    if err != nil {
        return realms, errors.New("error decoding response")
    }

    for i := range realms.Realms {
        v := &realms.Realms[i]
        v.Region = region
    }

    return realms, nil
}

func StoreRealmList(realms *common.RealmList) error {
    mc := memcache.New("localhost:11211")

    raw, err := json.Marshal(realms)
    if err != nil {
        return errors.New("error encoding realm data")
    }

    err = mc.Set(&memcache.Item{Key: "realms", Value: raw})
    if err != nil {
        return errors.New("error storing realm data")
    }

    return nil
}

func StoreRegionList(regions []common.Region) error {
    mc := memcache.New("localhost:11211")

    raw, err := json.Marshal(regions)
    if err != nil {
        fmt.Printf("error encoding region data")
        return errors.New("error encoding region data")
    }

    err = mc.Set(&memcache.Item{Key: "regions", Value: raw})
    if err != nil {
        fmt.Printf("error storing region data")
        return errors.New("error storing region data")
    }

    return nil
}

func StoreLocaleList(locales []common.Locale) error {
    mc := memcache.New("localhost:11211")

    raw, err := json.Marshal(locales)
    if err != nil {
        return errors.New("error encoding locale data")
    }

    err = mc.Set(&memcache.Item{Key: "locales", Value: raw})
    if err != nil {
        return errors.New("error storing locale data")
    }

    return nil
}

func ParseAndStoreLocales(realms *common.RealmList) error {
    lmap := make(map[string]struct{})

    for _, v := range realms.Realms {
        _, ok := lmap[v.Locale]
        if !ok {
            lmap[v.Locale] = struct{}{}
        }
    }

    locales := make([]common.Locale, len(lmap))
    i := 0
    for k, _ := range lmap {
        locales[i] = common.Locale(k)
        i++
    }

    err := StoreLocaleList(locales)
    if err != nil {
        return err
    }

    return nil
}

func StoreLastUpdateTimestamp() error {
    mc := memcache.New("localhost:11211")

    b := make([]byte, 8)
    binary.LittleEndian.PutUint64(b, uint64(time.Now().UnixNano()))

    err := mc.Set(&memcache.Item{Key: "lastupdate", Value: b})
    if err != nil {
        fmt.Printf("error storing last update timestamp")
        return errors.New("error storing last update timestamp")
    }

    return nil
}

type RealmsByName []common.Realm

func (this RealmsByName) Len() int           { return len(this) }
func (this RealmsByName) Swap(i, j int)      { this[i], this[j] = this[j], this[i] }
func (this RealmsByName) Less(i, j int) bool { return this[i].Name < this[j].Name }

func UpdateData(apiKey string) {
    fmt.Printf("UpdateData(apiKey: *)\n")

    regions := []common.Region{
        "us",
        "eu",
    }
    go StoreRegionList(regions)

    usRealms, usErr := FetchRealmList("us", apiKey)
    if usErr != nil {
        fmt.Printf("error fetching us realms: %s\n", usErr)
    }

    euRealms, euErr := FetchRealmList("eu", apiKey)
    if euErr != nil {
        fmt.Printf("error fetching eu realms: %s\n", euErr)
    }

    if usErr != nil && euErr != nil {
        fmt.Printf("failed to fetch any realms\n")
        return
    }

    realms := usRealms.Merge(&euRealms)
    sort.Sort(RealmsByName(realms.Realms))

    go ParseAndStoreLocales(&realms)

    go StoreRealmList(&realms)

    go StoreLastUpdateTimestamp()
}

func main() {
    interval := 60
    apiKey := "db4cbpaavc55rsddbcdbs8nbxug3ndgp"

    for true {
        go UpdateData(apiKey)
        time.Sleep(time.Second * time.Duration(interval))
    }
}
