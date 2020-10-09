package main

import (
	"fmt"
	"time"
	"net"
	"flag"
	"strconv"
	"math/rand"
	"github.com/AlekSi/zabbix-sender"
)

func nextIP(ip net.IP, inc uint) net.IP {
        i := ip.To4()
        v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
        v += inc
        v3 := byte(v & 0xFF)
        v2 := byte((v >> 8) & 0xFF)
        v1 := byte((v >> 16) & 0xFF)
        v0 := byte((v >> 24) & 0xFF)
        return net.IPv4(v0, v1, v2, v3)
}

func makeHosts(qty int) ([]string) {
	hosts := make([]string, qty)
	ip := net.ParseIP("203.0.113.1")
	for i := 0; i < qty; i++ {
		hosts[i] = ip.String()
		ip = nextIP(ip, 1)
	}
	return hosts
}

func random(min, max float64) float64 {
    rand.Seed(time.Now().UnixNano())
    return rand.Float64()*(max-min) + min
}

func sendValue(hosts []string, itemPerHost int, zabbixHost string) {
	data := map[string]interface{}{}
	di := zabbix_sender.MakeDataItems(data, hosts[0])
	for _, host := range hosts {
		for i := 0; i < itemPerHost; i++ {
			key := fmt.Sprintf("key.%d", i)
			val := strconv.FormatFloat(random(0.0, 100.0), 'f', 2, 64)
			di = append(di, zabbix_sender.DataItem{Hostname: host, Key: key, Value: val})
		}
	}
	addr, _ := net.ResolveTCPAddr("tcp", zabbixHost + ":10051")
	res, err := zabbix_sender.Send(addr, di)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(res)
}

func main() {
	var (
                hostNum int
                itemPerHost int
		zabbixHost string
        )
	flag.IntVar(&hostNum, "host", 1, "host num")
        flag.IntVar(&itemPerHost, "item", 1, "item num per host")
	flag.StringVar(&zabbixHost, "zabbix", "", "zabbix host addr")
	flag.Parse()

	hosts := makeHosts(hostNum)
	fmt.Println(hosts)
	t := time.NewTicker(1 * time.Second)
	for {
		start := time.Now()
		select {
		case <-t.C:
			fmt.Println("SEND")
			sendValue(hosts, itemPerHost, zabbixHost)
		}
		end := time.Now()
		fmt.Printf("%f秒\n", (end.Sub(start)).Seconds())
	}
	t.Stop()
}

