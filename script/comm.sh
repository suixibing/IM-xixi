#########################################################################
# File Name: script/comm.sh
# Author: shiawaseli
# Desc: common values
# Created Time: Sun Apr  4 11:32:42 2021
#########################################################################
# !/bin/bash

APP_PATH=/usr/local/app
IM_NAME=IM-xixi
IM_PATH=${APP_PATH}/${IM_NAME}
IM_PACKAGE=${IM_NAME}.tar.gz

START_ECHO_FRONT="\033[37m"
CHECK_ECHO_FRONT="\033[33m"
SUCC_ECHO_FRONT="\033[32m"
ERR_ECHO_FRONT="\033[31m"
ECHO_BACK="\033[0m"

function check_dir() {
    for DIR in $* ; do
        if [ ! -d "${DIR}" ] ; then
            echo -e "${ERR_ECHO_FRONT}${DIR} 目录不存在!${ECHO_BACK}"
            exit 1
        fi
        TMP_FILE=${DIR}/touch_test_$(date +%s)
        touch ${TMP_FILE}
        if [ "$?" != "0" ] ; then
            echo -e "${ERR_ECHO_FRONT}${DIR} 目录没有写权限!${ECHO_BACK}"
            exit 1
        fi
        rm -rf ${TMP_FILE}
    done
}

function check_ret() {
    if [ "$?" != "0" ] ; then
        echo -e "${ERR_ECHO_FRONT}$1 运行失败!${ECHO_BACK}"
        exit 1
    fi
}
