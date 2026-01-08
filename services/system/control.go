/*
   file:           services/system/control.go
   description:    Layanan backend untuk control
   created:        220711663@students.uajy.ac.id 04-11-2025
*/

package system

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"go-ecb/configs"
	"go-ecb/pkg/logging"
)

const (
	commandTimeout  = 2 * time.Minute
	commandFastWait = 2 * time.Second
)

func Shutdown() error {
	return runCommandFromEnv("SIMO_SHUTDOWN_COMMAND", "sudo shutdown -h now")
}

// Reboot adalah fungsi untuk reboot.
func Reboot() error {
	return runCommandFromEnv("SIMO_REBOOT_COMMAND", "sudo reboot")
}

// runCommandFromEnv adalah fungsi untuk menjalankan command from env.
func runCommandFromEnv(envKey, fallback string) error {
	if runtime.GOOS != "linux" {
		return fmt.Errorf("command not supported on %s", runtime.GOOS)
	}
	spec := configs.GetEnv(envKey, fallback)
	parts := strings.Fields(strings.TrimSpace(spec))
	if len(parts) == 0 {
		return fmt.Errorf("no command configured for %s", envKey)
	}

	logSystemEvent("%s: executing command '%s'", envKey, strings.Join(parts, " "))

	ctx, cancel := context.WithTimeout(context.Background(), commandTimeout)
	cmd := exec.CommandContext(ctx, parts[0], parts[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		cancel()
		logSystemEvent("%s: failed to start: %v", envKey, err)
		return err
	}

	waitCh := make(chan error, 1)
	go func() {
		defer cancel()
		waitCh <- cmd.Wait()
	}()

	select {
	case err := <-waitCh:
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				logSystemEvent("%s: timed out after %s", envKey, commandTimeout)
			} else {
				logSystemEvent("%s: command exited with error: %v", envKey, err)
			}
			return err
		}
		logSystemEvent("%s: command completed successfully", envKey)
	case <-time.After(commandFastWait):
		// Keep waiting in the background and log result; return immediately to avoid UI hang.
		go func() {
			if err := <-waitCh; err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					logSystemEvent("%s: timed out after %s", envKey, commandTimeout)
				} else {
					logSystemEvent("%s: command exited with error: %v", envKey, err)
				}
				return
			}
			logSystemEvent("%s: command completed successfully", envKey)
		}()
	}

	return nil
}

func logSystemEvent(format string, args ...interface{}) {
	logging.Logger().Infof(format, args...)

	dir := filepath.Join("storage", "logs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		logging.Logger().Errorf("[system] failed creating log dir: %v", err)
		return
	}

	line := fmt.Sprintf("%s %s\n", time.Now().Format(time.RFC3339), fmt.Sprintf(format, args...))
	logPath := filepath.Join(dir, "system.log")
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		logging.Logger().Errorf("[system] failed opening log file: %v", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(line); err != nil {
		logging.Logger().Errorf("[system] failed writing log append: %v", err)
	}
}
