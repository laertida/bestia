##
# la bestia
#
# @file
# @version 0.1

build:
	GOOS=linux GOARCH=arm GOARM=6 go build

transfer: build
	scp bestia canon:/home/laertida/

exec:
	ssh -t laertida@canon ./bestia

run: build transfer exec


# end
