#!/bin/bash

for f in $@; do
  echo $f
  identify -verbose $f | grep Orientation
done
