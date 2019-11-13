all:nftmain

nftmain:
	go build -ldflags "-s" -o $@  $@.go 
	#strip $@
	mv $@ $(HOME)/bin
clean:
	go clean -x
