// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program that performs a series of I/O related tasks to
// better understand tracing in Go.
package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"runtime/trace"
	"strings"
	"sync"
	"sync/atomic"
)

type (
	item struct {
		XMLName     xml.Name `xml:"item"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
	}

	channel struct {
		XMLName xml.Name `xml:"channel"`
		Items   []item   `xml:"item"`
	}

	document struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)

func main() {
	// pprof.StartCPUProfile(os.Stdout)
	// defer pprof.StopCPUProfile()

	trace.Start(os.Stdout)
	defer trace.Stop()

	docs := make([]string, 4000)
	for i := range docs {
		docs[i] = fmt.Sprintf("newsfeed-%.4d.xml", i)
	}

	topic := "president"
	n := freq(topic, docs)
	// n := freqConcurrent(topic, docs)
	// n := freqConcurrentSem(topic, docs)
	// n := freqNumCPU(topic, docs)
	// n := freqNumCPUTasks(topic, docs)
	// n := freqActor(topic, docs)

	log.Printf("Searching %d files, found %s %d times.", len(docs), topic, n)
}

func freq(topic string, docs []string) int {
	var found int32

	g := len(docs)
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for i := 0; i < g; i++ {
		go func() {
			var lfound int32
			defer func() {
				atomic.AddInt32(&found, lfound)
				wg.Done()
			}()

			for doc := range ch {
				file := fmt.Sprintf("%s.xml", doc[:8])
				f, err := os.OpenFile(file, os.O_RDONLY, 0)
				if err != nil {
					log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
					return
				}

				data, err := io.ReadAll(f)
				if err != nil {
					log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
					return
				}

				var d document
				if err := xml.Unmarshal(data, &d); err != nil {
					log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
					return
				}
				f.Close()
				for _, item := range d.Channel.Items {
					if strings.Contains(item.Title, topic) {
						lfound++
						continue
					}

					if strings.Contains(item.Description, topic) {
						lfound++
					}
				}
			}
		}()
	}

	for _, doc := range docs {
		ch <- doc
	}
	close(ch)

	wg.Wait()
	return int(found)
}

func freqConcurrent(topic string, docs []string) int {
	var found int32

	g := len(docs)
	var wg sync.WaitGroup
	wg.Add(g)

	for _, doc := range docs {
		go func(doc string) {
			var lFound int32
			defer func() {
				atomic.AddInt32(&found, lFound)
				wg.Done()
			}()

			file := fmt.Sprintf("%s.xml", doc[:8])
			f, err := os.OpenFile(file, os.O_RDONLY, 0)
			if err != nil {
				log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
				return
			}
			defer f.Close()

			data, err := io.ReadAll(f)
			if err != nil {
				log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
				return
			}

			var d document
			if err := xml.Unmarshal(data, &d); err != nil {
				log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
				return
			}

			for _, item := range d.Channel.Items {
				if strings.Contains(item.Title, topic) {
					lFound++
					continue
				}

				if strings.Contains(item.Description, topic) {
					lFound++
				}
			}
		}(doc)
	}

	wg.Wait()
	return int(found)
}

func freqConcurrentSem(topic string, docs []string) int {
	var found int32

	g := len(docs)
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan bool, runtime.NumCPU())

	for _, doc := range docs {
		go func(doc string) {
			ch <- true
			{
				var lFound int32
				defer func() {
					atomic.AddInt32(&found, lFound)
					wg.Done()
				}()

				file := fmt.Sprintf("%s.xml", doc[:8])
				f, err := os.OpenFile(file, os.O_RDONLY, 0)
				if err != nil {
					log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
					return
				}
				defer f.Close()

				data, err := io.ReadAll(f)
				if err != nil {
					log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
					return
				}

				var d document
				if err := xml.Unmarshal(data, &d); err != nil {
					log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
					return
				}

				for _, item := range d.Channel.Items {
					if strings.Contains(item.Title, topic) {
						lFound++
						continue
					}

					if strings.Contains(item.Description, topic) {
						lFound++
					}
				}
			}
			<-ch
		}(doc)
	}

	wg.Wait()
	return int(found)
}

func freqNumCPU(topic string, docs []string) int {
	var found int32

	g := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for i := 0; i < g; i++ {
		go func() {
			var lFound int32
			defer func() {
				atomic.AddInt32(&found, lFound)
				wg.Done()
			}()

			for doc := range ch {
				file := fmt.Sprintf("%s.xml", doc[:8])
				f, err := os.OpenFile(file, os.O_RDONLY, 0)
				if err != nil {
					log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
					return
				}

				data, err := io.ReadAll(f)
				if err != nil {
					f.Close()
					log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
					return
				}
				f.Close()

				var d document
				if err := xml.Unmarshal(data, &d); err != nil {
					log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
					return
				}

				for _, item := range d.Channel.Items {
					if strings.Contains(item.Title, topic) {
						lFound++
						continue
					}

					if strings.Contains(item.Description, topic) {
						lFound++
					}
				}
			}
		}()
	}

	for _, doc := range docs {
		ch <- doc
	}
	close(ch)

	wg.Wait()
	return int(found)
}

func freqNumCPUTasks(topic string, docs []string) int {
	var found int32

	g := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for i := 0; i < g; i++ {
		go func() {
			var lFound int32
			defer func() {
				atomic.AddInt32(&found, lFound)
				wg.Done()
			}()

			for doc := range ch {
				func() {
					file := fmt.Sprintf("%s.xml", doc[:8])
					ctx, tt := trace.NewTask(context.Background(), doc)
					defer tt.End()

					reg := trace.StartRegion(ctx, "OpenFile")
					f, err := os.OpenFile(file, os.O_RDONLY, 0)
					if err != nil {
						log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
						return
					}
					reg.End()

					reg = trace.StartRegion(ctx, "ReadAll")
					data, err := io.ReadAll(f)
					if err != nil {
						f.Close()
						log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
						return
					}
					f.Close()
					reg.End()

					reg = trace.StartRegion(ctx, "Unmarshal")
					var d document
					if err := xml.Unmarshal(data, &d); err != nil {
						log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
						return
					}
					reg.End()

					reg = trace.StartRegion(ctx, "Search")
					for _, item := range d.Channel.Items {
						if strings.Contains(item.Title, topic) {
							lFound++
							continue
						}

						if strings.Contains(item.Description, topic) {
							lFound++
						}
					}
					reg.End()
				}()
			}
		}()
	}

	for _, doc := range docs {
		ch <- doc
	}
	close(ch)

	wg.Wait()
	return int(found)
}

func freqActor(topic string, docs []string) int {
	files := make(chan *os.File, 100)
	go func() {
		for _, doc := range docs {
			file := fmt.Sprintf("%s.xml", doc[:8])
			f, err := os.OpenFile(file, os.O_RDONLY, 0)
			if err != nil {
				log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
				break
			}
			files <- f
		}
		close(files)
	}()

	data := make(chan []byte, 100)
	go func() {
		for f := range files {
			defer f.Close()
			d, err := io.ReadAll(f)
			if err != nil {
				log.Printf("Reading Document [%s] : ERROR : %v", f.Name(), err)
				break
			}
			data <- d
		}
		close(data)
	}()

	rss := make(chan document, 100)
	go func() {
		for dt := range data {
			var d document
			if err := xml.Unmarshal(dt, &d); err != nil {
				log.Printf("Decoding Document : ERROR : %v", err)
				break
			}
			rss <- d
		}
		close(rss)
	}()

	var wg sync.WaitGroup
	wg.Add(1)
	var found int
	go func() {
		for d := range rss {
			for _, item := range d.Channel.Items {
				if strings.Contains(item.Title, topic) {
					found++
					continue
				}

				if strings.Contains(item.Description, topic) {
					found++
				}
			}
		}
		wg.Done()
	}()

	wg.Wait()
	return found
}
