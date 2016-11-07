
# How to Run The Veillance Data Capture Simulartor

* Install Go from https://golang.org/
* Set the GOPATH environment variable to $HOME/gocode or $HOME (up to you).
* Test Go using `go version` from the command line.
* Use the command "go get" to fetch the following dependencies:
  * `go get github.com/antage/eventsource`
* Build and run the simulator:
  * `go run simulator.go`
* Test the simulator using the command line CURL client:
  * `curl http://localhost:8080/events`

.. this will then start printing out the live stream of sample data.

# Data Format

The data sent out via http://localhost:8080/events is a HTTP Server Sent Events
stream also known as an EventSource stream. This is just a continual stream of
data written to a standard HTTP connection with a simple format. Unlike a
standard HTTP connection though it's designed to be kept open long term so that
the client can continue to receive data over time. The formal for Veillance is:

`id: 123` # This is the message number since start of the stream NOT fragment ID.
`event:` fragment # This is the message type, one of fragment, domain, selection
`data: {}` # JSON formated data payload. This should be parsed as JSON.

## JSON Payload

> WARNING: This element will evolve over time!

Fragment JSON: `{"Id":123,"Name":"Bob","Type":"User","Text":"Hello World!"}`

```
Id: Fragment ID of se
Name: Identity of originating data source.
Type: Type of oringinating data source.
Text: The captured text.
```

Selection JSON: `{"Fid":31,"Words":[1,3]}`

```
Fid: Fragment ID.
Words: Selected words.
```

Domain JSON: `{"Id":321,"Name":"Bill","Type":"User","Text":"twitter.com"}`

```
Id: Seen domain identifier.
Name: Identity of originating data source.
Type: Type of oringinating data source.
Text: The captured text.
```

Command JSON: TODO!


## Example Veillance Messages

```
id: 207
event: fragment
data: {"Id":16,"Name":"SimUser","Type":"User","Text":"Well the barrel broke my shoulder as on to the ground it sped"}

id: 206
event: selection
data: {"Fid":31,"Words":[1,3]}

id: 208
event: domain
data: {"Id":12,"Name":"SimUser","Type":"User","Text":"twitter.com"}
```


# Server Sent Events (SSE) also known as EventSource

Server Sent Events (SSE, EventSource are really just long running HTTP
connections (similar to downloading a infinite file) that follow a specific
newline delimited format. See the following for further details:

https://en.wikipedia.org/wiki/Server-sent_events
http://html5doctor.com/server-sent-events/
https://www.html5rocks.com/en/tutorials/eventsource/basics/


## Go Implementations

https://github.com/donovanhide/eventsource
https://github.com/kljensen/golang-html5-sse-example
https://github.com/antage/eventsource
https://github.com/stuartnelson3/golang-eventsource


## C++ Implementations

One of these two are recommended:

https://pocoproject.org/
https://curl.haxx.se/

OpenFrameworks ofURLFileLoader which uses libcurl interanlly might work, it
depends if it can let you not close the connection:

http://openframeworks.cc/documentation/utils/ofURLFileLoader/

Another one worth a look at is:

https://pocoproject.org/docs/package-Net.HTTPClient.html

For more options see:

https://curl.haxx.se/libcurl/competitors.html
http://kukuruku.co/hub/cpp/a-cheat-sheet-for-http-libraries-in-c


# JSON

POCO also has a JSON parser which you'll also need for the data payload:

https://pocoproject.org/docs/Poco.JSON.html


