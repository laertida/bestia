##
# la bestia
#
# @file
# @version 0.1

build:
	GOOS=linux GOARCH=arm GOARM=6 go build

transfer:
	scp bestia ratito:/home/laertida/
exec:
	ssh -t laertida@ratito ./bestia

run: build transfer exec


# end
