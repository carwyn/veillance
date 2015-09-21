from bs4 import BeautifulSoup, SoupStrainer
import cStringIO
#from PIL import Image
from libmproxy.script import concurrent
from libmproxy.protocol.http import decoded
from pymongo import MongoClient
from urlparse import urlparse

client = MongoClient("mongodb://localhost:3001/")

only_tags = [
    #'a','p','ul','ol','dl','tt','code','blockquote','pre',
    'h1','h2','h3','h4','h5','h6'
]

strainer = SoupStrainer(only_tags)

@concurrent
def response(context, flow):
   
    # flow contains:
    #
    # request: HTTPRequest object
    # response: HTTPResponse object
    # error: Error object
    # server_conn: ServerConnection object
    # client_conn: ClientConnection object

    # flow.response contains:
    #
    # status_code: HTTP response status code
    # msg: HTTP response message
    # headers: Headers object
    # content: Content of the request, None, or CONTENT_MISSING if there
    # is content associated, but not present. CONTENT_MISSING evaluates
    # to False to make checking for the presence of content natural.

    if flow.response.code != 200:
        return
 
    cont_type = flow.response.headers.get_first("content-type", "")
   
    if cont_type.startswith("text/html"):

        if flow.response.content:
            print "CAPTURED"
            content = flow.response.get_decoded_content()
            soup = BeautifulSoup(content, "lxml", parse_only=strainer)

            cleaned = soup.get_text(' ', strip=True)
            if cleaned:
                print cleaned
                client.meteor.fragments.insert({'text': cleaned})

    elif cont_type.startswith("image"):

        # Flip all images upsidedown.
        pass
        #with decoded(flow.response):
        #    try:
        #        s = cStringIO.StringIO(flow.response.content)
        #        img = Image.open(s).rotate(180)
        #        s2 = cStringIO.StringIO()
        #        img.save(s2, "png")
        #        flow.response.content = s2.getvalue()
        #        flow.response.headers["content-type"] = ["image/png"]
        #    except:  # Unknown image types etc.
        #        pass

 
