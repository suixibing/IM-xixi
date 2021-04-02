# !/bin/bash 

APP_PATH=/usr/local/app
IM_NAME=IM-xixi
IM_PATH=${APP_PATH}/${IM_NAME}

echo "${IM_NAME} building..."
mkdir -p ./bin
go build -o ./bin/${IM_NAME}

mkdir -p ${IM_PATH}/mnt
cp -r ./view ${IM_PATH}
cp -r ./conf ${IM_PATH}
cp -r ./script ${IM_PATH}
cp -r ./asset ${IM_PATH}
cp ./bin/${IM_NAME} ${IM_PATH}

# 创建软链接
rm -rf ${APP_PATH}/bin/${IM_NAME}
ln -s ${IM_PATH}/script/control.sh ${APP_PATH}/bin/${IM_NAME}
