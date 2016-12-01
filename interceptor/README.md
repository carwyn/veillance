
# How to Run The Veillance Data Capture Simulator

## Install and Test Go

* Install Go from https://golang.org/
* Set the GOPATH environment variable to $HOME/gocode or $HOME (up to you).
* Test Go using `go version` from the command line.

## Make Sure libpcap In Installed and pcap.h is in your INCLUDE path.

On Linux this is in the libpcap/libpcap-devel packages.
On OS X I'm not sure where this is but should be installed if XCode is.

## Fetch and Build Interceptor

* Use the command "go get" to fetch everything including dependencies:
  * `go get github.com/carwyn/veillance/interceptor`
  * This will pull down all the dependencies into the $GOPATH/src tree.
* Build by running the following in the veillance/interceptor directory:
  * `go build`
* Run the interceptor from the veillance/interceptor directory:
  * `sudo ./interceptor -i eth0 -nooptcheck -ignorefsmerr "tcp port 80"`
  * NOTE: Change the network interface passed to `-i` to match your system.
* Test the interceptor via the test web page:
  * `curl http://localhost:8080/test.html`
  * This will print a live stream to the Javascript console.


## Connecting to the Websocket Server Programatically

* Open a websocket to http://localhost:8080/entry
* Messages are sent to client as UTF-8 encoded JSON strings.
* Format is as described below.


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

