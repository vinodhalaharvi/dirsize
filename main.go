package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"text/tabwriter"
)

type SizeFormatter func(int64) string

type DirSize struct {
	Path string
	Size int64
}

var sema = make(chan struct{}, 20) // Semaphore to limit concurrency

// BlockSize We don't use this currently, we should probably move up in sizes of 4K as that's the minimum block size
const BlockSize = 4096 // Standard block size on many filesystems

func main() {
	recursive := flag.Bool("recursive", false, "Calculate sizes recursively")
	human := flag.Bool("human", false, "Display sizes in human-readable format")
	flag.Parse()

	dirs := flag.Args()
	if len(dirs) == 0 {
		fmt.Println("Please provide at least one directory")
		os.Exit(1)
	}

	var formatter SizeFormatter
	if *human {
		formatter = formatHuman
	} else {
		formatter = formatBytes
	}

	results := make(chan DirSize)
	var wg sync.WaitGroup

	for _, dir := range dirs {
		wg.Add(1)
		go func(d string) {
			defer wg.Done()
			sema <- struct{}{} // Acquire token
			size := calculateSize(d, *recursive, formatter)
			<-sema // Release token
			results <- DirSize{d, size}
		}(dir)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	var total int64
	for result := range results {
		fmt.Fprintf(w, "%s\t%s\n", formatter(result.Size), result.Path)
		total += result.Size
	}

	fmt.Fprintf(w, "%s\t%s\n", formatter(total), "Total")
	w.Flush()
}

func calculateSize(dir string, recursive bool, formatter SizeFormatter) int64 {
	var size int64
	var mu sync.Mutex
	var wg sync.WaitGroup

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("Error accessing %s: %v\n", path, err)
			return nil // Continue walking despite the error
		}

		mu.Lock()
		if !info.IsDir() {
			// Yes, in the current version we don't count directories for size calculation
			// du probably adds 4K for each directory to it's size
			size += info.Size()
		}
		mu.Unlock()

		if path != dir && recursive {
			wg.Add(1)
			go func(p string, s int64) {
				defer wg.Done()
				fmt.Printf("%s: %s\n", p, formatter(s))
			}(path, info.Size())
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", dir, err)
	}

	wg.Wait()
	return size
}

func formatBytes(size int64) string {
	return fmt.Sprintf("%d", size/1024) + "K"
}

func formatHuman(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%c", float64(size)/float64(div), "KMGTPE"[exp])
}
