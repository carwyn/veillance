import urllib.request
import zlib

URL = "http://www.sueddeutsche.de/"

headers = {'Accept-Encoding': 'deflate'}


req = urllib.request.Request(URL, headers=headers)
resp = urllib.request.urlopen(req)

content = resp.read()

dobj = zlib.decompressobj(wbits=-zlib.MAX_WBITS)

print(dobj.decompress(content))
