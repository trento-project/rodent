#!/bin/sh
OUT=$( git describe --tags --always --dirty | grep -o '.....$' )

if [[ ${OUT} == 'dirty' ]] ; then
  echo 'dirty'
fi
