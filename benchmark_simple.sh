#!/bin/bash

ab -p post_json -T application/json -H 'Content-Type: application/json' -c 100 -n 1000 http://localhost:8080/THE_TOKEN