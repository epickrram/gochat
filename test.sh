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

for i in `find . -iname '*_test.go'`;
do
	package_dir=`dirname $i`
	if [[ $package_dir =~ .*epickrram.* ]]
	then
		echo go test $FLAGS `dirname $i`
		go test $FLAGS `dirname $i`
	fi
done
