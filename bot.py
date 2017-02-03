import urllib2
import time

BOT_SITE = "http://news.bbc.co.uk/"

SLEEP_TIME = 60

def main():
    
    while True:
        f = urllib2.urlopen(BOT_SITE)
        content = f.readlines()

        print f.getcode(), f.geturl(), "Length: ", len(content)
    
        time.sleep(SLEEP_TIME)

if __name__ == "__main__": main()

