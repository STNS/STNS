#!/bin/bash

if [[ $1 == "test" ]]; then
  echo "ok"
  exit 0
fi
exit 1
