#!/bin/bash

set -e

. ./env.sh
. ./clean.sh

FLAGS="$*"

if [ "-v" == "$1" ];
then
	FLAGS="-v"
fi

echo "flags: $FLAGS"

for i in `find $GOPATH/src/epickrram.com -iname '*_test.go'`;
do
	echo go test $FLAGS `dirname $i`
	go test $FLAGS `dirname $i`
done
