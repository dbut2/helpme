# autoupgrade

Use `autoupgrade` to install the latest version of your binary, quietly in the background. Zero dependencies.

## Features

- **Zero dependencies** - Uses only Go standard library
- **Context support** - Full cancellation support for upgrade operations
- **Background upgrades** - Non-blocking upgrade operations via goroutines
- **Build info access** - Get build information from both current and upgraded binaries
- **Smart skipping** - Automatically skips upgrades for development builds

## Installation

```bash
go get github.com/melt-inc/autoupgrade
```

## Usage

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/melt-inc/autoupgrade"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Upgrade the current binary
    result := autoupgrade.Upgrade(ctx, "cmd/myapp")

    if result.ExitError != nil {
        fmt.Printf("Upgrade failed: %v\n", result.ExitError)
        return
    }

    if result.DidUpgrade() {
        newInfo, _ := result.NewBuildInfo()
        fmt.Printf("Upgraded from %s to %s\n",
            result.CurrentInfo.Main.Version,
            newInfo.Main.Version)
    } else {
        fmt.Println("No upgrade needed")
    }
}
```

### Background Upgrades

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/melt-inc/autoupgrade"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Start upgrade in background
    resultCh := autoupgrade.UpgradeBackground(ctx, "cmd/myapp")

    // Do other work while upgrade happens
    fmt.Println("Upgrade started in background...")

    // Wait for result
    result := <-resultCh

    if result.ExitError != nil {
        fmt.Printf("Background upgrade failed: %v\n", result.ExitError)
        return
    }

    if result.DidUpgrade() {
        fmt.Println("Background upgrade completed successfully!")
    }
}
```

### Checking Upgrade Status

```go
result := autoupgrade.Upgrade(ctx, "")

// Check if upgrade actually occurred
if result.DidUpgrade() {
    // Get new build info
    newInfo, err := result.NewBuildInfo()
    if err == nil {
        fmt.Printf("New version: %s\n", newInfo.Main.Version)
    }
}

// Check for errors
if result.ExitError != nil {
    fmt.Printf("Upgrade error: %v\n", result.ExitError)
}

// Access current build info
if result.CurrentInfo != nil {
    fmt.Printf("Current version: %s\n", result.CurrentInfo.Main.Version)
}
```

## API Reference

### Types

#### `UpgradeResult`

Contains the result of an upgrade operation.

```go
type UpgradeResult struct {
    CurrentInfo *debug.BuildInfo // Current build info of running process
    ExitError   error            // Error from upgrade process, if any
}
```

### Functions

#### `Upgrade(ctx context.Context, packagePath string) *UpgradeResult`

Attempts to upgrade the current binary to the latest version using `go install`. The upgrade is skipped if the current version is a development build or build info is unavailable.

- `ctx`: Context for cancellation support
- `packagePath`: Relative path from module root to package (use `""` for root)

#### `UpgradeBackground(ctx context.Context, packagePath string) <-chan *UpgradeResult`

Runs `Upgrade` in a goroutine and returns a channel that receives the result. Supports context cancellation.

### Methods

#### `(u *UpgradeResult) DidUpgrade() bool`

Returns `true` if the upgrade actually occurred. Returns `false` if:
- Build information is not available
- Current version is a development build
- Upgrade was not necessary (already at latest version)
- Upgrade failed

#### `(u *UpgradeResult) NewBuildInfo() (*debug.BuildInfo, error)`

Returns the build information of the newly installed binary. The result is cached using `sync.Once`, so the file is only read once per `UpgradeResult`.

## How It Works

1. Reads current build information using `runtime/debug.ReadBuildInfo()`
2. Skips upgrade if current version is `"(devel)"` or build info unavailable
3. Runs `go install module/path/package@latest` to install latest version
4. Provides access to new build information via lazy-loaded `NewBuildInfo()` method

## License

MIT License
