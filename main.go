package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"time"

	manet "github.com/jbenet/go-multiaddr-net"
	ma "github.com/jbenet/go-multiaddr-net/Godeps/_workspace/src/github.com/jbenet/go-multiaddr"
	udt "github.com/jbenet/go-multiaddr-net/vendor/go-udtwrapper-v1.0.0/udt"
)

func main() {
	fmt.Println("UDT lib called through manet")
	manetlib()
	fmt.Println("UDT lib called directly")
	udtlib()
}

func runtest(accept func() (net.Conn, error), dial func() (net.Conn, error)) {
	txsize := 4096
	nloops := 50000

	go func() {
		c, err := dial()
		if err != nil {
			panic(err)
		}

		defer c.Close()

		buf := make([]byte, txsize)
		for i := 0; i < nloops; i++ {
			n, err := c.Write(buf)
			if err != nil {
				panic(err)
			}
			if n != txsize {
				fmt.Println(n, txsize)
				panic("failed to write correct size")
			}
		}
	}()

	oc, err := accept()
	if err != nil {
		panic(err)
	}

	before := time.Now()
	n, err := io.Copy(ioutil.Discard, oc)
	if err != nil {
		panic(err)
	}
	if n != int64(txsize*nloops) {
		panic("not enough bytes")
	}
	took := time.Now().Sub(before)
	fmt.Printf("took %s\n", took)
	fmt.Printf("bandwidth = %f b/s\n", float64(txsize*nloops)/took.Seconds())
}

func udtlib() {
	list, err := udt.Listen("udt", "localhost:5555")
	if err != nil {
		panic(err)
	}
	defer list.Close()

	runtest(
		func() (net.Conn, error) {
			return list.Accept()
		},
		func() (net.Conn, error) {
			return udt.Dial("udt", "localhost:5555")
		},
	)
}

func manetlib() {
	addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/5556/udt")
	if err != nil {
		panic(err)
	}

	list, err := manet.Listen(addr)
	if err != nil {
		panic(err)
	}

	defer list.Close()

	runtest(
		func() (net.Conn, error) {
			return list.Accept()
		},
		func() (net.Conn, error) {
			return manet.Dial(addr)
		},
	)
}
