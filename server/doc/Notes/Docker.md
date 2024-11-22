# docker file语法
## 1. FROM
-指定基础镜像（必须有的指令，并且必须是第一条指令）
`FROM golang:1.21-alpine`

## 2. WORKDIR
使用 WORKDIR 指令可以来指定工作目录（或者称为当前目录）
`WORKDIR /app`

## 3. COPY
COPY 指令将从构建上下文目录中 <源路径> 的文件/目录复制到新的一层的镜像内的 <目标路径> 位置
`COPY go.mod go.sum ./`

## 4. RUN
用于执行命令行命令
`RUN go mod download`


# 5. EXPOSE
 EXPOSE 指令是声明运行时容器提供服务端口
`EXPOSE 8080`

# 6. ENTRYPOINT
 ENTRYPOINT 指令是指定容器启动时执行的程序及参数
`ENTRYPOINT ["./main"]`


# docker 部署nginx并进行配置
##　１．打包静态资源，并COPY到容器nginx中：
```dockerfile
WORKDIR /app/front
COPY front .
RUN pnpm install && pnpm build

WORKDIR /app/admin
COPY admin .
RUN pnpm install && pnpm build

```
## 2.静态资源部署在nginx中
```dockerfile
FROM nginx:1.24.0-alpine

COPY --from=BUILD /app/front/dist /usr/share/nginx/html/front
COPY --from=BUILD /app/admin/dist /usr/share/nginx/html/admin
```
## 3.准备nginx配置文件
```template
server {
    listen 80;
    server_name ${SERVER_NAME};
    root /usr/share/nginx/html;

    ######### 反向代理 start ########
    location /api {
        proxy_pass http://${BACKEND_HOST}:${SERVER_PORT};
        # 获取真实 IP 
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header X-Forwarded-Proto $remote_addr;

        client_max_body_size 40M;
    }
    # Gin 处理静态资源, Nginx 反向代理 Gin (本地文件上传才需要)
    location /public/uploaded {
        proxy_pass http://${BACKEND_HOST}:${SERVER_PORT}/public/uploaded;
    }
    location /admin/public/uploaded {
        proxy_pass http://${BACKEND_HOST}:${SERVER_PORT}/public/uploaded;
    }
    ######### 反向代理 end ########

    ######### 静态资源 start ########
    location / {
        try_files $uri $uri/ /index.html;
    }
    location /admin {
        try_files $uri $uri/ /admin/index.html;
    }
    location ~ .*\.(js|json|css)$ {
        gzip on;
        gzip_static on; 
        gzip_min_length 1k;
        gzip_http_version 1.1;
        gzip_comp_level 9;
        gzip_types text/css application/javascript application/json;
        root /usr/share/nginx/html;
    }
    ######### 静态资源 end ########
}
```
##　把配置文件拷贝在容器中，运行初始化脚本并运行Nginx
```dockerfile
COPY deploy/build/web/default.conf.template /etc/nginx/conf.d/default.conf.template
COPY deploy/build/web/default.conf.ssl.template /etc/nginx/conf.d/default.conf.ssl.template
COPY deploy/build/web/run.sh /docker-entrypoint.sh
RUN chmod a+x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]
## docker中喜欢把需要初始化的脚本重命名为　/docker-entrypoint.sh 
## 然后entrypoint 指向这个脚本　完成初始化
CMD [ "nginx", "-g", "daemon off;" ]

EXPOSE 80
EXPOSE 443
```


# docker 部署mysql
---
## 定义初始化文件并执行
`mysql镜像里面　/docker-entrypoint-initdb.d/` 这个目录下的脚本会自动执行,存放初始脚本
```dockerfile
FROM mysql:8.0

ENV WORK_PATH  /usr/local/work
ENV AUTO_RUN_DIR /docker-entrypoint-initdb.d
ENV RUN_SHELL run.sh

COPY ./mysql.cnf /etc/mysql/conf.d
COPY ./gvb.sql ${WORK_PATH}
COPY ./${RUN_SHELL} ${AUTO_RUN_DIR}

RUN chmod a+x ${AUTO_RUN_DIR}/${RUN_SHELL}

```
```bash
#!/bin/bash

mysql -uroot -p$MYSQL_ROOT_PASSWORD << EOF
system echo '================Start create database gvb====================';
source $WORK_PATH/gvb.sql
system echo '================OK!====================';
EOF
```
- 到这里　基本的mysql容器就配置完成了.具体怎么使用？请看下节　docker compose　组合各个容器


# docker compose 组合容器
---
docker 容器单个常常不满足我们的使用,因此我们需要组合多个容器来完成我们的项目 如mysql等 因此docker compose出现
## 使用方法
首先 docker compose 服务由docker-compose.yml 这个文件来管理

`docker-compose.yml 同目录下的 .env 文件会被加载为其环境变量`
因此我们现将其环境变量写入到.env 文件中

## 容器间网络通讯
```yaml
networks:  // 定义网络   
  gvb-network:
    driver: bridge   //网络驱动默认选择桥接网络
    ipam:   //ip管理 通常定义子网
      config:
        - subnet: 172.12.0.0/24
```
- 在docker-compose.yml定义了networks 网络 然后在服务中使用就可以完成通讯:

```yaml
services:
  gvb-redis:
    networks:
      gvb-network:
        ipv4_address: ${REDIS_HOST}

  gvb-mysql:
    networks:
      gvb-network:
        ipv4_address: ${MYSQL_HOST}

  gvb-server:
    networks:
      gvb-network:
        ipv4_address: ${BACKEND_HOST}

  gvb-web:
    networks:
      gvb-network:
        ipv4_address: ${FRONTEND_HOST}
```
## 数据库挂载数据卷
倘若不做任何干涉，在每次启动一个 Mysql或其他数据库 容器时，数据库都是空的。另外容器删除之后，数据就丢失了（还有各类意外情况），非常糟糕！

数据卷 是被设计用来持久化数据的，它的生命周期独立于容器，Docker 不会在容器被删除后自动删除 数据卷，并且也不存在垃圾回收这样的机制来处理没有任何容器引用的 数据卷。
- 多个容器可以共享同一个数据卷，从而实现数据共享
- 通过将数据卷挂载到主机上，你可以方便地进行数据备份和恢复

```yaml
volumes:
  - ${DATA_DIRECTORY}/data/redis/:/data  ##把主机的${DATA_DIRECTORY}/data/redis/目录 挂载到 容器的/data目录
```
## 容器间的依赖关系
```yaml
  gvb-mysql:
    build: ../build/mysql
    container_name: gvb-mysql
    restart: always
    volumes:
      - ${DATA_DIRECTORY}/data/mysql/:/var/lib/mysql
    environment:
      - MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - TZ=Asia/Shanghai
    ports:
      - ${MYSQL_PORT}:3306
    networks:
      gvb-network:
        ipv4_address: ${MYSQL_HOST}
  
  gvb-server:
    build: ../../server
    container_name: gvb-server
    restart: always
    depends_on: 
      gvb-redis:
        condition: service_started
      gvb-mysql:
        condition: service_started
    volumes:
      - ${DATA_DIRECTORY}/file/uploaded:/gvb/public/uploaded
    environment:
    ## 环境变量的配置会覆盖后端中的 config.docker.yml
      - TZ=Asia/Shanghai
      - SERVER_PORT=:${SERVER_PORT}
      - MYSQL_HOST=${MYSQL_HOST}
      - MYSQL_PORT=3306
      - MYSQL_DBNAME=gvb
      - MYSQL_USERNAME=root
      - MYSQL_PASSWORD=${MYSQL_ROOT_PASSWORD}
      - REDIS_ADDR=${REDIS_HOST}:6379
      - REDIS_PASSWORD=${REDIS_PASSWORD}
    ports:
      - ${SERVER_PORT}:${SERVER_PORT}
    networks:
      gvb-network:
        ipv4_address: ${BACKEND_HOST}
```

在例子中 server这个后端服务依赖于  mysql 和 redis 两个服务,因此在server 服务中 加入了 `depends_on` 依赖. 并且加入了 `condition: service_started` 条件,表示只有当依赖的服务启动后，才会启动当前服务。

此时多个容器就可以实现合作运行一个完整的项目.

# 遇到的问题
## go mod下载超时
这是我的dockerfile指令,但是在下载mod的时候总是超时 显示timeout
```dockerfile
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download
```
**尝试换了很多次代理也没有用**

**解决方法**:
`docker build --network=host -t`
使用本机的网络栈来运行容器即可解决. 因为多次切换代理问题仍然存在,即问题处在容器内部的网络设置.

