#########################################################################
# File Name: build.sh
# Author: suixibing
# Desc: build package
# Created Time: Sun Apr  4 11:47:10 2021
#########################################################################
# !/bin/bash

# 引入常用变量
source script/comm.sh

# 编译
echo -e "${START_ECHO_FRONT}stage(1/2) 开始编译...${ECHO_BACK}"
go build -o ${IM_NAME}
echo -e "${SUCC_ECHO_FRONT}stage(1/2) 编译成功!${ECHO_BACK}\n"

# 打包
echo -e "${START_ECHO_FRONT}stage(2/2) 开始打包...${ECHO_BACK}"
# 在临时目录中进行打包
TMP_BUILD_DIR=.tmp_build_$(date +%s)
PACKAGE_FILE="${IM_NAME} ./conf ./view ./script ./asset"

mkdir -p ${TMP_BUILD_DIR}/${IM_NAME}
check_dir $(pwd)
cp -r ${PACKAGE_FILE} ${TMP_BUILD_DIR}/${IM_NAME}

cd ${TMP_BUILD_DIR}
tar czvf ${IM_PACKAGE} *
if [ "$?" != "0"  ] ; then
    echo -e "${ERR_ECHO_FRONT}压缩文件失败!${ECHO_BACK}"
    exit 1
fi
cd -

cp ${TMP_BUILD_DIR}/${IM_PACKAGE} .
rm -rf ${TMP_BUILD_DIR}
echo -e "${SUCC_ECHO_FRONT}stage(2/2) 生成包文件 ${IM_PACKAGE}，打包成功!${ECHO_BACK}"

