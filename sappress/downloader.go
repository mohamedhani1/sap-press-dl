package sappress

import (
	
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	//"net/url"
	"os"
	"path/filepath"
	"sync"
	"log"
)

const (
    BaseURL       = "https://eba.sap-press.com/v1"
    LoginURL      = BaseURL + "/account/token"
    DownloadURL   = BaseURL + "/ebooks/%s/download?app_version=0&file_path=%s"
	AccountListURL= BaseURL + "/account/lists/?page_size=1000"

)

// OPF structures to parse the .opf file
type Package struct {
	XMLName  xml.Name `xml:"package"`
	Metadata Metadata `xml:"metadata"`
	Manifest Manifest `xml:"manifest"`
	Spine    Spine    `xml:"spine"`
}

type Metadata struct {
	Title       string `xml:"title"`
	Creator     string `xml:"creator"`
	Description string `xml:"description"`
}

type Manifest struct {
	Items []Item `xml:"item"`
}

type Item struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

type Spine struct {
	ItemRefs []ItemRef `xml:"itemref"`
}

type ItemRef struct {
	IDRef string `xml:"idref,attr"`
}

type Downloader struct {
	Config	*Config
	HttpClient	*http.Client
	Threads		int
}

func (d *Downloader) Download(bookID string) {
	log.SetFlags(0) // disable date/time, etc.

	// Input OPF URL
	log.Println("[*] Building the opf url...")
	opfURL := fmt.Sprintf(DownloadURL, bookID, "content.opf")
	

	// Fetch and parse the OPF file
	log.Println("[*] Fetching the opf content...")
	opfData, err := d.fetchFile(opfURL)
	if err != nil {
		log.Printf("[-] Error fetching OPF: %v\n", err)
		return
	}
	// parsing the opf to package
	log.Println("[*] Parsing the opf...")
	var pkg Package
	if err := xml.Unmarshal(opfData, &pkg); err != nil {
		log.Printf("[-] Error parsing OPF: %v\n", err)
		return
	}

	// Create a temporary directory for EPUB structure
	os.MkdirAll("tmp", 0755)
	tempDir, err := os.MkdirTemp("tmp", "")
	if err != nil {
		log.Printf("[-] Error creating temp dir: %v\n", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// Create EPUB directory structure
	if err := CreateEPUBStructure(tempDir); err != nil {
		log.Printf("Error creating EPUB structure: %v\n", err)
		return
	}

	

	var wg sync.WaitGroup
	sem := make(chan struct{}, d.Threads)
	lenOfDownloadedFiles := len(pkg.Manifest.Items)
	errors := make(chan error, lenOfDownloadedFiles)
	outputEPUB := CleanFilename(pkg.Metadata.Title + ".epub")
	progressBar := NewProgressBar(lenOfDownloadedFiles, 30, outputEPUB)

	log.Printf("Ready to download\n\tNumber of files: %d\n\tNumber of threads: %d\n\tDecrypting on flow\n\n", lenOfDownloadedFiles, d.Threads)

	for _, item := range pkg.Manifest.Items {
		fileURL := fmt.Sprintf(DownloadURL, bookID, item.Href)
		outputPath := filepath.Join(tempDir, "OEBPS", item.Href)
		wg.Add(1)

		go func(fileURL, outputPath string) {
			sem <- struct{}{} // acquire
			defer func() { <-sem }() // release
			defer wg.Done()

			//fmt.Printf("Downloading %s\n", fileURL)
			progressBar.Add()

			if err := d.downloadFile(fileURL, outputPath); err != nil {
				fmt.Printf("Error downloading %s: %v\n", fileURL, err)
				errors <- err
				return
			}
			
		}(fileURL, outputPath)
		
	}

	wg.Wait()
	close(errors)

	// Print any errors
	for err := range errors {
		log.Println("Error:", err)
	}

	// Save the OPF file locally
	opfPath := filepath.Join(tempDir, "OEBPS", "content.opf")
	if err := os.WriteFile(opfPath, opfData, 0644); err != nil {
		log.Printf("Error saving OPF file: %v\n", err)
		return
	}

	// Create the EPUB file
	log.Println("[*] Creating epub file...")
	if err := CreateEPUB(tempDir, outputEPUB); err != nil {
		log.Printf("Error creating EPUB: %v\n", err)
		return
	}

	log.Printf("EPUB created successfully: %s\n", outputEPUB)
}

// fetchFile downloads a file from a URL and returns its contents
func (d *Downloader) fetchFile(fileURL string) ([]byte, error) {
	resp, err := d.HttpClient.Get(fileURL)
	if err != nil {
		log.Fatal("API request failed:", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// check the content-key for decryption
	contentKey := resp.Header.Get("X-CONTENT-KEY")
	if contentKey != ""{
		decryptedBytes, err := DecryptIt(data, d.Config.UserKey[:40], contentKey)
		return decryptedBytes, err
	}
	return data, nil
}

// downloadFile downloads a file and saves it to the specified path
func (d *Downloader) downloadFile(fileURL, destPath string) error {
	data, err := d.fetchFile(fileURL)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(destPath, data, 0644)
}

