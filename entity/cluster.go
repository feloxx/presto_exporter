package entity

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

type ClusterEntity struct {
	RunningQueries   float64 `json:"runningQueries"`
	BlockedQueries   float64 `json:"blockedQueries"`
	QueuedQueries    float64 `json:"queuedQueries"`
	ActiveWorkers    float64 `json:"activeWorkers"`
	RunningDrivers   float64 `json:"runningDrivers"`
	ReservedMemory   float64 `json:"reservedMemory"`
	TotalInputRows   float64 `json:"totalInputRows"`
	TotalInputBytes  float64 `json:"totalInputBytes"`
	TotalCPUTimeSecs float64 `json:"totalCpuTimeSecs"`
}

func NewClusterEntity(prestoUrl, prestoPort string) (ClusterEntity, error) {
	var cluster ClusterEntity

	// url地址补充
	url := fmt.Sprintf("http://%s:%s/v1/cluster", prestoUrl, prestoPort)

	// get请求接口
	resp, err := http.Get(url)
	if err != nil {
		return cluster, errors.Wrap(err, "failed to get cluster metrics")
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return cluster, errors.Wrap(err, "failed to read cluster response body")
	}
	defer resp.Body.Close()

	// 判断返回码
	if resp.StatusCode != 200 {
		return cluster, fmt.Errorf("failed to get metrics: %s %d", string(data), resp.StatusCode)
	}

	// 请求结果序列化为对象
	if err := json.Unmarshal(data, &cluster); err != nil {
		return cluster, errors.Wrapf(err, "failed to unmarshal cluster metrics output: %s", string(data))
	}

	return cluster, nil
}



