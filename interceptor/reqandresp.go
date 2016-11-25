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
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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

func isRequest(reader *bufio.Reader) bool {
	tp := bufio.NewReader(reader)
	firstLine, _ := tp.Peek(10)
	strFL := string(firstLine)
	arr := strings.Split(strFL, " ")

	switch strings.TrimSpace(arr[0]) {
	case "GET", "HEAD", "POST", "PUT", "DELETE", "TRACE", "CONNECT":
		return true
	default:
		return false
	}
}

func isResponse(reader *bufio.Reader) bool {
	tp := bufio.NewReader(reader)

	firstLine, err := tp.Peek(5)
	if err != nil {
		if err != io.EOF {
			log.Println("Err Peeking", err)
		}
	}
	strFL := string(firstLine)
	ret := strings.HasPrefix(strings.TrimSpace(strFL), "HTTP/")
	return ret
}

func (h *httpStream) run() {
	buf := bufio.NewReader(&h.r)
	var req *http.Request
	var err error
	for {
		if isRequest(buf) {
			log.Println("It's a request!")
			req, err = http.ReadRequest(buf)
			if err == io.EOF {
				// We must read until we see an EOF... very important!
				return
			} else if err != nil {
				log.Println("Error reading stream", h.net, h.transport, ":", err)
			} else {
				bodyBytes := tcpreader.DiscardBytesToEOF(req.Body)
				req.Body.Close()
				log.Println("Received request from stream", h.net, h.transport, ":", req, "with", bodyBytes, "bytes in request body")
			}
		} else if isResponse(buf) {
			log.Println("It's a response!")
			resp, err := http.ReadResponse(buf, req)
			if err == io.EOF {
				return
			} else if err != nil {
				log.Println("Error reading stream", h.net, h.transport, ":", err)
			} else {
				log.Println("Received resp from stream", h.net, h.transport, ":", resp, "with", " resp body: ", resp.Body, " content length:", resp.ContentLength, " transfer encoding: ", resp.TransferEncoding)
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Println("\nFailed to read response body, err: ", err)
				} else {
					log.Println("Reading response body..")
					bodystr := string(body)
					log.Println("resp body is: ", bodystr)
					resp.Body.Close()
				}

			}

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
	log.Println("Handle is: ", handle)
	// Read in packets, pass to assembler.
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			//log.Println("In select packet..")
			if packet == nil {
				log.Println("nil packet!")
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
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			log.Println("In ticker ..")
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
