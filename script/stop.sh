# !/bin/bash 

APP_PATH=/usr/local/app                                                                                                    
IM_NAME=IM-xixi
IM_PATH=${APP_PATH}/${IM_NAME}

pid=`ps aux | grep ${IM_NAME} | grep -v grep | awk '{print $2}'`

if [ "${pid}" != "" ] ; then
  kill -9 ${pid}
fi
