package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"time"

	manet "github.com/jbenet/go-multiaddr-net"
	ma "github.com/jbenet/go-multiaddr-net/Godeps/_workspace/src/github.com/jbenet/go-multiaddr"
	udt "github.com/jbenet/go-multiaddr-net/vendor/go-udtwrapper-v1.0.0/udt"
)

func main() {
	fmt.Println("UDT lib called through manet")
	maina()
	fmt.Println("UDT lib called directly")
	mainb()
}

func mainb() {
	list, err := udt.Listen("udt", "localhost:5555")
	if err != nil {
		panic(err)
	}
	defer list.Close()

	txsize := 4096
	nloops := 50000

	go func() {
		c, err := udt.Dial("udt", "localhost:5555")
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

	oc, err := list.Accept()
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

func maina() {
	addr, err := ma.NewMultiaddr("/ip4/127.0.0.1/udp/5556/udt")
	if err != nil {
		panic(err)
	}

	list, err := manet.Listen(addr)
	if err != nil {
		panic(err)
	}
	defer list.Close()

	txsize := 4096
	nloops := 50000

	go func() {
		c, err := manet.Dial(addr)
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

	oc, err := list.Accept()
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
