package main

import (
    "strings"
    "net/http"
    "html/template"
    "github.com/googollee/go-socket.io"
    "time"
    "encoding/json"
)

type Controller struct { 
    *http.ServeMux
    socket *socketio.Server
    templates *template.Template
    dataClient *DataClient
}

func NewController() *Controller {
    c := Controller{
        http.NewServeMux(),
        nil,
        nil,
        NewDataClient("localhost", 11211),
    }

    // Load all our templates
    var err error
    c.templates, err = template.New("template.html").ParseFiles(
        "templates/template.html",
        "templates/list.html",
        "templates/single.html",
    )
    if err != nil {
        // TODO: handle it
    }

    // Set up the websocket handler
    c.socket, err = socketio.NewServer(nil)
    if err != nil {
        // TODO: handle it
    }
    c.socket.On("connection", func(so socketio.Socket) {
        so.Join("realms")
    })
    go c.HandleSocket()

    // Register server handlers for our dynamic methods and static files
    c.HandleFunc("/", c.HandleIndex)
    c.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
    c.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
    c.Handle("/socket.io/", c.socket)
    c.HandleFunc("/wow/", c.HandleList)
    c.HandleFunc("/wow/realm/", c.HandleSingle)

    return &c
}

type RealmStatus struct {
    Status bool
    Queue bool
}

func (this *Controller) generateStatusMap() (map[string]RealmStatus, error) {
    results := make(map[string]RealmStatus)
    realms, err := this.dataClient.GetAllRealms()
    if err != nil {
        return results, err
    }

    for _, v := range realms {
        results[v.GetRealSlug()] = RealmStatus{
            v.Status,
            v.Queue,
        }
    }

    return results, nil
}

func (this *Controller) generateStatusDiff(old map[string]RealmStatus, new map[string]RealmStatus) map[string]RealmStatus {
    results := make(map[string]RealmStatus)

    for i, v := range old {
        if v.Status != new[i].Status || v.Queue != new[i].Queue {
            results[i] = new[i]
        }
    }

    return results
}

func (this *Controller) HandleSocket() {
    lastUpdate, _ := this.dataClient.GetLastUpdateTimestamp()
    lastData, _ := this.generateStatusMap()
    for ;; {
        timestamp, _ := this.dataClient.GetLastUpdateTimestamp()
        if(timestamp > lastUpdate) {
            lastUpdate = timestamp
            newData, _ := this.generateStatusMap()
            data := this.generateStatusDiff(lastData, newData)
            encoded, _ := json.Marshal(data)
            this.socket.BroadcastTo("realms", string(encoded))
        }

        time.Sleep(time.Second * 5)
    }
}

func (this *Controller) HandleIndex(w http.ResponseWriter, req *http.Request) {
    // If not actually requesting index, give not found
    if req.URL.Path != "/" {
        http.NotFound(w, req)
        return
    }

    // Redirect to the /wow/ page
    http.Redirect(w, req, "/wow/", http.StatusFound)
}

func (this *Controller) HandleList(w http.ResponseWriter, req *http.Request) {
    // Load the view data
    data, err := NewViewData(this.dataClient, req)
    if err != nil {
        http.Error(w, "data loading error (" + err.Error() + ")", http.StatusInternalServerError)
        return
    }

    // Set page title
    data.Title = "Live WoW Realm Status - Search"

    // Add list page javascript
    data.Scripts = make([]string, 2)
    data.Scripts[0] = "/js/lib/jquery.tablesorter.min.js"
    data.Scripts[1] = "/js/list.js"

    // Render the page
    err = this.templates.Execute(w, data)
    if err != nil {
        http.Error(w, "error rendering page (" + err.Error() + ")", http.StatusInternalServerError)
        return
    }
}

func (this *Controller) HandleSingle(w http.ResponseWriter, req *http.Request) {
    // Load the view data
    data, err := NewViewData(this.dataClient, req)
    if err != nil {
        http.Error(w, "data loading error (" + err.Error() + ")", http.StatusInternalServerError)
        return
    }

    // Set the Realm property for single view mode
    data.SetRealm(strings.TrimPrefix(req.URL.Path, "/wow/realm/"))

    // Set page title
    data.Title = data.Realm.Name + " " + data.Realm.GetRegion() + " Realm Status"

    // Add single page javascript
    data.Scripts = make([]string, 3)
    data.Scripts[0] = "/js/lib/moment.min.js"
    data.Scripts[1] = "/js/lib/moment-timezone-with-data-2010-2020.min.js"
    data.Scripts[2] = "/js/single.js"

    // Render the page
    err = this.templates.Execute(w, data)
    if err != nil {
        http.Error(w, "error rendering page (" + err.Error() + ")", http.StatusInternalServerError)
        return
    }
}
