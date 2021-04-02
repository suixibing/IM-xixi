# !/bin/bash 

APP_PATH=/usr/local/app                                                                                                    
IM_NAME=IM-xixi
IM_PATH=${APP_PATH}/${IM_NAME}

function help()
{
  echo "Usage: ${IM_NAME} [option]"
  echo ""
  echo "option: help(default) | run | stop"
}

if [ $# != 1 ] ; then 
  help
  exit
fi

if [ $1 = "run" ] ; then
  ${IM_PATH}/script/start.sh 
elif [ "$1" = "stop" ] ; then
  ${IM_PATH}/script/stop.sh 
else
  help
fi
