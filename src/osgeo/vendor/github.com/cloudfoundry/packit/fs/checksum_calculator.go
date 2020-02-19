package fs

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

type ChecksumCalculator struct{}

func NewChecksumCalculator() ChecksumCalculator {
	return ChecksumCalculator{}
}

type calculatedFile struct {
	path     string
	checksum []byte
	err      error
}

func (c ChecksumCalculator) Sum(path string) (string, error) {
	info, err := os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	if !info.IsDir() {
		file, err := os.Open(path)
		if err != nil {
			return "", fmt.Errorf("failed to calculate checksum: %w", err)
		}
		defer file.Close()

		hash := sha256.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			return "", fmt.Errorf("failed to calculate checksum: %w", err)
		}

		return hex.EncodeToString(hash.Sum(nil)), nil
	}

	//Finds all files in directoy
	var filesFromDir []string
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Mode().IsRegular() {
			filesFromDir = append(filesFromDir, path)
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to calculate checksum: %w", err)
	}

	//Gather all checksums into one byte array and check for checksum calculation errors
	hash := sha256.New()
	for _, f := range getParallelChecksums(filesFromDir) {
		if f.err != nil {
			return "", fmt.Errorf("failed to calculate checksum: %w", f.err)
		}

		_, err := hash.Write(f.checksum)
		if err != nil {
			return "", fmt.Errorf("failed to calculate checksum: %w", err)
		}
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func getParallelChecksums(filesFromDir []string) []calculatedFile {
	var checksumResults []calculatedFile
	numFiles := len(filesFromDir)
	files := make(chan string, numFiles)
	calculatedFiles := make(chan calculatedFile, numFiles)

	//Spawns workers
	for i := 0; i < runtime.NumCPU(); i++ {
		go fileChecksumer(files, calculatedFiles)
	}

	//Puts files in worker queue
	for _, f := range filesFromDir {
		files <- f
	}

	close(files)

	//Pull all calculated files off of result queue
	for i := 0; i < numFiles; i++ {
		checksumResults = append(checksumResults, <-calculatedFiles)
	}

	//Sort calculated files for consistent checksuming
	sort.Slice(checksumResults, func(i, j int) bool {
		return checksumResults[i].path < checksumResults[j].path
	})

	return checksumResults
}

func fileChecksumer(files chan string, calculatedFiles chan calculatedFile) {
	for path := range files {
		result := calculatedFile{path: path}

		file, err := os.Open(path)
		if err != nil {
			result.err = err
			calculatedFiles <- result
			continue
		}

		hash := sha256.New()
		_, err = io.Copy(hash, file)
		if err != nil {
			result.err = err
			calculatedFiles <- result
			continue
		}

		if err := file.Close(); err != nil {
			result.err = err
			calculatedFiles <- result
			continue
		}

		result.checksum = hash.Sum(nil)
		calculatedFiles <- result
	}
}
