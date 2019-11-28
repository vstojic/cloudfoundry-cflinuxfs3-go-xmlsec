#!/bin/bash

cf plugins | grep ^cflocal >/dev/null
[[ $? -ne 0 ]]  && echo "Please, install the cflocal plugin by runinng: cf install-plugin cflocal" && exit 1

cf local stage xmldsig
