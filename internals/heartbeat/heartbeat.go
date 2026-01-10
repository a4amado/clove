// The heartbreat module
// is a module that reports the currunt machine resource usage to redis
// CPU, RAM
// this information is used to loadbalance connections

package heartbeat

import (
	"clove/internals/heartbeat/dogpile"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type HeartbeatStatus struct{}

const ramWeight = 0.6
const cpuWeight = 0.4
const waitFor = time.Second * 5

var dogpileInstance = dogpile.New()

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
		usedPercent := cpuStats[0]

		localCpuWeight := 0.3 // 40% weight
		localRamWeight := 0.5  // 60% weight (since I think RAM is more important)
		// numOfConn := dogpileInstance.GetCurrentCount()
		baseScore := (usedPercent * localCpuWeight) + (memoryUsage.UsedPercent * localRamWeight)

		// Factor in connections as a separate multiplier

		fmt.Println(baseScore)
		<-ticker.C
	}

}

// New creates a new HeartbeatStatus used to report machine resource usage (CPU and RAM).
func New() *HeartbeatStatus {
	return &HeartbeatStatus{}
}
