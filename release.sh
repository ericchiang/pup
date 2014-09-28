#!/bin/bash

gox -output "dist/{{.Dir}}_{{.OS}}_{{.Arch}}"
