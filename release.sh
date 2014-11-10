#!/bin/bash

set -e

rm -f dist/*

gox -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"

cd dist

for file in `ls`; do
    if [[ $file == *.exe ]]
    then
        mv $file pup.exe
        file=${file%.exe}
        zip $file pup.exe
        rm pup.exe 
    else
        mv $file pup
        zip $file pup
        rm pup 
    fi
done
