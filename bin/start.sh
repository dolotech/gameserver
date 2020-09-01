#!/bin/bash

basedir=$(dirname $(readlink -f $0))
robot_name="${basedir}/gameserver"


if [ ! -f "${robot_name}" ] 
then
  echo "${robot_name} file is not exist!"
  exit 0
fi

chmod 755 $robot_name

count=`ps -ef | grep $robot_name | grep -v 'grep' | wc -l`

if [ $count != 0 ] 
then 
  echo "gameserver already started !"
  exit 0
fi

echo 'start gameserver...'

sleep 1

nohup $robot_name  > /dev/null 2>&1 &
echo 'start success'
count=`ps -ef | grep $robot_name | grep -v 'grep' | wc -l`

if [ $count != 0 ] 
then 
  echo "gameserver started !"
  exit 0
else
	echo "run fail!run fail!!!"
fi

exit 0
