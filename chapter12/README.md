# 说明： 
通过`istioctl`安装好`istio`之后，参考了官方提供的[例子](https://github.com/istio/istio/tree/master/samples/helloworld)部署`httpserver`

## istio服务的状态
```
root@k8snode1:/home/CNCF_TEST# k get all -n istio-system
NAME                                        READY   STATUS    RESTARTS        AGE
pod/istio-egressgateway-687f4db598-j2l2h    1/1     Running   3 (2d23h ago)   7d1h
pod/istio-ingressgateway-78f69bd5db-pbfm9   1/1     Running   3 (2d23h ago)   7d1h
pod/istiod-76d66d9876-n48kl                 1/1     Running   4 (2d23h ago)   7d1h

NAME                           TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                      AGE
service/istio-egressgateway    ClusterIP      10.98.176.106   <none>        80/TCP,443/TCP                                                               7d1h
service/istio-ingressgateway   LoadBalancer   10.101.91.161   <pending>     15021:32519/TCP,80:31129/TCP,443:30629/TCP,31400:32309/TCP,15443:31891/TCP   7d1h
service/istiod                 ClusterIP      10.99.44.236    <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP                                        7d1h

NAME                                   READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/istio-egressgateway    1/1     1            1           7d1h
deployment.apps/istio-ingressgateway   1/1     1            1           7d1h
deployment.apps/istiod                 1/1     1            1           7d1h

NAME                                              DESIRED   CURRENT   READY   AGE
replicaset.apps/istio-egressgateway-687f4db598    1         1         1       7d1h
replicaset.apps/istio-ingressgateway-78f69bd5db   1         1         1       7d1h
replicaset.apps/istiod-76d66d9876                 1         1         1       7d1h

```


## 文件说明和验证
`conf.json`内容是`httpserver`需要的配置文件，`cm.sh`会利用`conf.json`生成configmap，以卷的形式mount到pod里面

`deploy-gateway.yaml`的内容完全参考了官方例子，整体比较简单，记录下INGRESS_PORT是怎么获取的：
```
export INGRESS_PORT=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
```

服务运行完成之后，在浏览器中输入：  `http://192.168.160.4:31129/hello`   
因为设置的replicas为2，所以请求会随机发送到其中1个pod里, 网页显示如下：

```
gracedeploy-5f8fbb7bc-jvrqf, hello world!
```
或
```
gracedeploy-5f8fbb7bc-d6cmb, hello world!
```

## 改进
本次只是部署istio之后，参考官方例子将服务基于istio运行了起来，对于istio的理解依然比较浅，作业要求中提到的几点也没有完全做完：  
* 如何实现安全保证
* 七层路由规则
* 考虑 open tracing 的接入    

通过本次实践，学习了istio的基本概念和用法，高级特性还需要多摸索和尝试；其内部用到的`envoy`也是值得花时间研究一下