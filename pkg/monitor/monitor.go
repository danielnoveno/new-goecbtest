/*
   file:           pkg/monitor/monitor.go
   description:    Resource monitoring dashboard untuk ECB desktop app
   created:        optimization 05-01-2026
*/

package monitor

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"go-ecb/configs"
)

var (
	monitorOnce sync.Once
	// GpioOpsCount adalah counter untuk GPIO operations (increment dari external)
	GpioOpsCount    int64
	gpioOpsMu       sync.Mutex
	lastGpioOpsTime time.Time
)

// Stats menyimpan statistik resource usage
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

// IncrementGpioOps increment GPIO operation counter (dipanggil dari gpio package)
func IncrementGpioOps() {
	gpioOpsMu.Lock()
	defer gpioOpsMu.Unlock()
	GpioOpsCount++
	lastGpioOpsTime = time.Now()
}

// Start memulai resource monitoring dashboard jika diaktifkan
func Start(db *sql.DB) {
	if !configs.GetEnableResourceMonitor() {
		return
	}

	monitorOnce.Do(func() {
		go runMonitor(db)
	})
}

// runMonitor menjalankan monitoring loop
func runMonitor(db *sql.DB) {
	interval := time.Duration(configs.GetResourceMonitorInterval()) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Println("\n╔══════════════════════════════════════════════════════════╗")
	fmt.Println("║         ECB Test Resource Monitor Started              ║")
	fmt.Println("║  Press Ctrl+C to stop application                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════╝\n")

	// Baseline stats
	var lastCPUTime uint64
	var lastSysTime time.Time
	lastGpioOps := int64(0)

	for range ticker.C {
		stats := collectStats(db, &lastCPUTime, &lastSysTime, &lastGpioOps)
		displayStats(stats)
	}
}

// collectStats mengumpulkan statistik resource saat ini
func collectStats(db *sql.DB, lastCPUTime *uint64, lastSysTime *time.Time, lastGpioOps *int64) Stats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	stats := Stats{
		CPUPercent:     calculateCPUPercent(lastCPUTime, lastSysTime),
		MemoryUsedMB:   m.Alloc / 1024 / 1024,
		MemoryTotalMB:  m.Sys / 1024 / 1024,
		GoroutineCount: runtime.NumGoroutine(),
		Timestamp:      time.Now(),
	}

	// DB connection stats
	if db != nil {
		dbStats := db.Stats()
		stats.DBConnActive = dbStats.InUse
		stats.DBConnIdle = dbStats.Idle
	}

	// GPIO ops per second
	gpioOpsMu.Lock()
	currentOps := GpioOpsCount
	gpioOpsMu.Unlock()

	if !lastSysTime.IsZero() {
		elapsed := time.Since(*lastSysTime).Seconds()
		if elapsed > 0 {
			stats.GpioOpsPerSec = float64(currentOps-*lastGpioOps) / elapsed
		}
	}
	*lastGpioOps = currentOps

	return stats
}

// calculateCPUPercent menghitung persentase penggunaan CPU
func calculateCPUPercent(lastCPUTime *uint64, lastSysTime *time.Time) float64 {
	// CPU calculation simplified - could be enhanced with actual CPU sampling
	// For now, use goroutine count as proxy (very rough estimate)
	goroutines := runtime.NumGoroutine()
	
	// Rough estimate: assume baseline 10 goroutines, scale up to 100%
	// This is very approximate - for real CPU%, would need to sample over time
	percent := float64(goroutines-10) * 2.0
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	*lastSysTime = time.Now()
	return percent
}

// displayStats menampilkan statistik di terminal
func displayStats(stats Stats) {
	// Clear screen untuk refresh (ANSI escape codes)
	// fmt.Print("\033[H\033[2J") // Disable untuk kompatibilitas Windows

	// Buat progress bar
	cpuBar := makeProgressBar(int(stats.CPUPercent), 50)
	memPercent := 0.0
	if stats.MemoryTotalMB > 0 {
		memPercent = float64(stats.MemoryUsedMB) / float64(stats.MemoryTotalMB) * 100
	}
	memBar := makeProgressBar(int(memPercent), 50)

	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("  ECB Test Resource Monitor - %s\n", stats.Timestamp.Format("15:04:05"))
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	
	fmt.Printf("\n  CPU Usage:      %s %.1f%%\n", cpuBar, stats.CPUPercent)
	fmt.Printf("  Memory:         %s %.1f%% (%d / %d MB)\n", 
		memBar, memPercent, stats.MemoryUsedMB, stats.MemoryTotalMB)
	fmt.Printf("  Goroutines:     %d\n", stats.GoroutineCount)
	
	if stats.DBConnActive > 0 || stats.DBConnIdle > 0 {
		fmt.Printf("  DB Connections: %d active, %d idle\n", 
			stats.DBConnActive, stats.DBConnIdle)
	}
	
	if stats.GpioOpsPerSec > 0 {
		fmt.Printf("  GPIO Ops/sec:   %.1f\n", stats.GpioOpsPerSec)
	}
	
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println()
}

// makeProgressBar membuat ASCII progress bar
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

// LogStats mencatat stats ke file untuk analisis
func LogStats(stats Stats, filename string) error {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	line := fmt.Sprintf("%s,%.1f,%d,%d,%d,%d,%d,%.2f\n",
		stats.Timestamp.Format(time.RFC3339),
		stats.CPUPercent,
		stats.MemoryUsedMB,
		stats.MemoryTotalMB,
		stats.GoroutineCount,
		stats.DBConnActive,
		stats.DBConnIdle,
		stats.GpioOpsPerSec,
	)

	if _, err := f.WriteString(line); err != nil {
		return err
	}

	return nil
}

// PrintSystemInfo mencetak informasi sistem saat startup
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
