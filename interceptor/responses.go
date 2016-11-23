// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/examples/util"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
)

var iface = flag.String("i", "eth0", "Interface to get packets from")
var fname = flag.String("r", "", "Filename to read from, overrides -i")
var snaplen = flag.Int("s", 1600, "SnapLen for pcap packet capture")
var filter = flag.String("f", "tcp and dst port 80", "BPF filter for pcap")
var logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")

// Build a simple HTTP request parser using tcpassembly.StreamFactory and tcpassembly.Stream interfaces

// httpStreamFactory implements tcpassembly.StreamFactory
type httpStreamFactory struct{}

// httpStream will handle the actual decoding of http requests.
type httpStream struct {
	net, transport gopacket.Flow
	r              tcpreader.ReaderStream
}

func (h *httpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &httpStream{
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
	}
	go hstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	for {
		resp, err := http.ReadResponse(buf, nil)
		if err == io.EOF {
			// We must read until we see an EOF... very important!
			return
		} else if err == io.ErrUnexpectedEOF {
			// TODO: need to establish if we get these in the header.
			//log.Println("UEOF IN RESP:", h.net, ":",, err)
			return
		} else if err != nil {
			log.Println("ERROR IN RESP:", h.net, ":", err)
			// TODO: What else?
		} else {

			contentType := resp.Header["Content-Type"]
			//contentEnc := resp.Header["Content-Encoding"]

			//log.Println("ENCODING:", resp.TransferEncoding, ":", contentEnc, ":", resp.Uncompressed)

			if len(contentType) != 0 {

				reader := resp.Body
				/*
					if len(contentEnc) != 0 {
						if contentEnc[0] == "gzip" {
							r, qerr := gzip.NewReader(resp.Body)
							if qerr != nil {
								log.Println("ERROR GZIP:", qerr)
							}
							reader = r
						}
					}
				*/

				switch contentType[0] {
				// TODO: ASCII, ANSI (Windows-1252)
				case "text/html", "text/html; charset=utf-8", "text/html; charset=UTF-8":
					// Default charset for HTML5
					log.Print("MATCHED:", contentType[0])

					b, err := ioutil.ReadAll(reader)
					if err != nil {
						log.Println(err)
					}
					fmt.Println(string(b))
					/*
						body, perr := html.Parse(resp.Body)
						if perr != nil {
							log.Println("PARSE ERROR:", perr)
							break
						} else {
							doc := goquery.NewDocumentFromNode(body)
							fmt.Println("DOC:", doc.Find("h1").Text())
						}
					*/
				case "text/html; charset=iso-8859-1", "text/html; charset=ISO-8859-1":
					// Default charset for HTML 2 to 4
					// TODO: Do something with it, e.g. convert with iconv.
					log.Print("MATCHED:", contentType[0])
					fallthrough
				default:
					//log.Println("UNUSED TYPE:", contentTyp
				}
			}

			for {
				bytes, err := tcpreader.DiscardBytesToFirstError(resp.Body)

				if err == io.EOF || err == io.ErrUnexpectedEOF {
					break

				} else if err != nil && bytes == 0 {
					log.Println("ERROR BUT ZERO:", h.net, ":", err, ":", bytes)
					break

				} else if err != nil {
					log.Println("ERROR NOT ZERO:", h.net, ":", err, ":", bytes)
				}
			}
			resp.Body.Close()
			//log.Println(h.net, ":", contentType, resp.Status)
		}
	}
}

func main() {
	defer util.Run()()
	var handle *pcap.Handle
	var err error

	// Set up pcap packet capture
	if *fname != "" {
		log.Printf("Reading from pcap dump %q", *fname)
		handle, err = pcap.OpenOffline(*fname)
	} else {
		log.Printf("Starting capture on interface %q", *iface)
		// TODO: Not sure about BlockForever.
		handle, err = pcap.OpenLive(*iface, int32(*snaplen), true, pcap.BlockForever)
	}
	if err != nil {
		log.Fatal(err)
	}

	if err := handle.SetBPFFilter(*filter); err != nil {
		log.Fatal(err)
	}

	// Set up assembly
	streamFactory := &httpStreamFactory{}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Println("reading in packets")
	// Read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Second * 10)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				return
			}
			if *logAllPackets {
				log.Println(packet)
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)

		case <-ticker:
			// Was: Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Second * -20))
		}
	}
}
