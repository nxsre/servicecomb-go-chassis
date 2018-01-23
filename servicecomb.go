package servicecomb

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/ServiceComb/service-center/server/core/proto"
	scerr "github.com/ServiceComb/service-center/server/error"
	"github.com/gorilla/websocket"
	"github.com/imroc/req"
	"github.com/tidwall/gjson"
	gws "golang.org/x/net/websocket"
)

type ServiceComb struct {
	config     *ServiceComb_Config
	domain     string
	consumerId string
	watcher    chan []byte
}

func NewServiceCom(config *ServiceComb_Config) *ServiceComb {
	return &ServiceComb{
		config:  config,
		domain:  "default",
		watcher: nil,
	}
}

func (this *ServiceComb) SetDomain(domain string) {
	this.domain = domain
}

func (this *ServiceComb) SetConsumerId(consumerId string) {
	this.consumerId = consumerId
}

func (this *ServiceComb) SetWatcher(ch chan []byte) {
	this.watcher = ch
}

// GetAllService 获取服务列表
func (this *ServiceComb) GetAllService() (AllServices, error) {
	url := this.config.SCAddr + Path(GETALLSERVICE, nil)
	header := req.Header{}
	header["x-domain-name"] = this.domain

	services := AllServices{}

	resp, err := req.Get(url, header)
	if err != nil {
		return services, err
	}
	if resp.Response().StatusCode == 200 {
		err = gjson.Unmarshal(resp.Bytes(), &services)
		return services, err
	} else {
		sce := scerr.Error{}
		if err := gjson.Unmarshal(resp.Bytes(), &sce); err != nil {
			return services, errors.New(resp.String())
		}
		return services, sce
	}
}

// RegisterService 注册服务
func (this *ServiceComb) RegisterService(service Service) (svc *proto.MicroService, err error) {
	url := this.config.SCAddr + Path(REGISTERMICROSERVICE, nil)
	header := req.Header{}
	header["x-domain-name"] = this.domain

	resp, err := req.Post(url, header, req.BodyJSON(&service))
	if err != nil {
		return nil, err
	}
	if resp.Response().StatusCode == 200 {

		if err := gjson.Unmarshal(resp.Bytes(), &svc); err != nil {
			return nil, errors.New(resp.String())
		}
		return
	} else {
		sce := scerr.Error{}
		if err := gjson.Unmarshal(resp.Bytes(), &sce); err != nil {
			return nil, errors.New(resp.String())
		}
		return nil, sce
	}
	return
}

// RegisterInstance 注册服务实例
func (this *ServiceComb) RegisterInstance(instance Instance, serviceId string) (ins *proto.MicroServiceInstance, err error) {
	url := this.config.SCAddr + Path(REGISTERINSTANCE, map[string]string{"serviceId": serviceId})
	header := req.Header{}
	header["x-domain-name"] = this.domain

	resp, err := req.Post(url, header, req.BodyJSON(&instance))
	if err != nil {
		return
	}
	if resp.Response().StatusCode == 200 {

		if err := gjson.Unmarshal(resp.Bytes(), &ins); err != nil {
			return nil, errors.New(resp.String())
		}
		return
	} else {
		sce := scerr.Error{}
		if err := gjson.Unmarshal(resp.Bytes(), &sce); err != nil {
			return nil, errors.New(resp.String())
		}
		return nil, sce
	}

	return
}

// InstanceHeartBeat 实例心跳
func (this *ServiceComb) InstanceHeartBeat(serviceId, instanceId string) error {
	url := this.config.SCAddr + Path(INSTANCEHEARTBEAT, map[string]string{"serviceId": serviceId, "instanceId": instanceId})
	header := req.Header{}
	header["x-domain-name"] = this.domain
	resp, err := req.Put(url, header)
	if err != nil {
		return err
	}
	if resp.Response().StatusCode != 200 {
		sce := scerr.Error{}
		if err := gjson.Unmarshal(resp.Bytes(), &sce); err != nil {
			return errors.New(resp.String())
		}
		return sce
	}
	return nil
}

// InstanceWatcher 监听实例
func (this *ServiceComb) InstanceWatcher(serviceId string) {
	tu, err := url.Parse(this.config.SCAddr + Path(INSTANCEWATCHER, map[string]string{"serviceId": serviceId}))
	if err != nil {
		log.Println(err)
	}
	initWebsocketClient(tu, this.domain)
	/* 	interrupt := make(chan os.Signal, 1)
	   	signal.Notify(interrupt, os.Interrupt)

	   	tu, err := url.Parse(this.config.SCAddr)
	   	if err != nil {
	   		log.Println(err)
	   	}
	   	scheme := "ws"
	   	if tu.Scheme == "https" {
	   		scheme = "wss"
	   	}
	   	u := url.URL{Scheme: scheme, Host: tu.Host, Path: Path(INSTANCEWATCHER, map[string]string{"serviceId": serviceId})}
	   	log.Printf("connecting to %s", u.String())

	   	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"x-domain-name": []string{this.domain}})
	   	if err != nil {
	   		log.Fatal("dial:", err)
	   	}
	   	defer c.Close()

	   	done := make(chan struct{})

	   	go func() {
	   		defer c.Close()
	   		defer close(done)
	   		for {
	   			_, message, err := c.ReadMessage()
	   			if err != nil {
	   				log.Println("read:", err)
	   				return
	   			}
	   			log.Printf("recv: %s", message)
	   			if this.watcher != nil {
	   				this.watcher <- message
	   			}
	   		}
	   	}()

	   	ticker := time.NewTicker(time.Second)
	   	defer ticker.Stop()

	   	for {
	   		select {
	   		case t := <-ticker.C:
	   			// err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
	   			// if err != nil {
	   			// 	log.Println("write:", err)
	   			// 	return
	   			// }
	   			log.Println("time", t.String())
	   		case <-interrupt:
	   			log.Println("interrupt")
	   			// To cleanly close a connection, a client should send a close
	   			// frame and wait for the server to close the connection.
	   			// err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	   			// if err != nil {
	   			// 	log.Println("write close:", err)
	   			// 	return
	   			// }
	   			select {
	   			case <-done:
	   			case <-time.After(time.Second):
	   			}
	   			c.Close()
	   			return
	   		}
	   	} */

}

func initWebsocketClient(u *url.URL, domain string) {
	fmt.Println("Starting Client")
	wsurl, _ := url.Parse(u.String())
	if u.Scheme == "http" {
		wsurl.Scheme = "ws"
	}
	wsc := &gws.Config{
		Location: wsurl,
		Header:   http.Header{"x-domain-name": []string{domain}},
		Origin:   u,
		Version:  13,
		// https://github.com/gorilla/websocket/blob/master/client.go
		// req.Header["Sec-WebSocket-Version"] = []string{"13"}
	}
	ws, err := gws.DialConfig(wsc)
	// ws, err := gws.Dial("ws://"+u.Host+u.Path, "", fmt.Sprintf("http://%s/", u.Host))
	if err != nil {
		fmt.Printf("Dial failed: %s\n", err.Error())
		os.Exit(1)
	}

	incomingMessages := make(chan string)
	go readClientMessages(ws, incomingMessages)
	i := 0
	for {
		select {
		case <-time.After(time.Duration(2e9)):
			i++
			fmt.Println(i)
		case message := <-incomingMessages:
			fmt.Println(`Message Received:`, message)
		}
	}
}

func readClientMessages(ws *gws.Conn, incomingMessages chan string) {
	for {
		var message string
		// err := websocket.JSON.Receive(ws, &message)
		err := gws.Message.Receive(ws, &message)
		if err != nil {
			fmt.Printf("Error::: %s\n", err.Error())
			return
		}
		incomingMessages <- message
	}
}

// InstanceListWatcher 监听实例
// service center(sc)提供的watcher api是需要使用ws方式请求订阅，
// 同时api中的serviceid指的是consumer的serviceid，而不是producer的。
// 正确的使用方式是，c端调用调用一次/v4/registry/instances?serviceName=P缓存一次实例集合，
// 然后发起watcher接口，接口的serviceid填写的是c的id，保持订阅变更消息刷新本地缓存。
func (this *ServiceComb) InstanceListWatcher(serviceId string) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	tu, err := url.Parse(this.config.SCAddr)
	if err != nil {
		log.Println(err)
	}
	scheme := "ws"
	if tu.Scheme == "https" {
		scheme = "wss"
	}
	u := url.URL{Scheme: scheme, Host: tu.Host, Path: Path(INSTANCELISTWATCHER, map[string]string{"serviceId": serviceId})}
	log.Printf("connecting to %s", u.String())

	c, res, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"x-domain-name": []string{this.domain}})
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	bs, err := ioutil.ReadAll(res.Body)
	log.Println(string(bs), err)
	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			// log.Printf("recv: %s", message)
			if this.watcher != nil {
				this.watcher <- message
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
			// log.Println(t.String())
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			c.Close()
			return
		}
	}
}

// FindInstance 查找实例，也是注册消费者
func (this *ServiceComb) FindInstance(serviceAppId, serviceName, serviceVersion string) error {
	// "?appId="+serviceAppId+"&serviceName="+serviceName+"&version="+serviceVersion
	// FINDINSTANCE
	url := this.config.SCAddr + Path(FINDINSTANCE, nil)
	header := req.Header{}
	header["x-domain-name"] = this.domain
	header["X-ConsumerId"] = this.consumerId

	qparam := req.QueryParam{}
	qparam["appId"] = serviceAppId
	qparam["serviceName"] = serviceName
	qparam["version"] = serviceVersion

	resp, err := req.Get(url, header, qparam)
	if err != nil {
		return err
	}
	log.Println(resp.Request())
	log.Println(resp.String())
	if resp.Response().StatusCode != 200 {
		sce := scerr.Error{}
		if err := gjson.Unmarshal(resp.Bytes(), &sce); err != nil {
			return errors.New(resp.String())
		}
		return sce
	}
	return nil
}
