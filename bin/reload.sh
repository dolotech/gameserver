#!/bin/bash

basedir=$(dirname $(readlink -f $0))
robot_name="${basedir}/gameserver"

count=`ps -ef | grep $robot_name | grep -v 'grep' | wc -l`

if [ 0 == $count ]; then
  echo "${robot_name} has not started !"
  exit 0
fi

pids=`ps -ef | grep $robot_name | grep -v 'grep' | awk '{print $2}'` 

echo 'reloading progress...'

for pid in $pids
do
  kill -1 $pid
  
  sleep 1
  if [ 0 == `ps -ef | grep $robot_name | grep -v 'grep' | grep $pid | wc -l` ]; then
    echo "reload ${pid} be failed"
  else
    echo "${pid} reloaded"
  fi
done

echo 'reload finish'

exit 0
