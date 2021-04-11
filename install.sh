#########################################################################
# File Name: install.sh
# Author: suixibing
# Desc: install package
# Created Time: Sun Apr  4 11:57:19 2021
#########################################################################
# !/bin/bash

# 引入常用变量
source script/comm.sh

IM_PACKAGE=$1

# 没传参数时尝试使用默认的文件
if [ "$#" == "0" ] ; then 
    IM_PACKAGE=${IM_NAME}.tar.gz
fi

# stage1 确认包文件存在
echo -e "${START_ECHO_FRONT}stage(1/4) 确认 ${IM_PACKAGE} 包文件存在...${ECHO_BACK}"
if [ ! -f "${IM_PACKAGE}" ] ; then
    echo -e "${ERR_ECHO_FRONT}${IM_PACKAGE} 包文件不存在! 请传入正确的包文件...${ECHO_BACK}"
    echo "Usage: $0 package_path"
    echo ""
    echo -e "必须是${ERR_ECHO_FRONT}tar.gz${ECHO_BACK}文件，例如: ${IM_NAME}.tar.gz"
    exit 1
fi
echo -e "${SUCC_ECHO_FRONT}stage(1/4) 包文件存在!${ECHO_BACK}"

# stage2 确认必要的目录存在
echo -e "${START_ECHO_FRONT}stage(2/4) 确认必要的目录存在并具有写权限...${ECHO_BACK}"
check_dir ${APP_PATH} ${APP_PATH}/bin /data
# 检查环境变量PATH中是否包含${APP_PATH}/bin
PATH_CHECK=$(env | grep "^PATH=" | grep "${APP_PATH}/bin" | wc -l | awk '{print $1}')
if [ "${PATH_CHECK}" != "1" ] ; then
    echo 'export PATH=$PATH:'${APP_PATH}/bin >> ~/.bashrc
    echo -e "${CHECK_ECHO_FRONT}请运行 source ~/.bashrc 更新环境变量!${ECHO_BACK}"
fi
echo -e "${SUCC_ECHO_FRONT}stage(2/4) 目录符合条件!${ECHO_BACK}"

# stage3 清理原有的文件并解压缩包到${APP_PATH}
echo -e "${START_ECHO_FRONT}stage(3/4) 清理原有的文件并解压缩包到 ${APP_PATH}...${ECHO_BACK}"
OLD_APP_BIN=$(whereis ${IM_NAME} | awk '{print $2}')
rm -rf ${IM_PATH} ${OLD_APP_BIN} ${APP_PATH}/bin/${IM_NAME}
if [ "$?" != "0" ] ; then
    echo -e "${ERR_ECHO_FRONT}清理原有的文件失败!${ECHO_BACK}"
    exit 1
fi

# 解压缩包到${APP_PATH}
tar zxvf ${IM_PACKAGE} -C ${APP_PATH}
if [ "$?" != "0" ] ; then
    echo -e "${ERR_ECHO_FRONT}解压缩包失败! 请检查压缩包是否出现损坏${ECHO_BACK}"
    exit 1
fi
echo -e "${SUCC_ECHO_FRONT}stage(3/4) 解压缩包成功!${ECHO_BACK}"

# stage4 创建服务和mnt目录的软链接
echo -e "${START_ECHO_FRONT}stage(4/4) 创建服务和mnt目录的软链接...${ECHO_BACK}"
ln -s ${IM_PATH}/script/control.sh ${APP_PATH}/bin/${IM_NAME}
SYMBOLIC_CHECK=$(ls -l ${APP_PATH}/bin | grep "${IM_NAME} -> ${IM_PATH}/script/control.sh" | wc -l)
if [ "${SYMBOLIC_CHECK}" == "0" ] ; then
    echo -e "${ERR_ECHO_FRONT}服务的软链接创建失败!${ECHO_BACK}"
    exit 1
fi

# 创建mnt目录的链接
IM_MNT_PATH=/data/${IM_NAME}/mnt
mkdir -p ${IM_MNT_PATH}
ln -s ${IM_MNT_PATH}/ ${IM_PATH}/mnt
SYMBOLIC_CHECK=$(ls -l ${IM_PATH} | grep "mnt -> ${IM_NNT_PATH}" | wc -l)
if [ "${SYMBOLIC_CHECK}" == "0" ] ; then
    echo -e "${ERR_ECHO_FRONT}mnt目录的软链接创建失败!${ECHO_BACK}"
    exit 1
fi

echo -e "${SUCC_ECHO_FRONT}stage(4/4) 软链接创建成功!${ECHO_BACK}"

