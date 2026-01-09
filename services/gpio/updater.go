/*
   file:           services/gpio/updater.go
   description:    Driver GPIO untuk updater
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package gpio

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"go-ecb/configs"
	"go-ecb/pkg/logging"

	"github.com/go-gorp/gorp"
)

var (
	ecbStateUpdaterOnce sync.Once
)

func StartEcbStateUpdater(dbmap *gorp.DbMap) {
	if dbmap == nil {
		return
	}
	ecbStateUpdaterOnce.Do(func() {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					logging.Logger().Errorf("[EcbStateUpdater] panic recovered: %v", r)
				}
			}()
			runEcbStateUpdater(dbmap)
		}()
	})
}

func runEcbStateUpdater(dbmap *gorp.DbMap) {
	pollIntervalMs := configs.GetGpioPollIntervalMs()
	adaptivePolling := configs.GetGpioAdaptivePolling()
	
	basePollInterval := time.Duration(pollIntervalMs) * time.Millisecond
	currentPollInterval := basePollInterval
	
	maxSlowInterval := basePollInterval * 4
	unchangedCount := 0
	
	ticker := time.NewTicker(currentPollInterval)
	defer ticker.Stop()

	var lastState string
	for range ticker.C {
		current := readEcbStateString()
		if current == "" {
			continue
		}
		
		if current != lastState && lastState != "" {
			logStateChange(lastState, current)
		}

		if current == lastState {
			unchangedCount++
			
			if adaptivePolling && unchangedCount > 4 {
				newInterval := currentPollInterval * 2
				if newInterval > maxSlowInterval {
					newInterval = maxSlowInterval
				}
				if newInterval != currentPollInterval {
					currentPollInterval = newInterval
					ticker.Stop()
					ticker = time.NewTicker(currentPollInterval)
				}
			}
			continue
		}
		
		if currentPollInterval != basePollInterval {
			currentPollInterval = basePollInterval
			ticker.Stop()
			ticker = time.NewTicker(currentPollInterval)
		}
		unchangedCount = 0
		lastState = current

		if err := insertEcbState(dbmap, current); err != nil {
			logging.Logger().Errorf("failed to insert ecbstate %q: %v", current, err)
		}
	}
}

func readEcbStateString() string {
	layout := GetPinLayout()
	pass := readLevel(layout.Pass)
	fail := readLevel(layout.Fail)
	undertest := readLevel(layout.UnderTest)
	line := readLevel(layout.LineSelect)
	return fmt.Sprintf("%s.%s.%s.%s", undertest, pass, fail, line)
}

func insertEcbState(dbmap *gorp.DbMap, value string) error {
	now := time.Now()
	_, err := dbmap.Exec(`
		INSERT INTO ecbstates (tgl, ecbstate, readstate, created_at, updated_at)
		VALUES (DATE(?), ?, '', ?, ?)
		ON DUPLICATE KEY UPDATE ecbstate = VALUES(ecbstate), readstate = '', updated_at = VALUES(updated_at)
	`, now, value, now, now)
	return err
}

func logStateChange(oldState, newState string) {
	oldParts := strings.Split(oldState, ".")
	newParts := strings.Split(newState, ".")
	if len(oldParts) < 4 || len(newParts) < 4 {
		return
	}

	labels := []string{"UNDERTEST", "PASS", "FAIL", "LINE"}
	layout := GetPinLayout()
	pins := []string{layout.UnderTest, layout.Pass, layout.Fail, layout.LineSelect}

	for i := 0; i < 3; i++ {
		if oldParts[i] != newParts[i] {
			action := "RELEASED (HIGH)"
			if newParts[i] == "0" {
				action = "PRESSED (LOW)"
			}
			phys := WiringPiToPhysical(pins[i])
			logging.Logger().Infof("[GPIO Input] Button %s (Pin %s, Phys %s) %s", labels[i], pins[i], phys, action)
		}
	}

	if oldParts[3] != newParts[3] {
		logging.Logger().Infof("[GPIO Input] LINE SELECT changed to %s", newParts[3])
	}
}
