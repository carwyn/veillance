import time
import functools
import httplib
import urllib2
import argparse

# GREEN = "10.42.0.104"
# PURPLE = "10.42.0.188"

class BoundHTTPHandler(urllib2.HTTPHandler):

    def __init__(self, source_address=None, debuglevel=0):
        urllib2.HTTPHandler.__init__(self, debuglevel)
        self.http_class = functools.partial(httplib.HTTPConnection,
                source_address=source_address)

    def http_open(self, req):
        return self.do_open(self.http_class, req)

def main():

    parser = argparse.ArgumentParser(description="Veillance Bot")
    parser.add_argument("--site", default="http://news.bbc.co.uk/")
    parser.add_argument("--source")
    parser.add_argument("--sleep", default=60, type=int)

    args = parser.parse_args()

    if args.source:
        print "Using Source IP:", args.source
        handler = BoundHTTPHandler(source_address=(args.source, 0))
        opener = urllib2.build_opener(handler)
        urllib2.install_opener(opener)
    
    while True:
        f = urllib2.urlopen(args.site)
        content = f.readlines()

        print f.getcode(), f.geturl(), "Length: ", len(content)
    
        time.sleep(args.sleep)

if __name__ == "__main__": main()

