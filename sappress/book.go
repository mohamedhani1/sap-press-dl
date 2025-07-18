package sappress

import (
	"archive/zip"
	"os"
	"path/filepath"
)

// createEPUBStructure sets up the EPUB directory structure
func CreateEPUBStructure(tempDir string) error {
	// Create mimetype file
	if err := os.WriteFile(filepath.Join(tempDir, "mimetype"), []byte("application/epub+zip"), 0644); err != nil {
		return err
	}

	// Create META-INF directory and container.xml
	if err := os.Mkdir(filepath.Join(tempDir, "META-INF"), 0755); err != nil {
		return err
	}
	containerXML := `<?xml version="1.0" encoding="UTF-8"?>
<container version="1.0" xmlns="urn:oasis:names:tc:opendocument:xmlns:container">
   <rootfiles>
      <rootfile full-path="OEBPS/content.opf" media-type="application/oebps-package+xml"/>
   </rootfiles>
</container>`
	if err := os.WriteFile(filepath.Join(tempDir, "META-INF", "container.xml"), []byte(containerXML), 0644); err != nil {
		return err
	}

	// Create OEBPS directory
	return os.Mkdir(filepath.Join(tempDir, "OEBPS"), 0755)
}

// createEPUB zips the directory into an EPUB file
func CreateEPUB(tempDir, outputEPUB string) error {
	outFile, err := os.Create(outputEPUB)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zw := zip.NewWriter(outFile)
	defer zw.Close()

	// Add mimetype first, uncompressed
	mimetypePath := filepath.Join(tempDir, "mimetype")
		mimetypeHeader := &zip.FileHeader{
			Name:   "mimetype",
			Method: zip.Store, // No compression
		}
		f, err := zw.CreateHeader(mimetypeHeader)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(mimetypePath)
		if err != nil {
			return err
		}
		if _, err := f.Write(data); err != nil {
			return err
		}


	// Walk the directory and add other files
	return filepath.Walk(tempDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "mimetype" {
			return nil // Skip mimetype, already added
		}

		relPath, err := filepath.Rel(tempDir, filePath)
		relPath = filepath.ToSlash(relPath)
		if err != nil {
			return err
		}
		f, err := zw.Create(relPath)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		_, err = f.Write(data)
		return err
	})
}