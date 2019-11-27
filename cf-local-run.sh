#!/bin/bash

cf local run xmldsig -i $(docker-machine ip) -p 8080 &
CF_LOCAL_RUN_PID=$!
export CF_LOCAL_RUN_PID
echo
echo "The process (PID=$CF_LOCAL_RUN_PID) is running."
#echo "Use kill -9 \$CF_LOCAL_RUN_PID      to get rid of it."
echo "Use docker kill \$(docker ps -alq)  to get rid of it."
echo "Use curl \$(docker-machine ip):8080 to get the response from it."
echo


