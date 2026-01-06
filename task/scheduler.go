/*
   file:           task/scheduler.go
   description:    Penjadwal sinkronisasi tugas ECB desktop
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package task

import (
	"fmt"
	"log"
	"sync"
	"time"

	"go-ecb/configs"
)

var startOnce sync.Once

// Start adalah fungsi untuk menjalankan.
func Start() {
	startOnce.Do(func() {
		fmt.Println("ECB Test Desktop Scheduler Running...")

		// Load scheduler intervals dari config
		cleanMutexInterval := time.Duration(configs.GetSchedulerCleanMutexInterval()) * time.Minute
		postDataInterval := time.Duration(configs.GetSchedulerPostDataInterval()) * time.Minute
		syncPoInterval := time.Duration(configs.GetSchedulerSyncPoInterval()) * time.Minute

		fmt.Printf("[Scheduler] Configured intervals - CleanMutex: %v, PostData: %v, SyncPo: %v\n",
			cleanMutexInterval, postDataInterval, syncPoInterval)

		go func() {
			CleanMutex()
			ticker := time.NewTicker(cleanMutexInterval)
			defer ticker.Stop()
			for range ticker.C {
				CleanMutex()
			}
		}()

		if err := ensureExternalDbsReachable(); err != nil {
			log.Printf("[scheduler] disable SIMO/BServ jobs: %v", err)
			return
		}

		go func() {
			allowedHours := map[int]bool{6: true, 9: true, 12: true, 15: true, 18: true, 21: true}
			runNow := true
			for {
				now := time.Now()
				if runNow || (allowedHours[now.Hour()] && now.Minute() == 0) {
					GetAllMasters()
					runNow = false
				}
				time.Sleep(1 * time.Minute)
			}
		}()

		go func() {
			PostEcbData()
			ticker := time.NewTicker(postDataInterval)
			defer ticker.Stop()
			for range ticker.C {
				PostEcbData()
			}
		}()

		go func() {
			SyncEcbPo()
			ticker := time.NewTicker(syncPoInterval)
			defer ticker.Stop()
			for range ticker.C {
				SyncEcbPo()
			}
		}()
	})
}

// ensureExternalDbsReachable adalah fungsi untuk memastikan external dbs reachable.
func ensureExternalDbsReachable() error {
	cfg := configs.LoadConfig()
	check := []struct {
		name string
		dsn  string
	}{
		{name: "simoprd", dsn: buildDSN(
			cfg.SimoprdUser,
			cfg.SimoprdPassword,
			fmt.Sprintf("%s:%s", cfg.SimoprdHost, cfg.SimoprdPort),
			cfg.SimoprdDatabase,
		)},
		{name: "bserv", dsn: buildDSN(
			cfg.BservUser,
			cfg.BservPassword,
			fmt.Sprintf("%s:%s", cfg.BservHost, cfg.BservPort),
			cfg.BservDatabase,
		)},
	}

	for _, target := range check {
		if err := pingDSN(target.dsn); err != nil {
			return fmt.Errorf("%s ping failed: %w", target.name, err)
		}
	}
	return nil
}
