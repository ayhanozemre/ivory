package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const ScanWaitSecond = 2

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func scanPort(ip string, port int) bool {
	result := true
	timeout := 500 * time.Millisecond
	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)

	if err != nil {
		if strings.Contains(err.Error(), "too many open files") {
			time.Sleep(timeout)
			result = scanPort(ip, port)
		} else {
			result = false
		}
	} else {
		conn.Close()
	}
	return result
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func writeCsv(columns []string, row []string) {
	filename := "result.csv"
	writeColumns := fileExists(filename)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("File creation error")
	}
	defer file.Close()
	w := csv.NewWriter(file)
	if !writeColumns {
		w.Write(columns)
	}
	w.Write(row)
	w.Flush()
}

func sendNotification(ip string, port int, callbackUrl string, label string) int {
	base, err := url.Parse(callbackUrl)
	if err != nil {
		return 0
	}
	params := url.Values{}
	params.Add("label", label)
	params.Add("ip", ip)
	params.Add("port", strconv.Itoa(port))
	base.RawQuery = params.Encode()
	rsp, err := http.Get(base.String())
	if err != nil {
		log.Println("Notification Service Error")
	}
	return rsp.StatusCode
}

func startScan(ip string, port int, storage string, label string) {
	result := scanPort(ip, port)
	if result {
		if storage == "csv" {
			columns := []string{"label", "ip", "port"}
			data := []string{label, ip, strconv.Itoa(port)}
			writeCsv(columns, data)
		} else {
			sendNotification(ip, port, storage, label)
		}
	}
}

func main() {
	/*
			go build -o scanner .
			./scanner  -port=9200    // required
		               -label=elasticSearch    // required
		               -storage=http://0.0.0.0:8080/portscanner/create    // optional
		               -first-block=35    			 // optional
		               -second-block=193 			 // optional
		               -third-block=104  			 // optional
		               -scan-size=50     			 // optional
	*/
	blocks := []int{0, 0, 0, 0}

	port := flag.Int("port", 8080, "scan port")
	label := flag.String("label", "UnknownService", "ElasticSearch/Redis/Mongodb etc...")
	storage := flag.String("storage", "csv", "notification url or csv")
	scanSize := flag.Int("scan-size", 50, "concurrent count")

	firstBlock := flag.Int("first-block", 0, "ip address first block")
	secondBlock := flag.Int("second-block", 0, "ip address second block")
	thirdBlock := flag.Int("third-block", 0, "ip address third block")
	flag.Parse()

	blocks[0] = *firstBlock
	blocks[1] = *secondBlock
	blocks[2] = *thirdBlock
	ranges := makeRange(0, 255)
	log.Println("Search starting for", *label, "...")
	targetBlock := fmt.Sprintf("%d.%d.%d.[0-255]:%d", blocks[0], blocks[1], blocks[2], *port)
	log.Println("Payload: target: ", targetBlock)
	log.Println("Storage:", *storage)
	log.Println("Search started for", *label, "!")

	for _, first := range ranges {
		for _, second := range ranges {
			for _, third := range ranges {
				for fIndex, fourth := range ranges {
					if blocks[0] != 0 {
						first = blocks[0]
					}
					if blocks[1] != 0 {
						second = blocks[1]
					}
					if blocks[2] != 0 {
						third = blocks[2]
					}
					if fIndex%*scanSize == 0 {
						time.Sleep(ScanWaitSecond * time.Second)
					}

					target := fmt.Sprintf("%d.%d.%d.%d", first, second, third, fourth)
					go startScan(target, *port, *storage, *label)
				}
				if blocks[2] != 0 {
					break
				}
			}
			if blocks[1] != 0 {
				break
			}
		}
		if blocks[0] != 0 {
			break
		}
	}
	log.Println("Search end...")
}
