package servicecomb

import (
	"github.com/ServiceComb/service-center/server/core/proto"
)

type AllServices struct {
	Services []*proto.MicroService `protobuf:"bytes,1,rep,name=services" json:"services,omitempty"`
}

type Service struct {
	Service *proto.MicroService `protobuf:"bytes,1,opt,name=service" json:"service,omitempty"`
}

type Instance struct {
	Instance *proto.MicroServiceInstance `protobuf:"bytes,4,opt,name=instance" json:"instance,omitempty"`
}
