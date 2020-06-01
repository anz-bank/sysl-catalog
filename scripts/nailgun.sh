#!/bin/bash

trap "exit" INT TERM ERR
trap "kill 0" EXIT

java -cp plantuml.jar:nailgun-server-1.0.0-SNAPSHOT.jar -server com.facebook.nailgun.NGServer &

/usr/sysl-catalog -o /out/ $@

wait