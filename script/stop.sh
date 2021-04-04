#########################################################################
# File Name: script/stop.sh
# Author: shiawaseli
# Desc: stop IM service
# Created Time: Sun Apr  4 13:52:25 2021
#########################################################################
# !/bin/bash 

# 获取所有的IM进程的pid
PID=$(ps aux | grep ${IM_NAME} | grep -v grep | awk '{print $2}')

if [ "${PID}" != "" ] ; then
    kill -9 ${PID}
fi
