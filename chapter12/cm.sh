#!/bin/bash

kubectl create cm hscm --from-file=conf.json --from-literal=configpath=/home/config-file-path/conf.json
