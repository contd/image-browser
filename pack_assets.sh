#!/usr/bin/env bash

ansi --green "Removing old 'bindata*.go' files.."
rm bindata.go
#rm bindata_gzip.go

#ansi --cyan "Creating './templates' => 'bindata.go'"
#go-bindata ./templates/...

ansi --cyan "Creating './public' => 'bindata_gzip.go'"
bindata ./public/...

exit 0