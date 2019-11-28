#!/bin/bash

cf plugins | grep ^cflocal >/dev/null
[[ $? -ne 0 ]]  && echo "Please, install the cflocal plugin by runinng: cf install-plugin cflocal" && exit 1

cf local run xmldsig -i $(docker-machine ip) -p 8080 &
CF_LOCAL_RUN_PID=$!
export CF_LOCAL_RUN_PID
#echo "The process (PID=$CF_LOCAL_RUN_PID) is running."
#echo "Use kill -9 \$CF_LOCAL_RUN_PID      to get rid of it."

sleep 5

echo
docker ps -al
echo
echo "Use docker kill \$(docker ps -alq)  to get rid of it."
echo "Use curl \$(docker-machine ip):8080 to get the response from it."
echo


