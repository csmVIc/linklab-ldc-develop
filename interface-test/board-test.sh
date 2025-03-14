#!/bin/bash

while true; do
  curl -v \
  --header "Authorization: f41f6cd376219babade8ecbc0d77168d65215c65cfdf756843b1f72b66b7d17d" \
  http://10.214.149.214:30335/api/board/list
done