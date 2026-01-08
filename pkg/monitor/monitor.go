/*
   file:           pkg/monitor/monitor.go
   description:    Resource monitoring dashboard untuk ECB desktop app
   created:        optimization 05-01-2026
*/

package monitor

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"sync"
	"time"

	"go-ecb/configs"
)

var (
	monitorOnce sync.Once
	GpioOpsCount    int64
	gpioOpsMu       sync.Mutex
	lastGpioOpsTime time.Time
)

type Stats struct {
	CPUPercent      float64
	MemoryUsedMB    uint64
	MemoryTotalMB   uint64
	GoroutineCount  int
	DBConnActive    int
	DBConnIdle      int
	GpioOpsPerSec   float64
	Timestamp       time.Time
}

// func IncrementGpioOps() {
// 	gpioOpsMu.Lock()
// 	defer gpioOpsMu.Unlock()
// 	GpioOpsCount++
// 	lastGpioOpsTime = time.Now()
// }

func Start() {
	if !configs.GetEnableResourceMonitor() {
		return
	}

	monitorOnce.Do(func() {
		go runMonitor()
	})
}

func runMonitor() {
	interval := time.Duration(configs.GetResourceMonitorInterval()) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Println("======== ECB Test Resource Monitor Started ========")
	fmt.Println("======== Press Ctrl+C to stop application ========")
	fmt.Println()

	var lastSysTime time.Time
	lastGpioOps := int64(0)

	for range ticker.C {
		stats := collectStats(&lastSysTime, &lastGpioOps)
		displayStats(stats)
	}
}

func collectStats(lastSysTime *time.Time, lastGpioOps *int64) Stats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := Stats{
		CPUPercent:     calculateCPUPercent(lastSysTime),
		MemoryUsedMB:   m.Alloc / 1024 / 1024,
		MemoryTotalMB:  m.Sys / 1024 / 1024,
		GoroutineCount: runtime.NumGoroutine(),
		Timestamp:      time.Now(),
	}

	// if db != nil {
	// 	dbStats := db.Stats()
	// 	stats.DBConnActive = dbStats.InUse
	// 	stats.DBConnIdle = dbStats.Idle
	// }

	gpioOpsMu.Lock()
	currentOps := GpioOpsCount
	gpioOpsMu.Unlock()

	// apakah start pertama
	if !lastSysTime.IsZero() {
		elapsed := time.Since(*lastSysTime).Seconds() // hitung selisih waktu
		if elapsed > 0 { 
			stats.GpioOpsPerSec = float64(currentOps-*lastGpioOps) / elapsed // hitung ops per second
		}
	}
	*lastGpioOps = currentOps // simpan ops terakhir

	return stats
}

func calculateCPUPercent(lastSysTime *time.Time) float64 {
	goroutines := runtime.NumGoroutine() // jumlah goroutine yang sedang berjalan saat ini
	percent := float64(goroutines-10) * 2.0 // estimasi persentase CPU: (jumlah goroutine - 10) * 2

	if percent < 0 { // nilai persentase kurang dari 0
		percent = 0
	}
	if percent > 100 { // nilai persentase lebih dari 100
		percent = 100
	}

	*lastSysTime = time.Now() // waktu pengecekan terakhir ke waktu sekarang

	return percent
}

func displayStats(stats Stats) {
	// persentase memori
	memPercent := 0.0
	if stats.MemoryTotalMB > 0 {
		memPercent = float64(stats.MemoryUsedMB) / float64(stats.MemoryTotalMB) * 100 // hitung persentase memori: (digunakan / total) * 100
	}

	fmt.Printf("\n--- Monitor: %s ---\n", stats.Timestamp.Format("15:04:05"))

	// tampilkan visual progress bar
	fmt.Printf("CPU: %s %.1f%%\n", makeProgressBar(int(stats.CPUPercent), 40), stats.CPUPercent)
	fmt.Printf("MEM: %s %.1f%% (%d/%d MB)\n", makeProgressBar(int(memPercent), 40), memPercent, stats.MemoryUsedMB, stats.MemoryTotalMB)

	// ringkasan data
	fmt.Printf("Goroutines: %d | GPIO: %.1f ops/s | DB: %d active\n",
		stats.GoroutineCount, stats.GpioOpsPerSec, stats.DBConnActive)
}

func makeProgressBar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := (percent * width) / 100
	empty := width - filled

	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
}

// func LogStats(stats Stats, filename string) error {
// 	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil {
// 		return err
// 	}
// 	defer f.Close()

// 	line := fmt.Sprintf("%s,%.1f,%d,%d,%d,%d,%d,%.2f\n",
// 		stats.Timestamp.Format(time.RFC3339),
// 		stats.CPUPercent,
// 		stats.MemoryUsedMB,
// 		stats.MemoryTotalMB,
// 		stats.GoroutineCount,
// 		stats.DBConnActive,
// 		stats.DBConnIdle,
// 		stats.GpioOpsPerSec,
// 	)

// 	if _, err := f.WriteString(line); err != nil {
// 		return err
// 	}

// 	return nil
// }

func PrintSystemInfo() {
	log.Printf("[Monitor] System Info:")
	log.Printf("[Monitor]   - OS/Arch: %s/%s", runtime.GOOS, runtime.GOARCH)
	log.Printf("[Monitor]   - CPUs: %d", runtime.NumCPU())
	log.Printf("[Monitor]   - Go Version: %s", runtime.Version())
	
	if configs.IsRaspberryPi() {
		log.Printf("[Monitor]   - Platform: Raspberry Pi (detected)")
	} else {
		log.Printf("[Monitor]   - Platform: Desktop/Server")
	}
}
