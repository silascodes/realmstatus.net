package common

type Realm struct {
    Region string
    Name string
    Slug string
    Battlegroup string
    Locale string
    Type string
    Population string
    Queue bool
    Status bool
    Timezone string
    Connected_Realms []string
}