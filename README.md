
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

 I've broken down
 the technical components into the following:

 # Back End 

 This part gathers data from the data collectors (see below) and makes
 them available to the front ends. This can run on the wifi gateway or
 in the cloud.


 # Data Collectors 

 ## Local Sniffing Wifi Gateway (Prototype Complete) 

 This is what I was talking about last week with Ronan, this can be
 installed at any exhibition initially. Later there will be a way for
 anyone to create one of these with as little as a Raspberry pi and a
 USB Wifi Dongle. I'd suggest this be the main thing we demo at FACT?

 Under the hood this is basically a Wifi Hotspot with a modified
 intercepting http proxy and DNS packet sniffer. This grabs any passing
 HTTP (not HTTPS) traffic and DNS requests.


 ##Twitter Scraper 

 Anyone that follows the Veillance twitter account will have their
 public tweets added to the pool of data.

 Anyone that signs in with Twitter will also have their public posts grabbed.

 ## Facebook Stalker 

 Anyone who friends the Veillance account will have their public
 visible posts grabbed.

 Anyone who signs in with FB will have anything they post grabbed (not messages).


 ## Advertising Scheme Snooping 

 These depend on being a member of an advertising cartel like Amazon or
 Google Ad Words. These don't watch your actions as much as you might
 think.

 ## Native Mobile Apps 

 IF we can get these built AND get them approved for the app stores
 (the hard part) they will grab any data the user has granted access
 to. I doubt the app stores will approve these though.


 # Front Ends 

 ## Web Visitor 

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


 ## Exhibition Mode 

 This is just a full screen version of the current playing animation.

