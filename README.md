# realmstatus.net
A little site to show live WoW realm states, see http://realmstatus.net

Pulls data from Battle.net REST API. Two modules, one to grab and parse data, one to serve website with a memcached instance in the middle. 

Uses socket.io for live updates, has light and dark themes can be switched on the client. 

Needs the data grabbing module (its in a rough state ATM) rewriting so can support filtering the list view, also some errors aren't being handled. Some comments probably wouldn't go astray too!
