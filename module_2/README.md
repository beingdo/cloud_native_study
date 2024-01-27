##制作镜像方法

进入到这个目录，执行
docker build -t http_server .

##用这个镜像启动容器

docker run -it -d -p 8080:8080 http_server

##测试

curl 127.0.0.1:8080
curl 127.0.0.1:8080/healthz
curl -I 127.0.0.1:8080/panic
