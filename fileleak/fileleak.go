//go:build !solution

package fileleak

import (
	"fmt"
	"os"
	"path/filepath"
)

type testingT interface {
	Errorf(msg string, args ...interface{})
	Cleanup(func())
}

func VerifyNone(t testingT) {
	initialFiles := getOpenFiles()
	t.Cleanup(func() {
		finalFiles := getOpenFiles()
		leakedFiles := findLeaks(initialFiles, finalFiles)
		if len(leakedFiles) > 0 {
			t.Errorf("fileleak detected leaked file descriptors: %v", leakedFiles)
		}
	})
}

func getOpenFiles() map[string]struct{} {
	openFiles := make(map[string]struct{})
	fdDir := "/proc/self/fd"
	entries, err := os.ReadDir(fdDir)
	if err != nil {
		return openFiles
	}

	for _, entry := range entries {
		fdPath := filepath.Join(fdDir, entry.Name())
		target, err := os.Readlink(fdPath)
		if err == nil && shouldTrack(target) {
			openFiles[fmt.Sprintf("%s -> %s", entry.Name(), target)] = struct{}{}
		}
	}
	return openFiles
}

func shouldTrack(target string) bool {
	return target != "anon_inode:[eventfd]" && target != "anon_inode:[signalfd]"
}

func findLeaks(initial, final map[string]struct{}) []string {
	var leaks []string
	for fd := range final {
		if _, exists := initial[fd]; !exists {
			leaks = append(leaks, fd)
		}
	}
	return leaks
}
