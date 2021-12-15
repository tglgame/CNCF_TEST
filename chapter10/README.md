# 文件说明： 
1. `server/httpserver.go`新增了基于信号的优雅终止，并且配置和代码分离，实现了日志存储
2. `limitrange.yaml`限制了namespace下单个pod和容器的资源使用
3. `quota.yaml`限制了namespace下能使用的资源总量
4. `deploy.yaml`会启动三副本，并且加入了：
   * `preStart`和`postStop`的验证
   * livenessProbe
   * startupProbe
   * `cm.sh`会新建一个configmap对象，它包含了`conf.json`和一个字符串变量的信息
5. `ingress.yaml`定义了ingress对象
6. `svcnodeport.yaml`定义了service对象
7. `gen_tls.sh`会生成key和crt文件，并创建secret对象


# 验证流程：  
1. K8S环境是在自己机器上搭建了三台虚拟机，ingress nginx使用的nodeport，并加入了externalIPs属性（[参考的这里](）https://kubernetes.github.io/ingress-nginx/deploy/baremetal/#over-a-nodeport-service))
2. 创建 namespace gracehttp
3. 创建limitrange 和 resourcequota对象
4. 创建ingress 对象  
5. 创建deploy对象


服务的运行信息：  
![image](https://gitee.com/tglgame/tools/raw/master/images/m8-serviceinfo.PNG)


通过浏览器访问：  
![image](https://gitee.com/tglgame/tools/raw/master/images/m8-explorercheck.PNG)


# 遗留问题：
在ingress里加入了secret之后，并不能通过https进行访问，还不清楚原因是什么
```
spec:
  tls:
  - secretName: httpsecret
```

# 说明： 
在模块8的基础上完善代码，每次请求 `/hello` 会有0-2秒的随机延迟，并上传时长数据到prometheus。

http server代码在 `server/httpserver.go`，相关代码片段如下：
```
	metrics.Register()

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", HelloServer)
	mux.HandleFunc("/healthz", Healthz)
	mux.Handle("/metrics", promhttp.Handler())
```

```
func HelloServer(w http.ResponseWriter, req *http.Request) {
	logger.Info("====receive client header====")

	// promethus监控
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	// 随机延迟0-2秒
	delay := randInt(0, 2)
	time.Sleep(time.Second * time.Duration(delay))

```

`metrics/metrics.go` 是老师提供的例子，在此基础上修改了Buckets的值：
```
		prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "execution_latency_seconds",
			Help:      help,
			Buckets:   prometheus.DefBuckets,
			// Buckets: prometheus.ExponentialBuckets(0.001, 2, 15),
		}, []string{"step"},
```
对prometheus里的指标使用还不熟练，这里只展示了下`httpserver_execution_latency_second_bucket`的数据

prometheus页面：
![image](https://gitee.com/tglgame/tools/raw/master/images/m10-prom.PNG)

grafana页面：
![image](https://gitee.com/tglgame/tools/raw/master/images/m10-grafana.PNG)

