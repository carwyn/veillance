import time
import functools
import httplib
import urllib2

#BOT_SITE = "http://news.bbc.co.uk/"
BOT_SITE = "http://news.bbc.co.uk/"

#SOURCE_IP = "10.42.0.104"
SOURCE_IP = "10.42.0.188"

SLEEP_TIME = 60

class BoundHTTPHandler(urllib2.HTTPHandler):

    def __init__(self, source_address=None, debuglevel=0):
        urllib2.HTTPHandler.__init__(self, debuglevel)
        self.http_class = functools.partial(httplib.HTTPConnection,
                source_address=source_address)

    def http_open(self, req):
        return self.do_open(self.http_class, req)


def main():

    handler = BoundHTTPHandler(source_address=(SOURCE_IP, 0))
    opener = urllib2.build_opener(handler)
    urllib2.install_opener(opener)
    
    while True:
        f = urllib2.urlopen(BOT_SITE)
        content = f.readlines()

        print f.getcode(), f.geturl(), "Length: ", len(content)
    
        time.sleep(SLEEP_TIME)

if __name__ == "__main__": main()

