#!/bin/bash

k create cm hscm -n gracehttp --from-file=conf.json --from-literal=configpath=/home/config-file-path/conf.json
