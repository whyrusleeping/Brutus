#!/bin/sh

#apt-get install golang-go
#apt-get install git

git clone https://github.com/wernerd/Skein3Fish.git
cd Skein3Fish/go
export GOPATH=`pwd`
go install crypto/threefish
go install crypto/skein

go test crypto/threefish
go test crypto/skein

cd ../../
git clone https://github.com/TLane/Brutus.git
cd Brutus
make
./xkcd
