package main

import (
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	client := new(http.Client)

	//request, _ := http.NewRequest("GET", "http://www.bbc.co.uk/", nil)
	request, _ := http.NewRequest("GET", "http://www.sueddeutsche.de/", nil)
	request.Header.Add("Accept-Encoding", "deflate")

	response, _ := client.Do(request)

	defer response.Body.Close()

	// Check that the server actually sent compressed data
	var reader io.ReadCloser
	var err error
	switch response.Header.Get("Content-Encoding") {
	case "flate":
		fmt.Println("FLATE")
		reader = flate.NewReader(response.Body)
		//defer reader.Close()
	case "deflate":
		fmt.Println("DEFLATE")
		reader, err = zlib.NewReader(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		//defer reader.Close()
	case "gzip":
		fmt.Println("GZIP")
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			fmt.Println(err)
		}
		//defer reader.Close()
	default:
		fmt.Println("OTHER")
		reader = response.Body
	}

	io.Copy(os.Stdout, reader) // print html to standard out
}
