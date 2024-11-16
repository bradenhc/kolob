#!/bin/sh
here=$(dirname $0)

plantuml_jar='plantuml-asl-1.2024.7.jar'

plantuml_install="/opt/plantuml/$plantuml_jar"

if [ ! -f "$plantuml_install" ]
then
    echo "ERR: missing plantuml: install JAR to $plantuml_install"
    exit 1
fi

cd $here
mkdir -p img
java -jar "$plantuml_install" -o 'img' '*.plantuml.txt'
