#!/bin/bash
PLANTUML_VERSION=1.2019.10
wget "http://downloads.sourceforge.net/project/plantuml/${PLANTUML_VERSION}/plantuml.${PLANTUML_VERSION}.jar" -O plantuml.jar
curl -O https://github.com/facebook/nailgun/releases/download/nailgun-all-v1.0.0/nailgun-server-1.0.0-SNAPSHOT.jar
