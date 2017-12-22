#!/bin/sh
go get -u github.com/rogpeppe/godef
go get -u github.com/gitbufenshuo/godefcache
cd $GOPATH/src/github.com/rogpeppe/godef
go install .
cd $GOPATH/src/github.com/gitbufenshuo/godefcache
go install .
goodname=godef_raw
mv $GOPATH/bin/godef $GOPATH/bin/${goodname}
godefcache -s ${goodname}
mv $GOPATH/bin/godefcache $GOPATH/bin/godef
