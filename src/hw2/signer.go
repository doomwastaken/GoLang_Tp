package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}

	in := make(chan interface{})
	defer close(in)
	for _, jb := range jobs {
		out := make(chan interface{})

		wg.Add(1)
		go func(in, out chan interface{}, jb job) {
			jb(in, out)
			wg.Done()
			close(out)
		}(in, out, jb)
		in = out
	}

	wg.Wait()
}

func SingleHash(in, out chan interface{}) {
	wgOut := &sync.WaitGroup{}
	md5Mutex := &sync.Mutex{}
	for i := range in {
		wgOut.Add(1)
		go func(i interface{}) {
			defer wgOut.Done()
			convertedData := strconv.Itoa(i.(int))
			wgIn := &sync.WaitGroup{}
			wgIn.Add(2)
			var lHash, rHash string
			
			go func() {
				defer wgIn.Done()
				lHash = DataSignerCrc32(convertedData)
			}()
			go func() {
				defer wgIn.Done()
				md5Mutex.Lock()
				md5 := DataSignerMd5(convertedData)
				md5Mutex.Unlock()
				rHash = DataSignerCrc32(md5)
				
			}()
			wgIn.Wait()
			out <- lHash + "~" + rHash
		}(i)
	}
	wgOut.Wait()
}

func MultiHash(in, out chan interface{}) {
	wgOut := &sync.WaitGroup{}

	for i := range in {
		wgOut.Add(1)
		data := i.(string)
		go func() {
			defer wgOut.Done()
			mutex := &sync.Mutex{}
			wgIn := &sync.WaitGroup{}
			hashs := make([]string, 6)
			for th := 0; th < 6; th++ {
				wgIn.Add(1)
				go func(th int) {
					defer wgIn.Done()
					hashs[th] = DataSignerCrc32(strconv.Itoa(th) + data)
				}(th)
			}
			wgIn.Wait()
			mutex.Lock()
			out <- strings.Join(hashs, "")
			mutex.Unlock()
		}()
	}
	wgOut.Wait()
}

func CombineResults(in, out chan interface{}) {
	var results []string
	for i := range in {
		results = append(results, i.(string))
	}
	sort.Strings(results)
	out <- strings.Join(results, "_")
}