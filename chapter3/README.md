要求：  

1. 构建本地镜像  
2. 编写 Dockerfile 将练习 2.2 编写的 httpserver 容器化  
   * 请思考有哪些最佳实践可以引入到 Dockerfile 中来  
3. 将镜像推送至 docker 官方镜像仓库  
4. 通过 docker 命令本地启动 httpserver  
5. 通过 nsenter 进入容器查看 IP 配置  

<br>

**本地执行记录**

---

本地镜像和容器：
```
root@ubuntu:/home/dockerfiles/httpserver# docker images | grep http
tglgame/httpserver                                                            v1                  55da7c2873a2        6 days ago          11.7MB
root@ubuntu:/home/dockerfiles/httpserver# docker ps | grep -i http
b37657c49821        tglgame/httpserver:v1                                                      "./server"               6 days ago          Up 24 minutes       0.0.0.0:8003->80/tcp                       hardcore_lichterman
```

nsenter进入容器查看IP配置:
```
root@ubuntu:/home/dockerfiles/httpserver# docker inspect b37657c49821 | grep -i pid
            "Pid": 56244,
            "PidMode": "",
            "PidsLimit": null,
root@ubuntu:/home/dockerfiles/httpserver# nsenter -t 56244 -n ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
32: eth0@if33: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:0f brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.15/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
root@ubuntu:/home/dockerfiles/httpserver#
```


codehub上的镜像信息：  
![image](https://gitee.com/tglgame/tools/raw/master/images/dockerhub-httpserver.PNG)

client访问示例：
```
root@ubuntu:/home/dockerfiles/httpserver# docker ps -a | grep -i http
b37657c49821        tglgame/httpserver:v1                                                      "./server"               6 days ago           Up 15 seconds                   0.0.0.0:8003->80/tcp                       hardcore_lichterman
root@ubuntu:/home/dockerfiles/httpserver# 
root@ubuntu:/home/dockerfiles/httpserver# 
root@ubuntu:/home/dockerfiles/httpserver# docker logs b37657c49821
http_server:2021/10/10 13:18:32 httpserver.go:40: ====receive client header====
http_server:2021/10/10 13:18:32 httpserver.go:43: client addr:  192.168.87.1
http_server:2021/10/10 13:18:32 httpserver.go:46: receive and set: User-Agent : Go-http-client/1.1
http_server:2021/10/10 13:18:32 httpserver.go:46: receive and set: Ctag1 : value1
http_server:2021/10/10 13:18:32 httpserver.go:46: receive and set: Ctag2 : value2
http_server:2021/10/10 13:18:32 httpserver.go:46: receive and set: Accept-Encoding : gzip
http_server:2021/10/10 13:18:32 httpserver.go:57: return status code:  200
http_server:2021/10/10 13:18:32 httpserver.go:60: =====================
http_server:2021/10/10 13:18:32 httpserver.go:66: get health of server
http_server:2021/10/16 14:40:17 httpserver.go:40: ====receive client header====
http_server:2021/10/16 14:40:17 httpserver.go:43: client addr:  192.168.87.1
http_server:2021/10/16 14:40:17 httpserver.go:46: receive and set: User-Agent : Go-http-client/1.1
http_server:2021/10/16 14:40:17 httpserver.go:46: receive and set: Ctag1 : value1
http_server:2021/10/16 14:40:17 httpserver.go:46: receive and set: Ctag2 : value2
http_server:2021/10/16 14:40:17 httpserver.go:46: receive and set: Accept-Encoding : gzip
http_server:2021/10/16 14:40:17 httpserver.go:57: return status code:  200
http_server:2021/10/16 14:40:17 httpserver.go:60: =====================
http_server:2021/10/16 14:40:17 httpserver.go:66: get health of server
```
