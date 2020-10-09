package main

import (
	"github.com/AlekSi/zabbix"
	"fmt"
	"net"
	"flag"
)

type API struct {
	*zabbix.API
}

func (api *API) initHost(ip string) (nh *zabbix.Host) {
	i := zabbix.HostInterface{DNS: "",
                          IP: ip,
                          Main: 1,
                          Port: "10050",
                          Type: 1,
                          UseIP: 1}
	interfaces := zabbix.HostInterfaces{i}

	g := zabbix.HostGroupId{GroupId: "2"}
        groups := zabbix.HostGroupIds{g}

	h := zabbix.Host{Host: ip,
                 Available: 1,
                 GroupIds: groups,
                 Name: ip,
                 Status: 0,
                 Interfaces: interfaces}
	hs := zabbix.Hosts{h}
	err := api.HostsCreate(hs)
        if err != nil {
		fmt.Println(err)
        }
	nh, err = api.HostGetByHost(ip)
        if err != nil {
                fmt.Println(err)
        }
	return nh
}

func (api *API) initHostApplication(host *zabbix.Host) (a *zabbix.Application) {
	as := zabbix.Applications{{HostId: host.HostId, Name: fmt.Sprintf("App for %s", host.Host)}}
        err := api.ApplicationsCreate(as)
        if err != nil {
                fmt.Println(err)
        }

        a, err = api.ApplicationGetByHostIdAndName(host.HostId, fmt.Sprintf("App for %s", host.Host))
        if err != nil {
                fmt.Println(err)
        }
	return a
}

func (api *API) initItem(host *zabbix.Host, app *zabbix.Application, qty int) (nitms []interface{}) {
	itms := zabbix.Items{}
	for i := 0; i < qty; i++ {
		itm := zabbix.Item{HostId: app.HostId,
                                Key: fmt.Sprintf("key.%d", i),
                                Name: fmt.Sprintf("Item %d for %s", i, host.Host),
                                Type: zabbix.ZabbixTrapper,
                                ApplicationIds: []string{app.ApplicationId}}
		itms = append(itms, itm)
	}

	err := api.ItemsCreate(itms)
	if err != nil {
		fmt.Println(err)
        }

	resp, err := api.CallWithError("item.get", zabbix.Params{"output": "extend", "hostids": app.HostId})
        if err != nil {
                fmt.Println(err)
        }
	nitms = resp.Result.([]interface{})
        //fmt.Println(response.Result.([]interface{})[0].(map[string]interface{})["itemid"])
	return nitms
}

func (api *API) initTrigger(host *zabbix.Host, qty int) {
	for i := 0; i < qty; i++ {
		_, err := api.CallWithError("trigger.create", zabbix.Params{"expression": fmt.Sprintf("{%s:key.%d.last()}>90",host.Host, i),"description": fmt.Sprintf("trigger %s key.%d", host.Host, i)})
		if err != nil {
			fmt.Println(err)
		}
	}
}

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

func initZabbix(server string, user string, pass string, hostNum int, itemPerHost int) {
        a := zabbix.NewAPI(server + "/api_jsonrpc.php")
	api := &API{a}
        api.Login(user, pass)

	ip := net.ParseIP("203.0.113.1")
	for i := 0; i < hostNum; i++ {
		host := api.initHost(ip.String())
		fmt.Println("HOST:", host.Name, host.HostId, host.Interfaces)

		app := api.initHostApplication(host)
		fmt.Println("APRICATIONS:", app)

		items := api.initItem(host, app, itemPerHost)
		fmt.Println("ITEMS:" , items[0].(map[string]interface{})["itemid"])

		api.initTrigger(host, itemPerHost)
		ip = nextIP(ip, 1)
	}
}


func main() {
	var (
		hostNum int
		itemPerHost int
		zabbixHost string
		user string
		pass string
	)
	flag.IntVar(&hostNum, "host", 1, "host num")
	flag.IntVar(&itemPerHost, "item", 1, "item num per host")
	flag.StringVar(&zabbixHost, "zabbix", "", "zabbix host url")
	flag.StringVar(&user, "user", "Admin", "zabbix user")
	flag.StringVar(&pass, "pass", "password", "zabbix pass")
	flag.Parse()
	initZabbix(zabbixHost, user, pass, hostNum, itemPerHost)
}


