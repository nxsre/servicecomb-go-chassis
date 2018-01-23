package main

// 官方的 go-sc-client 使用示例
import (
	"log"
	"time"

	"github.com/ServiceComb/go-chassis/core/lager"
	client "github.com/ServiceComb/go-sc-client"
	"github.com/ServiceComb/go-sc-client/model"
)

func main() {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	r := &client.RegistryClient{}
	err := r.Initialize(client.Options{
		Addrs: []string{"10.10.99.177:30100"},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(r.GetAllApplications())
	r.WatchMicroService("ins-ansible-agent", func(mev *model.MicroServiceInstanceChangedEvent) {
		log.Println(mev)
	})

	time.Sleep(100 * time.Second)
}
