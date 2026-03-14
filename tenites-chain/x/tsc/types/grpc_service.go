package types

import (
	"context"

	"google.golang.org/grpc"
)

type QueryServer interface {
	Params(context.Context, *QueryParamsRequest) (*QueryParamsResponse, error)
	TotalSupply(context.Context, *QueryTotalSupplyRequest) (*QueryTotalSupplyResponse, error)
	DailyStats(context.Context, *QueryDailyStatsRequest) (*QueryDailyStatsResponse, error)
	MintAuthority(context.Context, *QueryMintAuthorityRequest) (*QueryMintAuthorityResponse, error)
	MintAuthorities(context.Context, *QueryMintAuthoritiesRequest) (*QueryMintAuthoritiesResponse, error)
	BurnAuthority(context.Context, *QueryBurnAuthorityRequest) (*QueryBurnAuthorityResponse, error)
	BurnAuthorities(context.Context, *QueryBurnAuthoritiesRequest) (*QueryBurnAuthoritiesResponse, error)
	Operation(context.Context, *QueryOperationRequest) (*QueryOperationResponse, error)
}

func RegisterMsgServer(s grpc.ServiceRegistrar, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func RegisterQueryServer(s grpc.ServiceRegistrar, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Msg_Mint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgMint)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Mint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Msg/Mint"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Mint(ctx, req.(*MsgMint))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Burn_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgBurn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Burn(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Msg/Burn"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Burn(ctx, req.(*MsgBurn))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_Transfer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgTransfer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).Transfer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Msg/Transfer"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).Transfer(ctx, req.(*MsgTransfer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddMintAuthority_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddMintAuthority)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddMintAuthority(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Msg/AddMintAuthority"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddMintAuthority(ctx, req.(*MsgAddMintAuthority))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveMintAuthority_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveMintAuthority)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveMintAuthority(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Msg/RemoveMintAuthority"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveMintAuthority(ctx, req.(*MsgRemoveMintAuthority))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_AddBurnAuthority_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgAddBurnAuthority)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).AddBurnAuthority(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Msg/AddBurnAuthority"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).AddBurnAuthority(ctx, req.(*MsgAddBurnAuthority))
	}
	return interceptor(ctx, in, info, handler)
}

func _Msg_RemoveBurnAuthority_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgRemoveBurnAuthority)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).RemoveBurnAuthority(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Msg/RemoveBurnAuthority"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).RemoveBurnAuthority(ctx, req.(*MsgRemoveBurnAuthority))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tenites.tsc.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Mint", Handler: _Msg_Mint_Handler},
		{MethodName: "Burn", Handler: _Msg_Burn_Handler},
		{MethodName: "Transfer", Handler: _Msg_Transfer_Handler},
		{MethodName: "AddMintAuthority", Handler: _Msg_AddMintAuthority_Handler},
		{MethodName: "RemoveMintAuthority", Handler: _Msg_RemoveMintAuthority_Handler},
		{MethodName: "AddBurnAuthority", Handler: _Msg_AddBurnAuthority_Handler},
		{MethodName: "RemoveBurnAuthority", Handler: _Msg_RemoveBurnAuthority_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tenites/tsc/v1/tx.proto",
}

func _Query_Params_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Params(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Query/Params"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Params(ctx, req.(*QueryParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_TotalSupply_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryTotalSupplyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).TotalSupply(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Query/TotalSupply"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).TotalSupply(ctx, req.(*QueryTotalSupplyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_DailyStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryDailyStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).DailyStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Query/DailyStats"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).DailyStats(ctx, req.(*QueryDailyStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_MintAuthority_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMintAuthorityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).MintAuthority(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Query/MintAuthority"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).MintAuthority(ctx, req.(*QueryMintAuthorityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_MintAuthorities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryMintAuthoritiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).MintAuthorities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Query/MintAuthorities"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).MintAuthorities(ctx, req.(*QueryMintAuthoritiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_BurnAuthority_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBurnAuthorityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BurnAuthority(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Query/BurnAuthority"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BurnAuthority(ctx, req.(*QueryBurnAuthorityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_BurnAuthorities_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryBurnAuthoritiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).BurnAuthorities(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Query/BurnAuthorities"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).BurnAuthorities(ctx, req.(*QueryBurnAuthoritiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Query_Operation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryOperationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Operation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{Server: srv, FullMethod: "/tenites.tsc.v1.Query/Operation"}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Operation(ctx, req.(*QueryOperationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "tenites.tsc.v1.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{MethodName: "Params", Handler: _Query_Params_Handler},
		{MethodName: "TotalSupply", Handler: _Query_TotalSupply_Handler},
		{MethodName: "DailyStats", Handler: _Query_DailyStats_Handler},
		{MethodName: "MintAuthority", Handler: _Query_MintAuthority_Handler},
		{MethodName: "MintAuthorities", Handler: _Query_MintAuthorities_Handler},
		{MethodName: "BurnAuthority", Handler: _Query_BurnAuthority_Handler},
		{MethodName: "BurnAuthorities", Handler: _Query_BurnAuthorities_Handler},
		{MethodName: "Operation", Handler: _Query_Operation_Handler},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tenites/tsc/v1/query.proto",
}

type QueryClient interface {
	Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error)
	TotalSupply(ctx context.Context, in *QueryTotalSupplyRequest, opts ...grpc.CallOption) (*QueryTotalSupplyResponse, error)
	DailyStats(ctx context.Context, in *QueryDailyStatsRequest, opts ...grpc.CallOption) (*QueryDailyStatsResponse, error)
	MintAuthority(ctx context.Context, in *QueryMintAuthorityRequest, opts ...grpc.CallOption) (*QueryMintAuthorityResponse, error)
	MintAuthorities(ctx context.Context, in *QueryMintAuthoritiesRequest, opts ...grpc.CallOption) (*QueryMintAuthoritiesResponse, error)
	BurnAuthority(ctx context.Context, in *QueryBurnAuthorityRequest, opts ...grpc.CallOption) (*QueryBurnAuthorityResponse, error)
	BurnAuthorities(ctx context.Context, in *QueryBurnAuthoritiesRequest, opts ...grpc.CallOption) (*QueryBurnAuthoritiesResponse, error)
	Operation(ctx context.Context, in *QueryOperationRequest, opts ...grpc.CallOption) (*QueryOperationResponse, error)
}

type queryClient struct {
	cc grpc.ClientConnInterface
}

func NewQueryClient(cc grpc.ClientConnInterface) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Params(ctx context.Context, in *QueryParamsRequest, opts ...grpc.CallOption) (*QueryParamsResponse, error) {
	out := new(QueryParamsResponse)
	err := c.cc.Invoke(ctx, "/tenites.tsc.v1.Query/Params", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) TotalSupply(ctx context.Context, in *QueryTotalSupplyRequest, opts ...grpc.CallOption) (*QueryTotalSupplyResponse, error) {
	out := new(QueryTotalSupplyResponse)
	err := c.cc.Invoke(ctx, "/tenites.tsc.v1.Query/TotalSupply", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) DailyStats(ctx context.Context, in *QueryDailyStatsRequest, opts ...grpc.CallOption) (*QueryDailyStatsResponse, error) {
	out := new(QueryDailyStatsResponse)
	err := c.cc.Invoke(ctx, "/tenites.tsc.v1.Query/DailyStats", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) MintAuthority(ctx context.Context, in *QueryMintAuthorityRequest, opts ...grpc.CallOption) (*QueryMintAuthorityResponse, error) {
	out := new(QueryMintAuthorityResponse)
	err := c.cc.Invoke(ctx, "/tenites.tsc.v1.Query/MintAuthority", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) MintAuthorities(ctx context.Context, in *QueryMintAuthoritiesRequest, opts ...grpc.CallOption) (*QueryMintAuthoritiesResponse, error) {
	out := new(QueryMintAuthoritiesResponse)
	err := c.cc.Invoke(ctx, "/tenites.tsc.v1.Query/MintAuthorities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) BurnAuthority(ctx context.Context, in *QueryBurnAuthorityRequest, opts ...grpc.CallOption) (*QueryBurnAuthorityResponse, error) {
	out := new(QueryBurnAuthorityResponse)
	err := c.cc.Invoke(ctx, "/tenites.tsc.v1.Query/BurnAuthority", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) BurnAuthorities(ctx context.Context, in *QueryBurnAuthoritiesRequest, opts ...grpc.CallOption) (*QueryBurnAuthoritiesResponse, error) {
	out := new(QueryBurnAuthoritiesResponse)
	err := c.cc.Invoke(ctx, "/tenites.tsc.v1.Query/BurnAuthorities", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *queryClient) Operation(ctx context.Context, in *QueryOperationRequest, opts ...grpc.CallOption) (*QueryOperationResponse, error) {
	out := new(QueryOperationResponse)
	err := c.cc.Invoke(ctx, "/tenites.tsc.v1.Query/Operation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
