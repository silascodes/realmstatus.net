package common

type RealmList struct {
    Realms []Realm
}

func (this RealmList) Merge(other *RealmList) RealmList {
    var list RealmList
    list.Realms = append(this.Realms, other.Realms...)
    return list
}