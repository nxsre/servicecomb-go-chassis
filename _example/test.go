package main

import (
	"fmt"
	"log"
	"time"

	// "regexp"
	"os"

	"github.com/ServiceComb/service-center/server/core/proto"
	sc "github.com/soopsio/servicecomb-go-chassis"
)

func main() {

	sc_config := sc.ServiceComb_Config{
		SCAddr: "http://10.10.99.177:30100",
	}

	scClient := sc.NewServiceCom(&sc_config)

	// 获取服务列表
	services, err := scClient.GetAllService()
	fmt.Println(services, err)

	// 注册服务
	serviceId := "ins-ansible-agent"
	serviceName := "ins-ansible-agent_2222" // 服务提供商名称
	appId := "ins-ansible-agent_111"
	serviceVersion := "0.0.3"
	service := sc.Service{
		Service: &proto.MicroService{
			ServiceId:   serviceId,
			AppId:       appId,
			ServiceName: serviceName,
			Version:     serviceVersion,
			Description: "ansible agent",
		},
	}
	fmt.Println(scClient.RegisterService(service))

	/* 	// 服务实例监听
	   	ch := make(chan []byte)
	   	scClient.SetWatcher(ch)
	   	go scClient.InstanceListWatcher(serviceId)
	   	go func() {
	   		for {
	   			select {
	   			case msg := <-ch:
	   				log.Println("ssss:", string(msg))
	   			}
	   		}

	   	}() */

	time.Sleep(2 * time.Second)
	// 注册服务实例
	hostname, _ := os.Hostname()
	instance := sc.Instance{
		Instance: &proto.MicroServiceInstance{
			Endpoints: []string{"ssh://10.10.99.88:2022"},
			HostName:  hostname,
			Status:    "UP", // ["STRATING", "UP", "DOWN", "OUT OF SERVICE"]
			Properties: map[string]string{
				"aaa": "bbb",
			},
			HealthCheck: &proto.HealthCheck{
				Mode:     proto.CHECK_BY_HEARTBEAT,
				Interval: 5,
				Times:    2,
			},
		},
	}
	ins, err := scClient.RegisterInstance(instance, serviceId)
	if err != nil {
		log.Fatalln("RegisterInstance Failed:", err)
	}
	log.Println(ins.InstanceId)

	scClient.SetConsumerId(serviceId)

	// 查找实例，即注册消费者
	scClient.FindInstance(appId, serviceName, "0.0.1+")

	// 服务实例监听
	ch := make(chan []byte)
	scClient.SetWatcher(ch)
	go scClient.InstanceWatcher(serviceId)
	go func() {
		for {
			select {
			case msg := <-ch:
				log.Println("ssss:", string(msg))
			}
		}

	}()
	go func() {
		// 服务心跳
		for i := 0; i <= 10; i++ {
			if err := scClient.InstanceHeartBeat(serviceId, ins.InstanceId); err != nil {
				fmt.Println("心跳失败:", err)
			}
			time.Sleep(3 * time.Second)
		}
	}()

	time.Sleep(600 * time.Second)

}
