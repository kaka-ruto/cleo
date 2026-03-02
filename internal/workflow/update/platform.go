package update

import (
	"fmt"
	"runtime"
)

func assetName(version string) (string, error) {
	arch, err := releaseArch()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "linux", "darwin":
		return fmt.Sprintf("cleo_%s_%s_%s.tar.gz", version, runtime.GOOS, arch), nil
	default:
		return "", fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}
}

func releaseArch() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "amd64", nil
	case "arm64":
		return "arm64", nil
	default:
		return "", fmt.Errorf("unsupported arch: %s", runtime.GOARCH)
	}
}
