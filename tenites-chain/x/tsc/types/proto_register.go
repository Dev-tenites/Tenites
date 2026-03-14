package types

import (
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func init() {
	registerTxFileDescriptor()
	registerQueryFileDescriptor()
}

func strPtr(s string) *string { return &s }
func boolPtr(b bool) *bool   { return &b }
func int32Ptr(i int32) *int32 { return &i }

func fieldDesc(name string, number int32, typeName string) *descriptorpb.FieldDescriptorProto {
	return &descriptorpb.FieldDescriptorProto{
		Name:     strPtr(name),
		Number:   int32Ptr(number),
		Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
		Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
		JsonName: strPtr(name),
	}
}

func msgDesc(name string, fields ...*descriptorpb.FieldDescriptorProto) *descriptorpb.DescriptorProto {
	return &descriptorpb.DescriptorProto{
		Name:  strPtr(name),
		Field: fields,
	}
}

func methodDesc(name, input, output string) *descriptorpb.MethodDescriptorProto {
	return &descriptorpb.MethodDescriptorProto{
		Name:       strPtr(name),
		InputType:  strPtr("." + input),
		OutputType: strPtr("." + output),
	}
}

func registerTxFileDescriptor() {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    strPtr("tenites/tsc/v1/tx.proto"),
		Package: strPtr("tenites.tsc.v1"),
		Syntax:  strPtr("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: strPtr("github.com/tenites/tenites-chain/x/tsc/types"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			msgDesc("MsgMint",
				fieldDesc("signer", 1, ""),
				fieldDesc("authority_id", 2, ""),
				fieldDesc("wallet_id", 3, ""),
				fieldDesc("amount", 4, ""),
				fieldDesc("purpose", 5, ""),
				fieldDesc("jurisdiction", 6, ""),
			),
			msgDesc("MsgMintResponse",
				fieldDesc("operation_id", 1, ""),
			),
			msgDesc("MsgBurn",
				fieldDesc("signer", 1, ""),
				fieldDesc("authority_id", 2, ""),
				fieldDesc("wallet_id", 3, ""),
				fieldDesc("amount", 4, ""),
				fieldDesc("reason", 5, ""),
				fieldDesc("jurisdiction", 6, ""),
			),
			msgDesc("MsgBurnResponse"),
			msgDesc("MsgTransfer",
				fieldDesc("signer", 1, ""),
				fieldDesc("from_wallet_id", 2, ""),
				fieldDesc("to_wallet_id", 3, ""),
				fieldDesc("amount", 4, ""),
				fieldDesc("purpose", 5, ""),
			),
			msgDesc("MsgTransferResponse"),
			msgDesc("MsgAddMintAuthority",
				fieldDesc("signer", 1, ""),
				fieldDesc("tenites_id", 2, ""),
				fieldDesc("name", 3, ""),
				fieldDesc("daily_limit", 4, ""),
				fieldDesc("single_tx_limit", 5, ""),
			),
			msgDesc("MsgAddMintAuthorityResponse",
				fieldDesc("authority_id", 1, ""),
			),
			msgDesc("MsgRemoveMintAuthority",
				fieldDesc("signer", 1, ""),
				fieldDesc("authority_id", 2, ""),
			),
			msgDesc("MsgRemoveMintAuthorityResponse"),
			msgDesc("MsgAddBurnAuthority",
				fieldDesc("signer", 1, ""),
				fieldDesc("tenites_id", 2, ""),
				fieldDesc("name", 3, ""),
				fieldDesc("daily_limit", 4, ""),
				fieldDesc("single_tx_limit", 5, ""),
			),
			msgDesc("MsgAddBurnAuthorityResponse",
				fieldDesc("authority_id", 1, ""),
			),
			msgDesc("MsgRemoveBurnAuthority",
				fieldDesc("signer", 1, ""),
				fieldDesc("authority_id", 2, ""),
			),
			msgDesc("MsgRemoveBurnAuthorityResponse"),
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strPtr("Msg"),
				Method: []*descriptorpb.MethodDescriptorProto{
					methodDesc("Mint", "tenites.tsc.v1.MsgMint", "tenites.tsc.v1.MsgMintResponse"),
					methodDesc("Burn", "tenites.tsc.v1.MsgBurn", "tenites.tsc.v1.MsgBurnResponse"),
					methodDesc("Transfer", "tenites.tsc.v1.MsgTransfer", "tenites.tsc.v1.MsgTransferResponse"),
					methodDesc("AddMintAuthority", "tenites.tsc.v1.MsgAddMintAuthority", "tenites.tsc.v1.MsgAddMintAuthorityResponse"),
					methodDesc("RemoveMintAuthority", "tenites.tsc.v1.MsgRemoveMintAuthority", "tenites.tsc.v1.MsgRemoveMintAuthorityResponse"),
					methodDesc("AddBurnAuthority", "tenites.tsc.v1.MsgAddBurnAuthority", "tenites.tsc.v1.MsgAddBurnAuthorityResponse"),
					methodDesc("RemoveBurnAuthority", "tenites.tsc.v1.MsgRemoveBurnAuthority", "tenites.tsc.v1.MsgRemoveBurnAuthorityResponse"),
				},
			},
		},
	}

	fileDesc, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("failed to create tx file descriptor: " + err.Error())
	}
	if _, err := protoregistry.GlobalFiles.FindFileByPath("tenites/tsc/v1/tx.proto"); err != nil {
		protoregistry.GlobalFiles.RegisterFile(fileDesc)
	}
}

func registerQueryFileDescriptor() {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    strPtr("tenites/tsc/v1/query.proto"),
		Package: strPtr("tenites.tsc.v1"),
		Syntax:  strPtr("proto3"),
		Options: &descriptorpb.FileOptions{
			GoPackage: strPtr("github.com/tenites/tenites-chain/x/tsc/types"),
		},
		MessageType: []*descriptorpb.DescriptorProto{
			msgDesc("QueryParamsRequest"),
			msgDesc("QueryParamsResponse"),
			msgDesc("QueryTotalSupplyRequest"),
			msgDesc("QueryTotalSupplyResponse",
				fieldDesc("total_supply", 1, ""),
			),
			msgDesc("QueryDailyStatsRequest"),
			msgDesc("QueryDailyStatsResponse",
				fieldDesc("daily_minted", 1, ""),
				fieldDesc("daily_burned", 2, ""),
			),
			msgDesc("QueryMintAuthorityRequest",
				fieldDesc("authority_id", 1, ""),
			),
			msgDesc("QueryMintAuthorityResponse"),
			msgDesc("QueryMintAuthoritiesRequest"),
			msgDesc("QueryMintAuthoritiesResponse"),
			msgDesc("QueryBurnAuthorityRequest",
				fieldDesc("authority_id", 1, ""),
			),
			msgDesc("QueryBurnAuthorityResponse"),
			msgDesc("QueryBurnAuthoritiesRequest"),
			msgDesc("QueryBurnAuthoritiesResponse"),
			msgDesc("QueryOperationRequest",
				fieldDesc("operation_id", 1, ""),
			),
			msgDesc("QueryOperationResponse"),
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: strPtr("Query"),
				Method: []*descriptorpb.MethodDescriptorProto{
					methodDesc("Params", "tenites.tsc.v1.QueryParamsRequest", "tenites.tsc.v1.QueryParamsResponse"),
					methodDesc("TotalSupply", "tenites.tsc.v1.QueryTotalSupplyRequest", "tenites.tsc.v1.QueryTotalSupplyResponse"),
					methodDesc("DailyStats", "tenites.tsc.v1.QueryDailyStatsRequest", "tenites.tsc.v1.QueryDailyStatsResponse"),
					methodDesc("MintAuthority", "tenites.tsc.v1.QueryMintAuthorityRequest", "tenites.tsc.v1.QueryMintAuthorityResponse"),
					methodDesc("MintAuthorities", "tenites.tsc.v1.QueryMintAuthoritiesRequest", "tenites.tsc.v1.QueryMintAuthoritiesResponse"),
					methodDesc("BurnAuthority", "tenites.tsc.v1.QueryBurnAuthorityRequest", "tenites.tsc.v1.QueryBurnAuthorityResponse"),
					methodDesc("BurnAuthorities", "tenites.tsc.v1.QueryBurnAuthoritiesRequest", "tenites.tsc.v1.QueryBurnAuthoritiesResponse"),
					methodDesc("Operation", "tenites.tsc.v1.QueryOperationRequest", "tenites.tsc.v1.QueryOperationResponse"),
				},
			},
		},
	}

	fileDesc, err := protodesc.NewFile(fd, nil)
	if err != nil {
		panic("failed to create query file descriptor: " + err.Error())
	}
	if _, err := protoregistry.GlobalFiles.FindFileByPath("tenites/tsc/v1/query.proto"); err != nil {
		protoregistry.GlobalFiles.RegisterFile(fileDesc)
	}
}
