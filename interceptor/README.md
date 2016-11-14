
# Instructions

* Install Go as for simulator.
* go get golang.org/x/net/websocket
* git clone git@github.com:BangorUniversity/veillance.git
OR if you don't have SSH working:
* git clone https://github.com/BangorUniversity/veillance.git
* cd interceptor
* go build
* ./server

# Connecting to the Websocket Server

* Open a websocket to http://localhost:8080/entry
* Messages are sent to client as UTF-8 encoded JSON strings.
* Currently these follow the same format as the simulator "data" fields.

