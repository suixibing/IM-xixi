# !/bin/bash 

set -v

APP_PATH=/usr/local/app
IM_NAME=IM-xixi
IM_PATH=${APP_PATH}/${IM_NAME}

# 删除原有的文件
OLD_APP_BIN=`whereis IM-xixi | awk '{print $2}'`
rm -rf ${IM_PATH} ${OLD_APP_BIN} bin

echo "${IM_NAME} building..."
mkdir -p bin ${IM_PATH}
go build -o bin/${IM_NAME}

cp -r ./view ${IM_PATH}
cp -r ./conf ${IM_PATH}
cp -r ./script ${IM_PATH}
cp -r ./asset ${IM_PATH}
cp ./bin/${IM_NAME} ${IM_PATH}

# 创建软链接
rm -rf ${APP_PATH}/bin/${IM_NAME}
ln -s ${IM_PATH}/script/control.sh ${APP_PATH}/bin/${IM_NAME}

IM_MNT_PATH=/data/${IM_NAME}/mnt

# 创建mnt目录的链接
mkdir -p ${IM_MNT_PATH}
ln -s ${IM_MNT_PATH}/ ${IM_PATH}/mnt
