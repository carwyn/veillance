from bs4 import BeautifulSoup, SoupStrainer
import cStringIO
from PIL import Image
from libmproxy.script import concurrent
from libmproxy.protocol.http import decoded
from pymongo import MongoClient

client = MongoClient("mongodb://localhost:3001/")

only_tags = [
    'a','p','ul','ol','dl','div','tt','code','blockquote','pre','table',
    'header','h1','h2','h3','h4','h5','h6'
]

only = SoupStrainer(only_tags)


#@concurrent
def response(context, flow):

    cont_type = flow.response.headers.get_first("content-type", "")

    if cont_type.startswith("image"):
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

    elif cont_type.startswith("text/html"):
        with decoded(flow.response):
            print "MATCHED: ", cont_type
            content = flow.response.content
            soup = BeautifulSoup(flow.response.content,"lxml",parse_only=only)
            cleaned = ' '.join(soup.get_text().split())
            if cleaned:
                print cleaned
                client.meteor.fragments.insert({'text': cleaned})


