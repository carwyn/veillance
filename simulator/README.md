

# Server Sent Events (SSE) also known as EventSource

Server Sent Events (SSE, EventSource are really just long running HTTP
connections (similar to downloading a infinite file) that follow a specific
newline delimited format. See the following for further details:

https://en.wikipedia.org/wiki/Server-sent_events
http://html5doctor.com/server-sent-events/
https://www.html5rocks.com/en/tutorials/eventsource/basics/


## Data Format

The data sent out via http://localhost:8080/events is a HTTP Server Sent Events
stream also known as an EventSource stream. This is just a continual stream of
data written to a standard HTTP connection with a simple format. Unlike a
standard HTTP connection though it's designed to be kept open long term so that
the client can continue to receive data over time. The formal for Veillance is:

`id: 123` # This is the message number since start of the stream NOT fragment ID.
`event:` fragment # This is the message type, one of fragment, domain, selection
`data: {}` # JSON formated data payload. This should be parsed as JSON.


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


