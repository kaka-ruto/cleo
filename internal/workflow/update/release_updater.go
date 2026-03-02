package update

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	defaultOwner = "cafaye"
	defaultRepo  = "cleo"
)

type ReleaseUpdater struct {
	owner  string
	repo   string
	client *http.Client
}

func NewReleaseUpdater() *ReleaseUpdater {
	return &ReleaseUpdater{
		owner: defaultOwner,
		repo:  defaultRepo,
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (u *ReleaseUpdater) UpdateLatest() error {
	rel, err := u.latestRelease()
	if err != nil {
		return err
	}
	asset, err := assetName(rel.TagName)
	if err != nil {
		return err
	}
	binaryURL, err := findAssetURL(rel.Assets, asset)
	if err != nil {
		return err
	}
	checksumsURL, err := findAssetURL(rel.Assets, "checksums.txt")
	if err != nil {
		return err
	}
	tmpDir, err := os.MkdirTemp("", "cleo-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	archivePath := filepath.Join(tmpDir, asset)
	if err := u.download(binaryURL, archivePath); err != nil {
		return err
	}
	checksumsPath := filepath.Join(tmpDir, "checksums.txt")
	if err := u.download(checksumsURL, checksumsPath); err != nil {
		return err
	}
	if err := verifyChecksum(archivePath, checksumsPath, asset); err != nil {
		return err
	}
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	newBin := filepath.Join(tmpDir, "cleo")
	if err := extractBinary(archivePath, newBin); err != nil {
		return err
	}
	if err := os.Chmod(newBin, 0o755); err != nil {
		return err
	}
	if err := os.Rename(newBin, exePath); err != nil {
		return fmt.Errorf("replace %s: %w", exePath, err)
	}
	return nil
}

func (u *ReleaseUpdater) latestRelease() (*githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", u.owner, u.repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return nil, fmt.Errorf("github releases api status=%d body=%s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return nil, err
	}
	if rel.TagName == "" {
		return nil, fmt.Errorf("latest release has empty tag")
	}
	return &rel, nil
}

func (u *ReleaseUpdater) download(url, path string) error {
	resp, err := u.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed status=%d url=%s", resp.StatusCode, url)
	}
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func extractBinary(archivePath, outPath string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()
	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		h, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if h.Typeflag != tar.TypeReg {
			continue
		}
		if !strings.Contains(h.Name, "cleo") {
			continue
		}
		out, err := os.Create(outPath)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, tr)
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	}
	return fmt.Errorf("cleo binary not found in archive")
}

func verifyChecksum(archivePath, checksumsPath, asset string) error {
	checksums, err := os.ReadFile(checksumsPath)
	if err != nil {
		return err
	}
	expected, err := checksumForAsset(string(checksums), asset)
	if err != nil {
		return err
	}
	got, err := fileSHA256(archivePath)
	if err != nil {
		return err
	}
	if !strings.EqualFold(expected, got) {
		return fmt.Errorf("checksum mismatch for %s", asset)
	}
	return nil
}

func checksumForAsset(content, asset string) (string, error) {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		parts := strings.Fields(strings.TrimSpace(line))
		if len(parts) < 2 {
			continue
		}
		if parts[1] == asset {
			return parts[0], nil
		}
	}
	return "", fmt.Errorf("checksum not found for %s", asset)
}

func fileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func findAssetURL(assets []githubAsset, name string) (string, error) {
	for _, a := range assets {
		if strings.TrimSpace(a.Name) == name {
			return a.BrowserDownloadURL, nil
		}
	}
	return "", fmt.Errorf("release asset not found: %s", name)
}

type githubRelease struct {
	TagName string        `json:"tag_name"`
	Assets  []githubAsset `json:"assets"`
}

type githubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}
