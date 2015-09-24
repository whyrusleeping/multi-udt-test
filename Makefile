LIBPATH=$(GOPATH)/src/github.com/jbenet/go-multiaddr-net/vendor/go-udtwrapper-v1.0.0
all:
	cd $(LIBPATH) && make go-deps
	cp $(LIBPATH)/udt/libudt.a .
	go build
