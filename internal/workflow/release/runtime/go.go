package runtime

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type GoTarget struct {
	OS   string
	Arch string
}

var DefaultGoTargets = []GoTarget{
	{OS: "linux", Arch: "amd64"},
	{OS: "linux", Arch: "arm64"},
	{OS: "darwin", Arch: "amd64"},
	{OS: "darwin", Arch: "arm64"},
}

func DetectGo(root string) bool {
	_, err := os.Stat(filepath.Join(root, "go.mod"))
	return err == nil
}

func ExpectedGoAssetNames(version, binaryName string) []string {
	names := make([]string, 0, len(DefaultGoTargets)+1)
	for _, t := range DefaultGoTargets {
		names = append(names, fmt.Sprintf("%s_%s_%s_%s.tar.gz", binaryName, version, t.OS, t.Arch))
	}
	names = append(names, "checksums.txt")
	return names
}

func BuildGoReleaseArtifacts(version, binaryName, buildTarget string) ([]string, error) {
	if _, err := runLocal("go", "version"); err != nil {
		return nil, err
	}
	distDir := filepath.Join("dist", "release", version)
	if err := os.RemoveAll(distDir); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(distDir, 0o755); err != nil {
		return nil, err
	}
	assets := make([]string, 0, len(DefaultGoTargets)+1)
	checksumLines := make([]string, 0, len(DefaultGoTargets))
	for _, t := range DefaultGoTargets {
		archivePath, checksum, err := buildTargetArchive(distDir, version, binaryName, buildTarget, t)
		if err != nil {
			return nil, err
		}
		assets = append(assets, archivePath)
		checksumLines = append(checksumLines, checksum+"  "+filepath.Base(archivePath))
	}
	checksumPath := filepath.Join(distDir, "checksums.txt")
	if err := os.WriteFile(checksumPath, []byte(strings.Join(checksumLines, "\n")+"\n"), 0o644); err != nil {
		return nil, err
	}
	assets = append(assets, checksumPath)
	return assets, nil
}

func buildTargetArchive(distDir, version, binaryName, buildTarget string, target GoTarget) (string, string, error) {
	binPath := filepath.Join(distDir, fmt.Sprintf("%s_%s_%s_%s", binaryName, version, target.OS, target.Arch))
	if _, err := runLocalEnv(
		[]string{"GOOS=" + target.OS, "GOARCH=" + target.Arch, "CGO_ENABLED=0"},
		"go", "build", "-ldflags", "-X main.version="+version, "-o", binPath, buildTarget,
	); err != nil {
		return "", "", err
	}
	archivePath := binPath + ".tar.gz"
	if err := writeTarGz(archivePath, binPath, binaryName); err != nil {
		return "", "", err
	}
	if err := os.Remove(binPath); err != nil {
		return "", "", err
	}
	sum, err := sha256File(archivePath)
	if err != nil {
		return "", "", err
	}
	return archivePath, sum, nil
}

func writeTarGz(archivePath, sourcePath, archiveName string) error {
	in, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer in.Close()
	info, err := in.Stat()
	if err != nil {
		return err
	}
	out, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer out.Close()
	gz := gzip.NewWriter(out)
	defer gz.Close()
	tw := tar.NewWriter(gz)
	defer tw.Close()
	hdr := &tar.Header{
		Name: archiveName,
		Mode: 0o755,
		Size: info.Size(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if _, err := io.Copy(tw, in); err != nil {
		return err
	}
	return nil
}

func runLocal(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}

func runLocalEnv(env []string, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
