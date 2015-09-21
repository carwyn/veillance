
# Veillance

The objective of this project is to create awareness of the pervasiveness of
these invisible tracking territories by rendering visible the landscapes of data
surveillance and, through an interactive typographic artwork, to re-imagine this
phenomenon creatively. In defined opposition to the surveillance idea of the
individual as ‘target’ or ‘consumer’, this artwork will suggest a new model for
experiencing the intersection between surveillance and ‘private’ everyday action
through the audience’s informed re-assumption of ownership over their personal
information, and the re-assertion of everyday creativity and informed agency as
we move through and engage with our public spaces.

## Prerequisites

To run the Veillance front end you must have a modern web browser that
supports WebGL. Note that there can be a small delay as the 3D image is
initially generated.

## Running The Development Version

The current front and back ends are combined as a single http://meteor.com/
applciation. To install meteor follow the instructions as per the Meteor
home page, namely:

```
$ curl https://install.meteor.com/ | sh
```

Then to start Veillance simply run the following in the `veilweb` directory.

```
$ meteor run
```

You can then connect to the front end by pointing a web browser at:

```
http://localhost:3000/
```

To start collecting data, you then need to run one or more of the
***Data Collectors*** mentioned below.

### Removing All Captured Data

To remove all existing captured data from the MongoDB database use the Meteor
MongoDB shell to get a MongoDB prompt. In the veilweb directory:

```
$ meteor mongo
```

You can then delete all the captured fragments using the following:

```
meteor:PRIMARY> db.fragments.remove({})
```

This will remove all existing entries.

## Technical Components

### Data Aggregator

This part gathers data from the data collectors (see below) and makes
them available to the front ends. This can run on the wifi gateway or
in the cloud. In the prototype this is effectively a MongoDB database.


### Data Collectors 

Currently the collectors all work by writing directly into the MongoDB
database that forms part of the Meteor application. This will change in the
future to allow distribution of the data collectors and the 

#### Intercepting HTTP Proxy (Complete)

Under the hood this is basically a modified intercepting http proxy. This
grabs any passing HTTP (not HTTPS) traffic stripping out and capturing the
plain text in the web pages requested.

This collector uses `mitmproxy` (http://mitmproxy.org/) and an inline script.

To install `mitmproxy` and all the requirement the collector script needs run
the following on OS X or Linux (you can also pull in some of the requiremetns
with Brew or RPM, see notes in `requirements.txt`):


```
$ pip install -r requirements.txt
```

Then run either `mitmproxy` or `mitmdump` using the http interceptor script:

```
$ mitmdump -q --stream 250k -s collectors/proxy/intercept.py
```

For development purposes you can now configure your browser to use localhost
on port 8080 as a HTTP proxy to test web page interception. Simply browse the
web with the browser you've configured to use the proxy.

***It's easier to use a different browser for capture to the one you are
using to look at the visualization on port 3000 when developing.*** For example
I use Chrome to view the visualization and Firefox to browse via the proxy.


#### Snooping DNS Server

TODO.


#### Twitter Scraper 

Anyone that follows the Veillance twitter account will have their
public tweets added to the pool of data.

Anyone that signs in with Twitter will also have their public posts grabbed.

#### Facebook Stalker 

Anyone who friends the Veillance account will have their public
visible posts grabbed.

Anyone who signs in with FB will have anything they post grabbed (not messages).


#### Advertising Scheme Snooping 

These depend on being a member of an advertising cartel like Amazon or
Google Ad Words. These don't watch your actions as much as you might
think.

#### Mobile Apps 

IF we can get these built AND get them approved for the app stores
(the hard part) they will grab any data the user has granted access
to. I doubt the app stores will approve these though.


## Front Ends 

### Web Visitor 

Available to desktop or mobile devices.

This is the public/kiosk web site that has a view of the current
animated text on it. Visitors will be able to sign up using Twitter or
Facebook hence allowing their respective data collectors (see above).

Note that there is a VERY limited amount of data that we can grab from
mobile browsers, for example:

Camera and Microphone:
http://www.html5rocks.com/en/tutorials/getusermedia/intro/

Geolocation:
http://www.html5rocks.com/en/tutorials/geolocation/trip_meter/

.. and only after a browser popup.


### Exhibition Mode 

This is just a full screen version of the current playing animation.

