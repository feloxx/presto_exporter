# presto_exporter

使用 `Go` 开发的 `presto_exporter`。

符合prometheus的metrics获取，在面向java应用中有常见的三种方式：

- 使用jmx_exporter以java agent的方式启动，程序开发与运维强耦合，调整、升级、维护成本较高。
- 本身应用支持符合prometheus的metrics，这需要在应用中引入prometheus的包来定制开发（比如java需要在pom.xml应用相关包，然后按照规范来开发一些监控项），调整、升级、维护成本归属到了应用开发中，与运维脱离。成本都在开发上。
- 开发应用的专属exporter，由这个单独的程序来提供符合prometheus的监控项，程序开发与运维低耦合，调整、升级、维护成本较低。

这里是使用应用专属exporter来实现的。所以我们描述一下数据的获取的大致流程（可以简单理解为我启动了一个HTTP服务器，然后去调用相关应用预留的运维接口获得数据，然后组装成metrics的形式返回）：

```
定义一个GET的metrics接口，然后像prometheus中注册一个collector。

我们可以注册多个collector，可以把它理解我们要监控的分类。

collector是prometheus定义监控的一种接口，在该接口中我们主要实现以下两个方法：

Describe(chan<- *Desc) // 收集监控项
Collect(chan<- Metric) // 具体的监控数据采集工作

最后将该接口与promhttp.Handler()绑定即可。
```

---

代码结构：

```
presto_exporter/
├── collector   // 采集接口实现
├── entity      // 具体的监控项实体，以及采集实现
├── main.go     // 入口
├── go.mod
└── go.sum
```

---

参数解释

- **web.listen-address** presto_exporter的工作地址和端口，一般建议是本机ip端口为9088
- **presto.coordinator.url** 监控coordinator时，需要配置此参数，填写coordinator的地址
- **presto.worker.url** 监控worker时，需要配置此参数，填写worker的地址
- **presto.port** presto的工作端口

注意：

以上4个参数都需要配置，`presto.coordinator.url` 与 `presto.worker.url` 参数互斥。

启动例子：

```
// 启动对presto的coordinator监控
./presto_exporter --web.listen-address 10.150.31.29:9088 --presto.coordinator.url 10.150.31.29 --presto.port 5797

// 启动对presto的worker监控
./presto_exporter --web.listen-address 10.150.31.31:9088 --presto.worker.url 10.150.31.31 --presto.port 5797
```