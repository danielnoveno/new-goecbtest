/*
   file:           task/cleanup_mutex.go
   description:    Utility pembersih mutex scheduler
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package task

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// CleanMutex adalah fungsi untuk membersihkan mutex.
func CleanMutex() {
	mutexDir := "./storage/framework"
	fmt.Println("[CleanMutex] Scanning:", mutexDir)

	files, err := filepath.Glob(filepath.Join(mutexDir, "schedule-*"))
	if err != nil {
		log.Println("[CleanMutex] Error listing files:", err)
		return
	}

	now := time.Now()
	count := 0

	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			log.Println("[CleanMutex] Error stat:", file, err)
			continue
		}
		diff := now.Sub(info.ModTime()).Minutes()
		if diff > 60 {
			fmt.Printf("[CleanMutex] Deleting %s (last modified %s)\n", file, info.ModTime().Format("2006-01-02 15:04:05"))
			if err := os.Remove(file); err == nil {
				count++
			} else {
				log.Println("[CleanMutex] Failed to delete:", file, err)
			}
		}
	}

	fmt.Printf("[CleanMutex] Done. Deleted %d file(s).\n", count)
}
