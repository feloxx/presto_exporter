package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"my/presto_exporter/entity"
	"sync"
)

type clusterCollector struct {
	mutex  sync.Mutex
	cluster entity.ClusterEntity

	runningQueries   *prometheus.Desc // 正在查询个数
	blockedQueries   *prometheus.Desc // 阻塞查询个数
	queuedQueries    *prometheus.Desc // 排序查询个数
	activeWorkers    *prometheus.Desc // 当前存活工作节点
	runningDrivers   *prometheus.Desc // 当前运行driver个数
	reservedMemory   *prometheus.Desc // 当前使用保留内存
	totalInputRows   *prometheus.Desc // 当前读入数据的总行数
	totalInputBytes  *prometheus.Desc // 当前读入数据的总大小
	totalCPUTimeSecs *prometheus.Desc // cpu运行时间
}

// presto 集群状态监控项的构造方法
func NewClusterCollector(cluster entity.ClusterEntity) prometheus.Collector {
	subsystem := "cluster"

	return &clusterCollector{
		cluster: cluster,
		runningQueries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "running_queries"),
			"Running requests of the presto cluster.",
			nil, nil,
		),
		blockedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "blocked_queries"),
			"Blocked queries of the presto cluster.",
			nil, nil,
		),
		queuedQueries: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "queued_queries"),
			"Queued queries of the presto cluster.",
			nil, nil,
		),
		activeWorkers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "active_workers"),
			"Active workers of the presto cluster.",
			nil, nil,
		),
		runningDrivers: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "running_drivers"),
			"Running drivers of the presto cluster.",
			nil, nil,
		),
		reservedMemory: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reserved_memory_bytes"),
			"Reserved memory of the presto cluster.",
			nil, nil,
		),
		totalInputRows: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "input_rows_total"),
			"Total input rows of the presto cluster.",
			nil, nil,
		),
		totalInputBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "input_bytes_total"),
			"Total input bytes of the presto cluster.",
			nil, nil,
		),
		totalCPUTimeSecs: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "cpu_seconds_total"),
			"Total CPU time of the presto cluster.",
			nil, nil,
		),
	}
}

// 实现Describe接口，描述
func (c *clusterCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.runningQueries
	ch <- c.blockedQueries
	ch <- c.queuedQueries
	ch <- c.activeWorkers
	ch <- c.runningDrivers
	ch <- c.reservedMemory
	ch <- c.totalInputRows
	ch <- c.totalInputBytes
	ch <- c.totalCPUTimeSecs
}

// 实现Collect接口，进行数据采集
func (c *clusterCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	if err := c.cluster.ClusterCheck(); err != nil {
		log.Fatal("cluster collect err: ", err)
	}
	defer c.mutex.Unlock()

	ch <- prometheus.MustNewConstMetric(c.runningQueries, prometheus.GaugeValue, c.cluster.Info.RunningQueries)
	ch <- prometheus.MustNewConstMetric(c.blockedQueries, prometheus.GaugeValue, c.cluster.Info.BlockedQueries)
	ch <- prometheus.MustNewConstMetric(c.queuedQueries, prometheus.GaugeValue, c.cluster.Info.QueuedQueries)
	ch <- prometheus.MustNewConstMetric(c.activeWorkers, prometheus.GaugeValue, c.cluster.Info.ActiveWorkers)
	ch <- prometheus.MustNewConstMetric(c.runningDrivers, prometheus.GaugeValue, c.cluster.Info.RunningDrivers)
	ch <- prometheus.MustNewConstMetric(c.reservedMemory, prometheus.GaugeValue, c.cluster.Info.ReservedMemory)
	ch <- prometheus.MustNewConstMetric(c.totalInputRows, prometheus.CounterValue, c.cluster.Info.TotalInputRows)
	ch <- prometheus.MustNewConstMetric(c.totalInputBytes, prometheus.CounterValue, c.cluster.Info.TotalInputBytes)
	ch <- prometheus.MustNewConstMetric(c.totalCPUTimeSecs, prometheus.CounterValue, c.cluster.Info.TotalCPUTimeSecs)
}

