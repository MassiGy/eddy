#! /bin/bash

tag=$(git tag | tail -1)
if [ ! -z "$tag" -a  "$tag" != "" ]; then
    echo $(git tag | tail -1) > ./VERSION
fi