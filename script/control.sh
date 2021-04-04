#########################################################################
# File Name: script/control.sh
# Author: shiawaseli
# Desc: control IM service
# Created Time: Sun Apr  4 13:27:07 2021
#########################################################################
# !/bin/bash 

APP_PATH=/usr/local/app
IM_NAME=IM-xixi
IM_PATH=${APP_PATH}/${IM_NAME}

function help()
{
    echo "Usage: $0 [option]"
    echo ""
    echo "option: help(default) | run | stop"
}

if [ $# != 1 ] ; then 
    help
    exit 1
fi

if [ $1 = "run" ] ; then
    source ${IM_PATH}/script/start.sh 
elif [ "$1" = "stop" ] ; then
    source ${IM_PATH}/script/stop.sh 
else
    help
fi
