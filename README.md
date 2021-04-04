# IM-xixi

一个简单的Web即时通信软件（Instant Messaging，IM）



#	一、使用源代码部署

##	1、首先拉取代码并进入目录

```bash
$ git clone https://github.com/shiawaseli/IM-xixi.git
$ cd IM-xixi
```

##	2、使用 build.sh 编译并打包编译结果

```bash
$ ./build.sh
```

##	3、使用 install.sh 安装程序

```bash
# 编译前请先准备好如下目录并保证执行者对其有写权限。
# 注意环境变量 export PATH=$PATH:/usr/local/app/bin
# /usr/local/app
# /usr/local/app/bin
# /data
$ ./install.sh
```

##	4、使用mysql数据库

**注意：**默认使用如下配置连接数据库，运行前请先准备好相关的数据库（程序会自动建表），或者更改默认的配置。

**root:123456@(127.0.0.1:3306)/chat_xixi?charset=utf8**

##	5、启动程序

**注意：**默认使用80端口进行启动，请保证端口不被占用且有相应的权限。

```bash
$ IM-xixi run
```

##	6、关闭程序

**注意：**会关闭所有的同名程序。

```bash
$ IM-xixi stop
```



#	二、使用Docker快速部署

##	1、拉取镜像

latest中仅支持最低限度运行的环境，要想可以支持编译，可以选择dev_base版本的镜像。

```bash
$ docker pull shiawaseli/chat_xixi:latest
```

##	2、使用镜像创建容器

```bash
$ docker run -itd --name chat_xixi -v /data/IM-xixi/mnt:/data/IM-xixi/mnt --network host shiawaseli/chat_xixi:latest
```

**参数说明:**

|            --name            |          容器的名称          |
| :--------------------------: | :--------------------------: |
| -v [宿主机地址]:[容器中地址] | 将宿主机上的目录映射到容器中 |
|          --network           |       共享主机的网络ip       |

可以根据需要自己配置合适的参数

##	3、进入容器中启动服务

```bash
$ docker exec -it chat_xixi /bin/bash #进入容器
$ IM-xixi run #启动服务
$ exit #退出容器
```



#		三、其他

##	1、可以使用docker部署mysql

```bash
# 这是一个最简单可以启动本服务的mysql镜像(version 5.7)，仅有创建有最基本的database
# 并且数据集使用了utf8，可以支持显示中文不乱码
$ docker pull shiawaseli/chat_xixi:mysql
# 将宿主机的3306端口映射到容器中，root密码为123456
$ docker run -itd -p 3306:3306 -e MYSQL\_ROOT\_PASSWORD=123456 --name chat_mysql shiawaseli/chat_xixi:mysql
# 使用如下命令可以进入mysql容器
$ docker exec -it chat_mysql /bin/bash
# 也可以直接在宿主机中直接登录mysql
$ mysql -h'127.0.0.1' -P'3306' -u'root' -p'123456'
```

##	2、从宿主机中登陆mysql发现中文出现乱码

此时使用的是宿主机的client登陆mysql，所以数据在本地出现了乱码，可以使用如下命令检查此时的字符集。

```mysql
mysql> show variables like '%char%';
```

如果出现**latin1**，可以通过在**my.cnf**中添加如下配置解决：

```ini
[mysqld]
character-set-server=utf8 

[client]
default-character-set=utf8 

[mysql]
default-character-set=utf8
```

##	3、关于上传的文件

上传的文件会保存到/data/${IM_NAME}/mnt下，服务会创建**软链接**到该目录，请保证运行用户有该目录的写权限。

##	4、关于修改IM变量

修改comm.sh中的变量时别忘了一起修改control.sh中的变量。

