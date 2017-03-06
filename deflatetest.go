package main

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

var site = flag.String("site", "http://www.sueddeutsche.de/", "Site to fetch.")
var encoding = flag.String("encoding", "deflate", "Encodings to accept.")

func main() {
	flag.Parse()
	client := new(http.Client)

	//request, _ := http.NewRequest("GET", "http://www.bbc.co.uk/", nil)
	request, _ := http.NewRequest("GET", *site, nil)
	request.Header.Add("Accept-Encoding", *encoding)

	response, _ := client.Do(request)
	//body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	//var reader io.Reader
	var err error

	ce := response.Header.Get("Content-Encoding")
	fmt.Println("TEST:", ce, *encoding)

	switch ce {
	case "deflate":
		fmt.Println("DEFLATE")
		reader, err = zlib.NewReader(response.Body)
		if err == nil {
			break
		} else {
			fmt.Println(err)
		}
		reader = flate.NewReader(response.Body)
		//defer reader.Close()
	case "gzip":
		fmt.Println("GZIP")
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		//defer reader.Close()
	case "compress", "br":
		fmt.Println("UNSUPPORTED: %s", ce)
		fallthrough
	default:
		fmt.Println("IDENTITIY")
		//reader = bytes.Reader(body)
	}

	io.Copy(os.Stdout, reader) // print html to standard out
}
