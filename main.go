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

func generateOctets(ip string) []int {
	octets := make([]int, 4)
	rawIp := strings.Split(ip, ".")
	for index, oct := range rawIp {
		newOct, err := strconv.Atoi(oct)
		if err == nil {
			// if an octet is 0, should the next ones be 0?
			// perhaps only one octet will be searched
			// if newOct == 0 {
			// 	 break
			// }
			octets[index] = newOct
		}
	}
	return octets
}

func main() {

	port := flag.Int("port", 8080, "scan port")
	label := flag.String("label", "UnknownService", "Redis/Mongodb etc...")
	storage := flag.String("storage", "csv", "notification url or csv")
	concurrentCount := flag.Int("concurrent-count", 50, "concurrent count")
	ip := flag.String("ip", "192.168.*.*", "ip address")
	flag.Parse()
	octets := generateOctets(*ip)

	log.Println("Search starting for", *label)
	targetBlock := fmt.Sprintf("%d.%d.%d.[0-255]:%d", octets[0], octets[1], octets[2], *port)
	log.Println("Payload: target: ", targetBlock)
	log.Println("Storage: ", *storage)
	log.Println("Search started for", *label)

	ranges := makeRange(0, 255)
	for _, first := range ranges {
		for _, second := range ranges {
			for _, third := range ranges {
				for fIndex, fourth := range ranges {
					if octets[0] != 0 {
						first = octets[0]
					}
					if octets[1] != 0 {
						second = octets[1]
					}
					if octets[2] != 0 {
						third = octets[2]
					}
					if fIndex%*concurrentCount == 0 {
						time.Sleep(ScanWaitSecond * time.Second)
					}

					target := fmt.Sprintf("%d.%d.%d.%d", first, second, third, fourth)
					go startScan(target, *port, "csv", "test")
				}
				if octets[2] != 0 {
					break
				}
			}
			if octets[1] != 0 {
				break
			}
		}
		if octets[0] != 0 {
			break
		}
	}
	log.Println("Search end...")
}
