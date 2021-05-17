# presto_exporter

使用 `Go` 开发的 `exporter`。

按照 `Prometheus` 的官方推荐的代码架构：

```
presto_exporter/
├── collector   // 符合prometheus的数据监控的代码
├── entity      // 监控信息的实体
├── main.go     // exporter 入口
├── go.mod
└── go.sum
```

参数解释

- **web.listen-address** presto_exporter的工作地址和端口，一般建议是本机ip端口为9088
- **presto.coordinator.url** 监控coordinator时，需要配置此参数，填写coordinator的地址
- **presto.worker.url** 监控worker时，需要配置此参数，填写worker的地址
- **presto.port** presto的工作端口

注意：

以上4个参数都需要配置，`presto.coordinator.url` 与 `presto.worker.url` 参数互斥。

启动例子：

```
./presto_exporter --web.listen-address 10.150.31.29:9088 --presto.coordinator.url 10.150.31.29 --presto.port 5797

./presto_exporter --web.listen-address 10.150.31.31:9088 --presto.worker.url 10.150.31.31 --presto.port 5797
```