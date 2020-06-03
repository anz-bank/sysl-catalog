#!/bin/sh

trap "kill 0" EXIT

java -cp plantuml.jar:nailgun-server-1.0.0-SNAPSHOT.jar -server com.facebook.nailgun.NGServer>/dev/null 2>/tmp/stderr &
server=0

for i in "$@"
do
case $i in
    -s|--serve)
        ./sysl-catalog -o /out/ $@ &
        server=1
        wait
        echo "Done"
        
esac
done

if [ $server -eq 0 ]
then
./sysl-catalog -o /out/ $@
fi

