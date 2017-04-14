#!/bin/bash

filename=$@
base="${filename%.*}"
ext="${filename##*.}"

r0=$filename
r1=$base-1.$ext
r2=$base-2.$ext
r3=$base-3.$ext
r4=$base-4.$ext
r5=$base-5.$ext
r6=$base-6.$ext
r7=$base-7.$ext
r8=$base-8.$ext

tag() {
  local f=$1
  local t=$2
  echo $f
  exiftool -n -Orientation=$t -overwrite_original -q $f
  identify -verbose $f | grep Orientation
}

rm -rf $base-*.$ext

convert -strip $r0 $r1
tag $r1 1

convert $r1 -flop $r2
tag $r2 2

convert $r1 -rotate 180 $r3
tag $r3 3

convert $r1 -flip $r4
tag $r4 4

convert $r2 -rotate -90 $r5
tag $r5 5

convert $r1 -rotate -90 $r6
tag $r6 6

convert $r2 -rotate 90 $r7
tag $r7 7

convert $r1 -rotate 90 $r8
tag $r8 8
