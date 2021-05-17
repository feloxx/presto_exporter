package entity

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type StatusEntity struct {
	Coordinator     bool        `json:"coordinator"`
	Environment     string      `json:"environment"`
	ExternalAddress string      `json:"externalAddress"`
	HeapAvailable   int64       `json:"heapAvailable"`
	HeapUsed        int64       `json:"heapUsed"`
	InternalAddress string      `json:"internalAddress"`
	MemoryInfo      MemoryInfo  `json:"memoryInfo"`
	NodeID          string      `json:"nodeId"`
	NodeVersion     NodeVersion `json:"nodeVersion"`
	NonHeapUsed     int64       `json:"nonHeapUsed"`
	ProcessCPULoad  float64     `json:"processCpuLoad"`
	Processors      int64       `json:"processors"`
	SystemCPULoad   float64     `json:"systemCpuLoad"`
	Uptime          string      `json:"uptime"`
}

type MemoryDetail struct {
	FreeBytes                        int64            `json:"freeBytes"`
	MaxBytes                         int64            `json:"maxBytes"`
	QueryMemoryReservations          map[string]int64 `json:"queryMemoryReservations"`
	QueryMemoryRevocableReservations map[string]int64 `json:"queryMemoryRevocableReservations"`
	ReservedBytes                    int64            `json:"reservedBytes"`
	ReservedRevocableBytes           int64            `json:"reservedRevocableBytes"`
}

type Pools struct {
	General  MemoryDetail `json:"general"`
	Reserved MemoryDetail `json:"reserved"`
}

type MemoryInfo struct {
	Pools           Pools  `json:"pools"`
	TotalNodeMemory string `json:"totalNodeMemory"`
}

type NodeVersion struct {
	Version string `json:"version"`
}

func NewStatusEntity(prestoUrl, prestoPort string) (StatusEntity, error) {
	var status StatusEntity

	// url地址补充
	url := fmt.Sprintf("http://%s:%s/v1/status", prestoUrl, prestoPort)

	// get请求接口
	resp, err := http.Get(url)
	if err != nil {
		return status, errors.Wrap(err, "failed to get cluster metrics")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return status, errors.Wrap(err, "failed to read cluster response body")
	}
	defer resp.Body.Close()

	// 判断返回码
	if resp.StatusCode != 200 {
		return status, fmt.Errorf("failed to get metrics: %s %d", string(data), resp.StatusCode)
	}

	// 请求结果序列化为对象
	if err := json.Unmarshal(data, &status); err != nil {
		return status, errors.Wrapf(err, "failed to unmarshal cluster metrics output: %s", string(data))
	}

	return status, nil
}