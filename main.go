package main

import (
	"github.com/synapse-rpc/nymph"
	"fmt"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/streadway/amqp"
	"flag"
	"os"
)

func main() {
	var host = flag.String("host", "", "MqHost")
	var user = flag.String("user", "", "MqUser")
	var pass = flag.String("pass", "", "MqPass")
	var sysName = flag.String("sys_name", "", "System Name")
	var debug = flag.Bool("debug", false, "Debug Mode")
	flag.Parse();
	if *host == "" || *user == "" || *pass == "" || *sysName == "" {
		fmt.Println("Usage: go run main.go --host MQ_HOST --user MQ_USER --pass MQ_PASS --sys_name SYSTEM_NAME [--debug]")
		os.Exit(255)
	}
	app := synapse.New()
	app.MqHost = *host
	app.MqUser = *user
	app.MqPass = *pass
	app.SysName = *sysName
	app.Debug = *debug
	app.AppName = "golang"
	app.RpcCallback = map[string]func(*simplejson.Json, amqp.Delivery) map[string]interface{}{
		"test": RpcTest,
	}

	//设置事件回调方法(不设置系统将不会启动事件监听器)
	app.EventCallback = map[string]func(*simplejson.Json, amqp.Delivery) bool{
		"dotnet.test": EventTest,
		"java.test":   EventTest,
		"golang.test": EventTest,
		"python.test": EventTest,
		"ruby.test":   EventTest,
		"php.test":    EventTest,
	}
	go app.Serve()
	showHelp()
	var i1, i2, i3, i4 string
	for {
		i1 = ""
		i2 = ""
		i3 = ""
		i4 = ""
		fmt.Print("input >> ")
		fmt.Scanln(&i1, &i2, &i3, &i4)
		if i1 == "event" {
			if i3 == "" {
				showHelp()
				continue
			}
			query := map[string]interface{}{"msg": i3}
			app.SendEvent(i2, query)
		} else if i1 == "rpc" {
			if i4 == "" {
				showHelp()
				continue
			}
			query := map[string]interface{}{"msg": i4}
			ret := app.SendRpc(i2, i3, query)
			retJson, _ := json.Marshal(ret)
			fmt.Printf("\n %s \n", retJson)
		} else {
			showHelp()
		}
	}
}

func showHelp() {
	fmt.Println("----------------------------------------------")
	fmt.Println("|   event usage:                             |")
	fmt.Println("|     > event [event] [msg]                  |")
	fmt.Println("|   rpc usage:                               |")
	fmt.Println("|     > rpc [app] [method] [msg]             |")
	fmt.Println("----------------------------------------------")
}

func EventTest(r *simplejson.Json, d amqp.Delivery) bool {
	jsonStr, _ := json.Marshal(r)
	fmt.Printf("**收到EVENT: %s@%s %s", d.Type, d.ReplyTo, jsonStr)
	return true
}

func RpcTest(r *simplejson.Json, d amqp.Delivery) map[string]interface{} {
	jsonStr, _ := json.Marshal(r)
	fmt.Printf("RPC有请求: %s\n", jsonStr)
	return map[string]interface{}{
		"from":   "GoLang",
		"m":      r.MustString("msg"),
		"number": 5233,
	}
}
