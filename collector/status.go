package collector

import (
	"github.com/prometheus/client_golang/prometheus"
	"my/presto_exporter/entity"
	"strconv"
	"sync"
)

type statusCollector struct {
	mutex  sync.Mutex
	status entity.StatusEntity

	uptime          *prometheus.Desc
	totalNodeMemory *prometheus.Desc // 当前节点分配的内存

	reservedMaxBytes                         *prometheus.Desc // reserved 区最大内存
	reservedReservedBytes                    *prometheus.Desc // reserved 区保留内存
	reservedRevocableBytes                   *prometheus.Desc // reserved 区可回收内存
	reservedQueryMemoryReservations          *prometheus.Desc // reserved 区正在查询的任务 map形式显示
	reservedQueryMemoryRevocableReservations *prometheus.Desc // reserved 区可撤销 map形式显示
	reservedFreeBytes                        *prometheus.Desc // reserved 区空闲内存

	generalMaxBytes                         *prometheus.Desc // general 区最大内存
	generalReservedBytes                    *prometheus.Desc // general 区保留内存
	generalRevocableBytes                   *prometheus.Desc // general 区可回收内存
	generalQueryMemoryReservations          *prometheus.Desc // general 区正在查询的任务 map形式显示
	generalQueryMemoryRevocableReservations *prometheus.Desc // general 区可撤销 map形式显示
	generalFreeBytes                        *prometheus.Desc // general 区空闲内存

	processors     *prometheus.Desc // cpu核数
	processCpuLoad *prometheus.Desc // 核数cpu使用
	systemCpuLoad  *prometheus.Desc // 系统cpu使用
	heapUsed       *prometheus.Desc // 当前节点堆内存使用
	heapAvailable  *prometheus.Desc // 当前节点堆内存可用容量
	nonHeapUsed    *prometheus.Desc // 非堆内存使用
}

// presto 集群状态监控项的构造方法
func NewStatusCollector(status entity.StatusEntity) prometheus.Collector {
	subsystem := "status"

	return &statusCollector{
		status: status,
		uptime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "uptime"),
			"presto node uptime.",
			nil, nil,
		),
		totalNodeMemory: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "total_node_memory"),
			"presto node total memory.",
			nil, nil,
		),
		reservedMaxBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reserved_max_bytes"),
			"reservedMaxBytes",
			nil, nil,
		),
		reservedReservedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reserved_reserved_bytes"),
			"reservedReservedBytes",
			nil, nil,
		),
		reservedRevocableBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reserved_revocable_bytes"),
			"reservedRevocableBytes",
			nil, nil,
		),
		reservedQueryMemoryReservations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reserved_query_memory_reservations"),
			"reservedQueryMemoryReservations",
			nil, nil,
		),
		reservedQueryMemoryRevocableReservations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reserved_query_memory_revocable_reservations"),
			"reservedQueryMemoryRevocableReservations",
			nil, nil,
		),
		reservedFreeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "reserved_free_bytes"),
			"reservedFreeBytes",
			nil, nil,
		),
		generalMaxBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "general_max_bytes"),
			"generalMaxBytes",
			nil, nil,
		),
		generalReservedBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "general_reserved_bytes"),
			"generalReservedBytes",
			nil, nil,
		),
		generalRevocableBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "general_revocable_bytes"),
			"generalRevocableBytes",
			nil, nil,
		),
		generalQueryMemoryReservations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "general_query_memory_reservations"),
			"generalQueryMemoryReservations",
			nil, nil,
		),
		generalQueryMemoryRevocableReservations: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "general_query_memory_revocable_reservations"),
			"generalQueryMemoryRevocableReservations",
			nil, nil,
		),
		generalFreeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "general_free_bytes"),
			"generalFreeBytes",
			nil, nil,
		),
		processors: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "processors"),
			"processors",
			nil, nil,
		),
		processCpuLoad: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "process_cpu_load"),
			"processCpuLoad",
			nil, nil,
		),
		systemCpuLoad: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "system_cpu_load"),
			"systemCpuLoad",
			nil, nil,
		),
		heapUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "heap_used"),
			"heapUsed",
			nil, nil,
		),
		heapAvailable: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "heap_available"),
			"heapAvailable",
			nil, nil,
		),
		nonHeapUsed: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, subsystem, "non_heap_used"),
			"nonHeapUsed",
			nil, nil,
		),
	}
}

// 实现Describe接口，描述
func (c *statusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.uptime
	ch <- c.totalNodeMemory

	ch <- c.reservedMaxBytes
	ch <- c.reservedReservedBytes
	ch <- c.reservedRevocableBytes
	ch <- c.reservedQueryMemoryReservations
	ch <- c.reservedQueryMemoryRevocableReservations
	ch <- c.reservedFreeBytes

	ch <- c.generalMaxBytes
	ch <- c.generalReservedBytes
	ch <- c.generalRevocableBytes
	ch <- c.generalQueryMemoryReservations
	ch <- c.generalQueryMemoryRevocableReservations
	ch <- c.generalFreeBytes

	ch <- c.processors
	ch <- c.processCpuLoad
	ch <- c.systemCpuLoad
	ch <- c.heapUsed
	ch <- c.heapAvailable
	ch <- c.nonHeapUsed
}

// 实现Collect接口，进行数据采集
func (c *statusCollector) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	uptime, _ := strconv.ParseFloat(c.status.Uptime[:len(c.status.Uptime)-1], 64)
	ch <- prometheus.MustNewConstMetric(c.uptime, prometheus.GaugeValue, uptime)

	totalNodeMemory, _ := strconv.ParseInt(c.status.MemoryInfo.TotalNodeMemory[:len(c.status.MemoryInfo.TotalNodeMemory)-1], 10, 64)
	ch <- prometheus.MustNewConstMetric(c.totalNodeMemory, prometheus.GaugeValue, float64(totalNodeMemory))

	// reserved
	ch <- prometheus.MustNewConstMetric(c.reservedMaxBytes, prometheus.GaugeValue, float64(c.status.MemoryInfo.Pools.Reserved.MaxBytes))
	ch <- prometheus.MustNewConstMetric(c.reservedReservedBytes, prometheus.CounterValue, float64(c.status.MemoryInfo.Pools.Reserved.ReservedBytes))

	// 线程池暂时用是否有占用来替代，后续调整
	reservedQueryMemoryReservations := 0
	if len(c.status.MemoryInfo.Pools.Reserved.QueryMemoryReservations) != 0 {
		reservedQueryMemoryReservations = 1
	}
	ch <- prometheus.MustNewConstMetric(c.reservedQueryMemoryReservations, prometheus.CounterValue, float64(reservedQueryMemoryReservations))
	reservedQueryMemoryRevocableReservations := 0
	if len(c.status.MemoryInfo.Pools.Reserved.QueryMemoryRevocableReservations) != 0 {
		reservedQueryMemoryRevocableReservations = 1
	}
	ch <- prometheus.MustNewConstMetric(c.reservedQueryMemoryRevocableReservations, prometheus.CounterValue, float64(reservedQueryMemoryRevocableReservations))

	ch <- prometheus.MustNewConstMetric(c.reservedFreeBytes, prometheus.CounterValue, float64(c.status.MemoryInfo.Pools.Reserved.FreeBytes))

	// general
	ch <- prometheus.MustNewConstMetric(c.generalMaxBytes, prometheus.GaugeValue, float64(c.status.MemoryInfo.Pools.General.MaxBytes))
	ch <- prometheus.MustNewConstMetric(c.generalReservedBytes, prometheus.CounterValue, float64(c.status.MemoryInfo.Pools.General.ReservedBytes))

	// 线程池暂时用是否有占用来替代，后续调整
	generalQueryMemoryReservations := 0
	if len(c.status.MemoryInfo.Pools.Reserved.QueryMemoryReservations) != 0 {
		generalQueryMemoryReservations = 1
	}
	ch <- prometheus.MustNewConstMetric(c.generalQueryMemoryReservations, prometheus.CounterValue, float64(generalQueryMemoryReservations))
	generalQueryMemoryRevocableReservations := 0
	if len(c.status.MemoryInfo.Pools.Reserved.QueryMemoryRevocableReservations) != 0 {
		generalQueryMemoryRevocableReservations = 1
	}
	ch <- prometheus.MustNewConstMetric(c.generalQueryMemoryRevocableReservations, prometheus.CounterValue, float64(generalQueryMemoryRevocableReservations))

	ch <- prometheus.MustNewConstMetric(c.generalFreeBytes, prometheus.CounterValue, float64(c.status.MemoryInfo.Pools.General.FreeBytes))

	// other
	ch <- prometheus.MustNewConstMetric(c.processors,prometheus.CounterValue, float64(c.status.Processors))
	ch <- prometheus.MustNewConstMetric(c.processCpuLoad,prometheus.CounterValue, float64(c.status.ProcessCPULoad))
	ch <- prometheus.MustNewConstMetric(c.systemCpuLoad,prometheus.CounterValue, float64(c.status.SystemCPULoad))
	ch <- prometheus.MustNewConstMetric(c.heapUsed,prometheus.CounterValue, float64(c.status.HeapUsed))
	ch <- prometheus.MustNewConstMetric(c.heapAvailable,prometheus.CounterValue, float64(c.status.HeapAvailable))
	ch <- prometheus.MustNewConstMetric(c.nonHeapUsed,prometheus.CounterValue, float64(c.status.NonHeapUsed))
}
