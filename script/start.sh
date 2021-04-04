#########################################################################
# File Name: script/start.sh
# Author: suixibing
# Desc: start IM service
# Created Time: Sun Apr  4 13:44:20 2021
#########################################################################
# !/bin/bash 

cd ${IM_PATH}
${IM_PATH}/${IM_NAME} &
cd -
