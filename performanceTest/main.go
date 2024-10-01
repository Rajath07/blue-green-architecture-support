package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

// Convert bytes to megabytes
func bytesToMegabytes(b uint64) float64 {
	return float64(b) / (1024 * 1024)
}

func main() {
	const numGoroutines = 100000
	var wg sync.WaitGroup
	channels := make([]chan int, numGoroutines)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	startAlloc := memStats.Alloc
	startCpuPercent, _ := cpu.Percent(0, false)
	startTime := time.Now()

	// Spawn the goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		ch := make(chan int)
		channels[i] = ch
		go func(ch chan int) {
			defer wg.Done()
			for {
				select {
				case <-ch:
					// Do nothing, just keep the goroutine blocked on this channel
				}
			}
		}(ch)
	}

	// Wait for all goroutines to start and hit their blocking state
	time.Sleep(5 * time.Second)

	// Measure resource utilization again
	runtime.ReadMemStats(&memStats)
	endAlloc := memStats.Alloc
	endCpuPercent, _ := cpu.Percent(0, false)
	endTime := time.Now()

	fmt.Printf("Memory Allocation Before: %d bytes (%.2f MB)\n", startAlloc, bytesToMegabytes(startAlloc))
	fmt.Printf("Memory Allocation After: %d bytes (%.2f MB)\n", endAlloc, bytesToMegabytes(endAlloc))
	fmt.Printf("Memory Allocation Difference: %d bytes (%.2f MB)\n", endAlloc-startAlloc, bytesToMegabytes(endAlloc-startAlloc))

	fmt.Printf("CPU Utilization Before: %.2f%%\n", startCpuPercent[0])
	fmt.Printf("CPU Utilization After: %.2f%%\n", endCpuPercent[0])

	fmt.Printf("Time Elapsed: %s\n", endTime.Sub(startTime))

}
