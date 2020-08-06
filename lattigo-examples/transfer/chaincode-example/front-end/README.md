#### 文档说明

* 在dist同一级下创建并编辑Dockerfile文件

```shell
FROM nginx
COPY /dist /usr/share/nginx/html
```

* 编译镜像

```shell
docker build -t digit-coin .
```

* 在dist同一级下创建并编辑nginx.conf文件

```shell
user  nginx;
worker_processes  1;

error_log  /var/log/nginx/error.log warn;
pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                      '$status $body_bytes_sent "$http_referer" '
                      '"$http_user_agent" "$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;

    keepalive_timeout  65;

    # include /etc/nginx/conf.d/*.conf;

    server {
        listen       8080;
        charset utf-8;
        server_name  localhost;# 服务器地址或绑定域名

        location / {
           root   /usr/share/nginx/html;
           index  index.html index.htm;
           try_files $uri $uri/ /index.html;
        }

        location @router {
          rewrite ^.*$ /index.html last;
        }

        location /v1/ {
           proxy_pass http://本地ip地址:8080;
        }

        error_page   500 502 503 504  /50x.html;
        location = /50x.html {
            root   /usr/share/nginx/html;
        }
   }
}
```

* 运行镜像

```shell
docker run --name digit-coin -p "8080:8080" -v $PWD/log:/var/log/nginx -v $PWD/nginx.conf:/etc/nginx/nginx.conf digit-coin:latest
```

