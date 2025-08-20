package autoupgrade

import (
	"context"
	"debug/buildinfo"
	"os"
	"os/exec"
	"path"
	"runtime/debug"
	"sync"
)

// UpgradeResult contains the result of an upgrade operation.
// It provides access to build information from both the current process
// and the newly installed binary.
type UpgradeResult struct {
	CurrentInfo *debug.BuildInfo // Current build information of the running process, if available
	ExitError   error            // Error encountered during the upgrade process, if any
	once        sync.Once
	newInfo     *debug.BuildInfo
	newInfoErr  error
}

// Upgrade attempts to upgrade the current binary to the latest version using
// 'go install'. The packagePath parameter specifies the relative path from the
// module root to the package. Upgrade is skipped if the current version is a
// development build or build info is unavailable.
// Context cancellation can be used to kill the go install process.
func Upgrade(ctx context.Context, packagePath string) *UpgradeResult {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return &UpgradeResult{}
	}

	// Don't upgrade if the current version is a development version
	if info.Main.Version == "(devel)" {
		return &UpgradeResult{CurrentInfo: info}
	}

	modulePath := info.Main.Path
	if modulePath == "" {
		return &UpgradeResult{CurrentInfo: info}
	}

	cmd := exec.CommandContext(ctx, "go", "install", fullPath(modulePath, packagePath, "latest"))
	// Suppress standard output and error
	cmd.Stdout = nil
	cmd.Stderr = nil

	err := cmd.Run()
	return &UpgradeResult{
		CurrentInfo: info,
		ExitError:   err,
	}
}

// UpgradeBackground runs Upgrade in a goroutine and returns a channel that will
// receive the UpgradeResult. The channel is closed after the result is sent.
// This allows for non-blocking upgrade operations. The context can be used to
// cancel the upgrade operation.
func UpgradeBackground(ctx context.Context, packagePath string) <-chan *UpgradeResult {
	ch := make(chan *UpgradeResult, 1)
	go func() {
		defer close(ch)
		select {
		case <-ctx.Done():
			ch <- &UpgradeResult{
				ExitError: ctx.Err(),
			}
		case ch <- Upgrade(ctx, packagePath):
		}
	}()
	return ch
}

// DidUpgrade returns false if the upgrade did not occur, this can happen when
// the build information is not available, the current version is a development
// version, or upgrade was not necessary (e.g., already at latest version).
func (u *UpgradeResult) DidUpgrade() bool {
	if u.CurrentInfo == nil {
		return false
	}
	if u.CurrentInfo.Main.Version == "(devel)" {
		return false
	}
	newInfo, _ := u.NewBuildInfo()
	return newInfo != nil && newInfo.Main.Version != u.CurrentInfo.Main.Version
}

// NewBuildInfo returns the build information of the newly installed binary.
// Returns nil if the executable path cannot be determined or the build info
// cannot be read.
func (u *UpgradeResult) NewBuildInfo() (*debug.BuildInfo, error) {
	u.once.Do(func() {
		execPath, err := os.Executable()
		if err != nil {
			u.newInfoErr = err
			return
		}
		u.newInfo, u.newInfoErr = buildinfo.ReadFile(execPath)
	})
	return u.newInfo, u.newInfoErr
}

// fullPath constructs the full module path with version for 'go install'.
// It combines the module path, package path, and version into the format
// expected by go install.
func fullPath(modulePath, packagePath, version string) string {
	ret := path.Join(modulePath, packagePath+"@"+version)
	return ret
}
