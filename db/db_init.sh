#/usr/bin/env sh

set -euxo pipefail

if psql -lqt | cut -d \| -f 1 | grep -w "traintrackdb"; then
  >&2 printf "\e[31mtraintrackdb exists already\e[m\n"
  exit 1
else
  createdb "traintrackdb"
fi

psql -d traintrackdb -f ./db.sql
