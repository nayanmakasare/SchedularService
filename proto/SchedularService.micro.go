// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: SchedularService.proto

package SchedularService

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for SchedularService service

type SchedularService interface {
	ApplySchedule(ctx context.Context, in *ApplyRequest, opts ...client.CallOption) (*ServiceResponse, error)
	MakeSchedule(ctx context.Context, in *Schedule, opts ...client.CallOption) (*ServiceResponse, error)
	GetSchedule(ctx context.Context, in *RequestSchedule, opts ...client.CallOption) (*UserScheduleResponse, error)
}

type schedularService struct {
	c    client.Client
	name string
}

func NewSchedularService(name string, c client.Client) SchedularService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "SchedularService"
	}
	return &schedularService{
		c:    c,
		name: name,
	}
}

func (c *schedularService) ApplySchedule(ctx context.Context, in *ApplyRequest, opts ...client.CallOption) (*ServiceResponse, error) {
	req := c.c.NewRequest(c.name, "SchedularService.ApplySchedule", in)
	out := new(ServiceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedularService) MakeSchedule(ctx context.Context, in *Schedule, opts ...client.CallOption) (*ServiceResponse, error) {
	req := c.c.NewRequest(c.name, "SchedularService.MakeSchedule", in)
	out := new(ServiceResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *schedularService) GetSchedule(ctx context.Context, in *RequestSchedule, opts ...client.CallOption) (*UserScheduleResponse, error) {
	req := c.c.NewRequest(c.name, "SchedularService.GetSchedule", in)
	out := new(UserScheduleResponse)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SchedularService service

type SchedularServiceHandler interface {
	ApplySchedule(context.Context, *ApplyRequest, *ServiceResponse) error
	MakeSchedule(context.Context, *Schedule, *ServiceResponse) error
	GetSchedule(context.Context, *RequestSchedule, *UserScheduleResponse) error
}

func RegisterSchedularServiceHandler(s server.Server, hdlr SchedularServiceHandler, opts ...server.HandlerOption) error {
	type schedularService interface {
		ApplySchedule(ctx context.Context, in *ApplyRequest, out *ServiceResponse) error
		MakeSchedule(ctx context.Context, in *Schedule, out *ServiceResponse) error
		GetSchedule(ctx context.Context, in *RequestSchedule, out *UserScheduleResponse) error
	}
	type SchedularService struct {
		schedularService
	}
	h := &schedularServiceHandler{hdlr}
	return s.Handle(s.NewHandler(&SchedularService{h}, opts...))
}

type schedularServiceHandler struct {
	SchedularServiceHandler
}

func (h *schedularServiceHandler) ApplySchedule(ctx context.Context, in *ApplyRequest, out *ServiceResponse) error {
	return h.SchedularServiceHandler.ApplySchedule(ctx, in, out)
}

func (h *schedularServiceHandler) MakeSchedule(ctx context.Context, in *Schedule, out *ServiceResponse) error {
	return h.SchedularServiceHandler.MakeSchedule(ctx, in, out)
}

func (h *schedularServiceHandler) GetSchedule(ctx context.Context, in *RequestSchedule, out *UserScheduleResponse) error {
	return h.SchedularServiceHandler.GetSchedule(ctx, in, out)
}
