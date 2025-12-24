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

func (*HeartbeatStatus) Run() {

	ticker := time.NewTicker(time.Second * 5)
	for {
		memoryUsage, err := mem.VirtualMemory()

		if err != nil {
			// fatal notify admin
		}

		cpuStats, err := cpu.Percent(0, false)

		if err != nil {
			// fatal notify admin
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
func New() *HeartbeatStatus {
	return &HeartbeatStatus{}
}
