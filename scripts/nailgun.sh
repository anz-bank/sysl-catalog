#!/bin/sh

trap "kill 0" EXIT

for i in "$@"
do
case $i in
    -s|--serve)
        java -cp plantuml.jar:nailgun-server-1.0.0-SNAPSHOT.jar -server com.facebook.nailgun.NGServer&
        sysl-catalog $@ &
        wait
        echo "Done"
        exit 0
        
esac
done

sysl-catalog -o /out/ $@ &
wait
echo "Done"
exit 0


