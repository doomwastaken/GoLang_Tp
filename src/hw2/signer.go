package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

const (
	thCount = 6
)

func ExecutePipeline(jobs ...job) {
	wg := &sync.WaitGroup{}

	in := make(chan interface{})
	defer close(in)
	for _, jb := range jobs {
		out := make(chan interface{})

		wg.Add(1)
		go func(in, out chan interface{}, jb job) {
			defer wg.Done()
			jb(in, out)
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
			data := strconv.Itoa(i.(int))
			wgIn := &sync.WaitGroup{}
			wgIn.Add(2)
			var lHash string
			var rHash string
			
			go func() {
				defer wgIn.Done()
				lHash = DataSignerCrc32(data)
			}()
			go func() {
				defer wgIn.Done()
				md5Mutex.Lock()
				md5 := DataSignerMd5(data)
				md5Mutex.Unlock()
				rHash = DataSignerCrc32(md5)
				
			}()
			wgIn.Wait()
			hash:= lHash + "~" + rHash
			out <- hash
		}(i)
	}
	wgOut.Wait()
}

func MultiHash(in, out chan interface{}) {
	wgOut := &sync.WaitGroup{}

	for i := range in {
		wgOut.Add(1)
		go func(i interface{}) {
			defer wgOut.Done()
			mutex := &sync.Mutex{}
			wgIn := &sync.WaitGroup{}
			hashs := make([]string, thCount)
			data := i.(string)
			for th := 0; th < thCount; th++ {
				wgIn.Add(1)
				go func(th int) {
					defer wgIn.Done()
					hashs[th] = DataSignerCrc32(strconv.Itoa(th) + data)
				}(th)
			}
			wgIn.Wait()
			mutex.Lock()
			hash := strings.Join(hashs, "")
			out <- hash
			mutex.Unlock()
		}(i)
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