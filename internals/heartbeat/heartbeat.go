// The heartbreat module
// is a module that reports the currunt machine resource usage to redis
// CPU, RAM
// this information is used to loadbalance connections

package heartbeat

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type HeartbeatStatus struct{}

const RamWight = 0.6
const CpuWeight = 0.4
const waitFor = time.Second * 5

func (*HeartbeatStatus) Run() {

	ticker := time.NewTicker(waitFor)
	for {
		memoryUsage, err := mem.VirtualMemory()

		if err != nil {
			time.Sleep(waitFor)

			continue
		}

		cpuStats, err := cpu.Percent(0, false)

		if err != nil {
			time.Sleep(waitFor)
			continue
		}
		UsedPercent := cpuStats[0]

		cpuWeight := 0.4 // 40% weight
		ramWeight := 0.6 // 60% weight (since I think RAM is more important)

		baseScore := (UsedPercent * cpuWeight) + (memoryUsage.UsedPercent * ramWeight)

		// Factor in connections as a separate multiplier

		fmt.Println(baseScore)
		<-ticker.C
	}

}

// New creates a new HeartbeatStatus used to report machine resource usage (CPU and RAM).
func New() *HeartbeatStatus {
	return &HeartbeatStatus{}
}
