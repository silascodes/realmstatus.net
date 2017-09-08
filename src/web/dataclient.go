package main

import (
    "errors"
    "encoding/json"

    "common"

    "github.com/bradfitz/gomemcache/memcache"
    "encoding/binary"
)

type DataClient struct { 
    Memcache *memcache.Client
}

func NewDataClient(address string, port int) *DataClient {
    d := DataClient{
        memcache.New("localhost:11211"),
    }

    return &d
}

func (this *DataClient) GetLastUpdateTimestamp() (int64, error) {
    var ts int64

    raw, err := this.Memcache.Get("lastupdate")
    if err != nil {
        return ts, errors.New("error fetching last update timestamp (" + err.Error() + ")")
    }

    ts = int64(binary.LittleEndian.Uint64(raw.Value))
    return ts, nil
}

func (this *DataClient) GetAllRealms() ([]RealmView, error) {
    list := struct{
        Realms []RealmView
    }{}

    raw, err := this.Memcache.Get("realms")
    if err != nil {
        return list.Realms, errors.New("error fetching realm data (" + err.Error() + ")")
    }

    err = json.Unmarshal(raw.Value, &list)
    if err != nil {
        return list.Realms, errors.New("error decoding realm data")
    }

    return list.Realms, nil
}

func (this *DataClient) GetAllRegions() ([]common.Region, error) {
    list := []common.Region{}

    raw, err := this.Memcache.Get("regions")
    if err != nil {
        return list, errors.New("error fetching region data")
    }

    err = json.Unmarshal(raw.Value, &list)
    if err != nil {
        return list, errors.New("error decoding region data")
    }

    return list, nil
}

func (this *DataClient) GetAllLocales() ([]common.Locale, error) {
    list := []common.Locale{}

    raw, err := this.Memcache.Get("locales")
    if err != nil {
        return list, errors.New("error fetching locale data")
    }

    err = json.Unmarshal(raw.Value, &list)
    if err != nil {
        return list, errors.New("error decoding locale data")
    }

    return list, nil
}