package main

import (
        "io/ioutil"
        "sync"
        "fmt"
        "strings"
        "net/http"
        // "encoding/json"
        "bytes"
        "time"
)

const maxConcurrency = 30// for example

var throttle = make(chan int, maxConcurrency)

// type DiagnosisCount struct {
//     diagnosis []string
// }

// type ImageResult struct {
//     malignant string
//     id int
//     diagnosis string
//     thumbnail string
// }

// type ImageResults struct {
//     image_results []ImageResult
// }

// type CbirData struct {
//     diagnosis_count []DiagnosisCount `json:"diagnosis_count"`
//     image_results []ImageResults `json:"image_results"`
//     image_hash_value string `json:"image_hash_value"`
// }


// type CbirResponse struct {
//     result string `json:"result"`
// }

func main() {
        content , err := ioutil.ReadFile("/Users/randy/Desktop/Dermoscopic Images.txt")
        if err != nil {
            panic(err)
        }
        lines := strings.Split(strings.Replace(strings.Replace(strings.Replace(strings.Replace(strings.Replace(string(content),"[","",-1),"]","",-1)," ","",-1),"'","",-1),"\n",",",-1), ",")
        fmt.Println(len(lines))
        const N = 100000
        var wg sync.WaitGroup
        t := time.Now()
        fmt.Println("current time",t.Format("15-04-05"))
        for i := 0; i < N; i++ {
                throttle <- 100
                wg.Add(1)
                go get_url(i,t , lines ,&wg, throttle)
        }
        wg.Wait()
}

func get_url(i int,start time.Time ,lines []string , wg *sync.WaitGroup, throttle chan int) {
        defer wg.Done() 
        var jsonStr = []byte(`{"data":{"image_url":"`+ lines[i] + `", "secret_key":"6PT8e9MPwSPHf6uS4khGe5Ts7jprTC7A4Cz8BWzHEYgKFmdAw7Sm3CCKxQ2g4XtWrCGRqRbEFhZbuQnAM6LEjY2FrSy5sRSu5aVkRWAn5s997nzv4VLcwCHL67GnzasuadqCLBAGSzsAxRUrpgdKJMmEvDDEDFJqJYCmr3DTdF5m8ApnHsH34Ej99yDC4dDqN6mPKAjN6dK5qX3HfABU5KU4T9hMXx2BtX23paskY3ZZQSyGhKtckXSPHLtx85Jb"}}`)
        url := "https://www.google.ca"
        resp, err := http.Post(url,"application/json", bytes.NewBuffer(jsonStr))
        if err != nil {
            panic(err)
        }
        if resp.Status != "200 OK"{
            fmt.Println(resp.Status)
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err.Error())
        }
        fmt.Println(string(body))
        // var cbirresponse = new(CbirResponse)
        // json.Unmarshal(body, &cbirresponse)

        t1 := time.Since(start)
        fmt.Println("thread",i,"elapsed time",t1)
        <-throttle
}
