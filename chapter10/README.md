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

prometheus的配置里添加了如下job, targets里是http server nodeport service 里的 nodeip:nodeport：
```
- job_name: 'my-http-server'
  static_configs:
    - targets: ['192.168.160.6:30123', '192.168.160.5:30123', '192.168.160.4:30123']
```

对prometheus里的指标使用还不熟练，这里只展示了下`httpserver_execution_latency_second_bucket`的数据  

prometheus页面：
![image](https://gitee.com/tglgame/tools/raw/master/images/m10-prom.PNG)

grafana页面：
![image](https://gitee.com/tglgame/tools/raw/master/images/m10-grafana.PNG)

