# !/bin/bash 

APP_PATH=/usr/local/app                                                                                                    
IM_NAME=IM-xixi
IM_PATH=${APP_PATH}/${IM_NAME}

cd ${IM_PATH}
${IM_PATH}/${IM_NAME} &
cd -
