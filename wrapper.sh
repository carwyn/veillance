#!/bin/bash

LOGFILE=/var/tmp/veillance.log

while true ; do

  echo "Interceptor started: `date`" | tee -a ${LOGFILE}
  interceptor -i wlp1s0 "port 80 or port 53"
  echo "Interceptor exited : `date`" | tee -a ${LOGFILE}
  sleep 3

done

