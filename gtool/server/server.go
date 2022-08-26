package server

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/protobuf/types/descriptorpb"
)

type RefServer struct {
	RefClient *grpcreflect.Client
	Service   ListInfo
	url       string
	channel   *grpc.ClientConn
}

type ListInfo struct {
	Service Service `json:"services"`
}

type Service struct {
	Name    string   `json:"name"`
	Methods []Method `json:"methods"`
	Key     string   `json:"key"`
}

type Method struct {
	Name       string `json:"name"`
	InputType  string `json:"inputType"`
	OutputType string `json:"outputType"`
	Father     string `json:"father"`
	Url        string `json:"url"`
}

func NewClient(ctx context.Context, url string) (*RefServer, error) {
	cc, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, errors.New("连接grpc失败")
	}

	refCli := &RefServer{
		RefClient: grpcreflect.NewClient(ctx, reflectpb.NewServerReflectionClient(cc)),
		url:       url,
		channel:   cc,
	}

	return refCli, nil
}

func (r *RefServer) ListService() error {
	strs, err := r.RefClient.ListServices()
	if err != nil {
		log.Println("ListError", err)
		return err
	}

	ret := ListInfo{}
	//获取service信息s
	for i := 0; i < len(strs); i++ {
		if strs[i] == "grpc.health.v1.Health" || strs[i] == "grpc.reflection.v1alpha.ServerReflection" {
			strs = append(strs[:i], strs[i+1:]...)
			i--
		}
	}
	s, err := r.RefClient.ResolveService(strs[0])
	if err != nil {
		log.Println("ResolveError", err)
		return err
	}
	ret.Service = Service{}
	ret.Service.Name = s.GetFullyQualifiedName()
	ret.Service.Key = r.url + "::" + ret.Service.Name
	methods := s.GetMethods()
	ret.Service.Methods = make([]Method, len(methods))
	for k, v := range methods {
		ret.Service.Methods[k].Name = v.GetName()
		ret.Service.Methods[k].InputType = v.GetInputType().GetName()
		ret.Service.Methods[k].OutputType = v.GetOutputType().GetName()
		ret.Service.Methods[k].Father = ret.Service.Name
		ret.Service.Methods[k].Url = r.url
	}
	r.Service = ret
	return nil

}

func (r *RefServer) Call(serviceName, methodName, jsonString string) (map[string]interface{}, error) {
	methodDesc, err := getMethodDescriptor(r.RefClient, serviceName, methodName)
	if err != nil {
		return nil, err
	}

	req := dynamic.NewMessage(methodDesc.GetInputType())
	req.UnmarshalJSON([]byte(jsonString))
	stub := grpcdynamic.NewStub(r.channel)

	resp, err := stub.InvokeRpc(context.Background(), methodDesc, req)
	if err != nil {
		return nil, err
	}
	ret := make(map[string]interface{})
	js := resp.(*dynamic.Message)
	ty, _ := js.MarshalJSON()

	json.Unmarshal(ty, &ret)
	return ret, nil
}

func getMethodDescriptor(refClient *grpcreflect.Client, serviceName, methodName string) (*desc.MethodDescriptor, error) {
	var st *desc.MethodDescriptor
	s, e := refClient.ResolveService(serviceName)
	if e != nil {
		return nil, e
	}
	for _, v := range s.GetMethods() {
		if v.GetName() == methodName {
			st = v
			break
		}
	}
	return st, nil
}

func (r *RefServer) Close() {
	r.channel.Close()
}

func (r *RefServer) GetParams(serviceName, methodName string) (map[string]interface{}, error) {
	methodDesc, err := getMethodDescriptor(r.RefClient, serviceName, methodName)
	if err != nil {
		return nil, err
	}

	up := make(map[string]interface{})
	down := make(map[string]interface{})
	for _, fieldDescriptor := range methodDesc.GetInputType().GetFields() {
		fieldName := fieldDescriptor.GetName()
		if fieldDescriptor.IsRepeated() {
			// 如果是一个数组的话，就返回 nil 吧
			down[fieldName] = nil
			continue
		}
		switch fieldDescriptor.GetType() {
		case descriptorpb.FieldDescriptorProto_Type(descriptor.FieldDescriptorProto_TYPE_MESSAGE):
			down[fieldName] = convertMessageToMap(fieldDescriptor.GetMessageType())
			continue
		}
		down[fieldName] = fieldDescriptor.GetDefaultValue()
	}

	up["methods"] = down

	/*
		m := make(map[string]interface{})


			for _, msgDescriptor := range methodDesc.GetInputType() {
				fmt.Println(msgDescriptor)
				if methodName == msgDescriptor.GetName() {
					m[msgDescriptor.GetName()] = convertMessageToMap(msgDescriptor)
				}
			}
			fmt.Println(m)
	*/

	return up, nil
}

func convertMessageToMap(message *desc.MessageDescriptor) map[string]interface{} {
	m := make(map[string]interface{})
	for _, fieldDescriptor := range message.GetFields() {
		fieldName := fieldDescriptor.GetName()
		if fieldDescriptor.IsRepeated() {
			// 如果是一个数组的话，就返回 nil 吧
			m[fieldName] = nil
			continue
		}
		switch fieldDescriptor.GetType() {
		case descriptorpb.FieldDescriptorProto_Type(descriptor.FieldDescriptorProto_TYPE_MESSAGE):
			m[fieldName] = convertMessageToMap(fieldDescriptor.GetMessageType())
			continue
		}
		m[fieldName] = fieldDescriptor.GetDefaultValue()
	}
	return m
}
