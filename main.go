package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"my/presto_exporter/collector"
	"my/presto_exporter/entity"
	"net/http"
	"os"
)


var (
	version = "1.0.0"

	listenAddress = flag.String("web.listen-address","127.0.0.1:9088","presto_export http url and port.")
	coordinatorUrl = flag.String("presto.coordinator.url","","Presto coordinator url, Such as 192.168.1.100")
	workerUrl = flag.String("presto.worker.url","","Presto worker url, Such as 192.168.1.101")
	Port = flag.String("presto.port","","Presto coordinator/worker port, Such as 5797")
)

func main() {
	log.Infof("starting presto_exporter %s...", version)

	// 读取参数
	flag.Parse()

	// 参数检查
	// TODO
	if *coordinatorUrl != "" && *workerUrl != "" {
		log.Fatalf("Choose between coordinator and worker")
		os.Exit(1)
	}

	// 实例化监控
	if *coordinatorUrl != "" {
		cluster, err := entity.NewClusterEntity(*coordinatorUrl, *Port)
		if err != nil {
			log.Fatalf("failed to get presto cluster: %v", err)
		}
		prometheus.MustRegister(collector.NewClusterCollector(cluster))
	}

	url := ""
	if *coordinatorUrl != "" {
		url = *coordinatorUrl
	} else {
		url = *workerUrl
	}

	status, err := entity.NewStatusEntity(url, *Port)
	if err != nil {
		log.Fatalf("failed to get presto status: %v", err)
	}

	// 向prometheus注册collector
	prometheus.MustRegister(collector.NewStatusCollector(status))

	// 启动一个http接口，用于prometheus拉取数据
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
			<head><title>Presto Exporter</title></head>
			<body>
				<h1>Presto Exporter</h1>
				<p>Available for Presto 0.203</p>
				<p><a href=/metrics>metrics</a></p>
			</body>
			</html>
		`)
	})
	log.Infof("presto_exporter listening on %s", *listenAddress)
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
