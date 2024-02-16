package updater_client

import (
	"archive/zip"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/inconshreveable/go-update"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

type UpdateFileMetadata struct {
	Updated string `json:"updated"`
}

type AutoUpdateConfiguration struct {
	ExpectedPublicKey      string
	UpdateLocationMetadata string
	UpdateLocationFile     string
	LocalVersionInfo       string
	ProductToken           string
}

var configuration *AutoUpdateConfiguration

func getUpdaterRequest(url string, moduleName string) (*http.Request, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf(url, moduleName), nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request for %s: %w", url, err)
	}

	if len(configuration.ProductToken) > 0 {
		req.Header.Set("Authorization", configuration.ProductToken)
	}

	return req, nil
}

func getHttpClient() *http.Client {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			for _, rawCert := range rawCerts {
				cert, err := x509.ParseCertificate(rawCert)
				if err != nil {
					return fmt.Errorf("failed to parse certificate: %v", err)
				}

				if len(configuration.ExpectedPublicKey) <= 1 {
					return nil
				}

				if hex.EncodeToString(cert.Signature) == configuration.ExpectedPublicKey {
					return nil
				}
			}

			return fmt.Errorf("certificate validation failed")
		},
	}

	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}
}

func relaunch() error {
	executablePath := os.Args[0]

	cmd := exec.Command(executablePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("unable to start new instance: %w", err)
	}

	os.Exit(0)

	return nil
}

func getFileMetadata(urlToFile string, moduleName string) (*UpdateFileMetadata, error) {
	httpRequest, err := getUpdaterRequest(urlToFile, moduleName)
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	client := getHttpClient()
	resp, err := client.Do(httpRequest)

	if err != nil {
		return nil, fmt.Errorf("unable to fetch metadata for %s: %w", urlToFile, err)
	}
	defer resp.Body.Close()

	tokenBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch metadata for %s: %w", urlToFile, err)
	}

	var downloadData UpdateFileMetadata
	err = json.Unmarshal(tokenBody, &downloadData)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch metadata for %s: %w", urlToFile, err)
	}

	return &downloadData, nil
}

func isShouldDownloadFile(metadataFromCloud *UpdateFileMetadata, localSnapshotInfo string) bool {
	currentMetadataFile, err := os.ReadFile(localSnapshotInfo)
	if err != nil {
		return true
	}

	var currentMetadata UpdateFileMetadata
	err = json.Unmarshal(currentMetadataFile, &currentMetadata)
	if err != nil {
		return true
	}
	fmt.Printf("Hi, my version is: %s \n", currentMetadata.Updated)

	cloudUpdateTime, err := time.Parse(time.RFC3339, metadataFromCloud.Updated)
	if err != nil {
		return true
	}
	localUpdateTime, err := time.Parse(time.RFC3339, currentMetadata.Updated)
	if err != nil {
		return true
	}

	if cloudUpdateTime.After(localUpdateTime) {
		return true
	}

	return false
}

func downloadByUrl(downloadUrl string, moduleName string) ([]byte, error) {
	httpRequest, err := getUpdaterRequest(downloadUrl, moduleName)
	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}

	client := getHttpClient()
	resp, err := client.Do(httpRequest)

	if err != nil {
		return nil, fmt.Errorf("update failed: %w", err)
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	return result, nil
}

func ConfigureAutoUpdateModule(config AutoUpdateConfiguration) {
	configuration = &config
}

func DoUpdate(moduleName string) error {
	downloadData, err := getFileMetadata(configuration.UpdateLocationFile, moduleName)
	if err != nil || downloadData == nil {
		return fmt.Errorf("update failed: %w", err)
	}

	if !isShouldDownloadFile(downloadData, configuration.LocalVersionInfo) {
		return nil
	}

	zipBytes, err := downloadByUrl(configuration.UpdateLocationFile, moduleName)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		return fmt.Errorf("update failed: %w, response is: %s", err, hex.EncodeToString(zipBytes))
	}

	for _, file := range zipReader.File {
		if file.Name == "update.exe" {
			fileReader, err := file.Open()
			if err != nil {
				return fmt.Errorf("error opening file inside zip: %w", err)
			}

			err = update.Apply(fileReader, update.Options{})
			if err != nil {
				if rerr := update.RollbackError(err); rerr != nil {
					fmt.Printf("Failed to rollback from failed update: %v\n", rerr)
				}
				return fmt.Errorf("error applying update: %w", err)
			}

			dataToSave, err := json.Marshal(downloadData)
			err = os.WriteFile(configuration.LocalVersionInfo, dataToSave, 0644)
			if err != nil {
				return fmt.Errorf("saving update result failed: %w", err)
			}

			fileReader.Close()

			return relaunch()
		}
	}

	return fmt.Errorf("binary not found in the zip file")
}
