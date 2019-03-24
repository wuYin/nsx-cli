package client

import (
	"encoding/json"
	"fmt"
	"nix/codec"
	"nix/service"
	"reflect"
	"time"
	"tron"
)

type Consumer struct {
	uri2Service   map[string]service.Service
	clientManager *NixClientManager
}

func NewConsumer(services []service.Service) *Consumer {
	c := &Consumer{uri2Service: make(map[string]service.Service)}

	var uris []string
	for _, s := range services {
		c.uri2Service[s.Uri] = s
		uris = append(uris, s.Uri)
	}
	c.clientManager = NewNixClientManager("localhost:8080", uris)

	for _, s := range services {
		c.makeRPCFunc(s)
	}

	return c
}

func (c *Consumer) GetService(uri string) interface{} {
	return c.uri2Service[uri].Instance
}

// 封装目标方法与网络调用
func (c Consumer) makeRPCFunc(s service.Service) {
	p := reflect.ValueOf(s.Instance).Elem() // Instance 其实是 Proxy
	pt := reflect.TypeOf(p.Interface())
	for i := 0; i < p.NumField(); i++ { // 全部是值为 nil 的方法字段
		m := p.Field(i)
		mName := pt.Field(i).Name
		mType := m.Type()

		f := func(uri string, method string, mType reflect.Type) func(in []reflect.Value) []reflect.Value {
			return func(in []reflect.Value) []reflect.Value {
				return c.call(uri, method, mType, in)
			}
		}(s.Uri, mName, mType)

		v := reflect.MakeFunc(mType, f)
		m.Set(v)
	}
}

// 请求网络发起调用并返回
func (c Consumer) call(uri string, method string, mType reflect.Type, in []reflect.Value) []reflect.Value {
	// 1. 选取已连接到指定服务的 client
	cli, err := c.clientManager.SelectClient(uri)
	if err != nil {
		panic("call: select client: " + err.Error())
	}

	// 2. 组装请求数据
	args := make([]interface{}, len(in))
	for i, v := range in {
		args[i] = v.Interface()
	}
	cmd := codec.CmdReq{
		ServiceUri: uri,
		Method:     method,
		Args:       args,
	}
	data, _ := json.Marshal(cmd)

	// 3. 打包请求数据
	reqPack := tron.NewReqPacket(data)

	// 4. 同步调用
	res, err := cli.SyncWrite(reqPack, 2*time.Second)
	if err != nil {
		panic("call: sync write: " + err.Error())
	}
	// fmt.Printf("SyncWriteAndRead finish: %+v\n", string(res.([]byte)))

	// 5. 处理调用结果
	var callResp codec.CallResp
	if err = json.Unmarshal(res.([]byte), &callResp); err != nil {
		panic("call: unmarshal failed:" + err.Error())
	}

	respV := reflect.ValueOf(callResp.Res)
	if !respV.Type().ConvertibleTo(mType.Out(0)) {
		panic(fmt.Sprintf("can't convert resp %+v to %+v", respV.Type(), mType.Out(0)))
	}

	return []reflect.Value{respV.Convert(mType.Out(0))}
}
