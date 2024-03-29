#!/bin/bash

cat URLs.txt | parallel "ab -p post_json -T application/json -H 'Content-Type: application/json' -c 10 -n 100 {}"