package types

import (
        "bytes"
        "compress/gzip"
        "sync"

        gogoproto "github.com/cosmos/gogoproto/proto"
        protov2 "google.golang.org/protobuf/proto"
        "google.golang.org/protobuf/types/descriptorpb"
)

var (
        txFileDescGzipOnce sync.Once
        txFileDescGzip     []byte
)

func getTxFileDescGzip() []byte {
        txFileDescGzipOnce.Do(func() {
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
                                msgDescWithRepeated("MsgAddMintAuthority",
                                        []*descriptorpb.FieldDescriptorProto{
                                                fieldDesc("signer", 1, ""),
                                                fieldDesc("tenites_id", 2, ""),
                                                fieldDesc("name", 3, ""),
                                                fieldDesc("daily_limit", 4, ""),
                                                fieldDesc("single_tx_limit", 5, ""),
                                        },
                                        []*descriptorpb.FieldDescriptorProto{
                                                repeatedFieldDesc("allowed_purposes", 6),
                                                repeatedFieldDesc("jurisdictions", 7),
                                        },
                                        fieldDescVarint("expires_at", 8),
                                ),
                                msgDesc("MsgAddMintAuthorityResponse",
                                        fieldDesc("authority_id", 1, ""),
                                ),
                                msgDesc("MsgRemoveMintAuthority",
                                        fieldDesc("signer", 1, ""),
                                        fieldDesc("authority_id", 2, ""),
                                ),
                                msgDesc("MsgRemoveMintAuthorityResponse"),
                                msgDescWithRepeated("MsgAddBurnAuthority",
                                        []*descriptorpb.FieldDescriptorProto{
                                                fieldDesc("signer", 1, ""),
                                                fieldDesc("tenites_id", 2, ""),
                                                fieldDesc("name", 3, ""),
                                                fieldDesc("daily_limit", 4, ""),
                                                fieldDesc("single_tx_limit", 5, ""),
                                        },
                                        []*descriptorpb.FieldDescriptorProto{
                                                repeatedFieldDesc("allowed_reasons", 6),
                                                repeatedFieldDesc("jurisdictions", 7),
                                        },
                                        fieldDescVarint("expires_at", 8),
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
                }
                raw, err := protov2.Marshal(fd)
                if err != nil {
                        panic("marshal tx file descriptor: " + err.Error())
                }
                var buf bytes.Buffer
                w := gzip.NewWriter(&buf)
                if _, err := w.Write(raw); err != nil {
                        panic("gzip tx file descriptor: " + err.Error())
                }
                w.Close()
                txFileDescGzip = buf.Bytes()
        })
        return txFileDescGzip
}

func repeatedFieldDesc(name string, number int32) *descriptorpb.FieldDescriptorProto {
        return &descriptorpb.FieldDescriptorProto{
                Name:     strPtr(name),
                Number:   int32Ptr(number),
                Type:     descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
                Label:    descriptorpb.FieldDescriptorProto_LABEL_REPEATED.Enum(),
                JsonName: strPtr(name),
        }
}

func fieldDescVarint(name string, number int32) *descriptorpb.FieldDescriptorProto {
        return &descriptorpb.FieldDescriptorProto{
                Name:     strPtr(name),
                Number:   int32Ptr(number),
                Type:     descriptorpb.FieldDescriptorProto_TYPE_INT64.Enum(),
                Label:    descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
                JsonName: strPtr(name),
        }
}

func msgDescWithRepeated(name string, singularFields []*descriptorpb.FieldDescriptorProto, repeatedFields []*descriptorpb.FieldDescriptorProto, extraFields ...*descriptorpb.FieldDescriptorProto) *descriptorpb.DescriptorProto {
        var allFields []*descriptorpb.FieldDescriptorProto
        allFields = append(allFields, singularFields...)
        allFields = append(allFields, repeatedFields...)
        allFields = append(allFields, extraFields...)
        return &descriptorpb.DescriptorProto{
                Name:  strPtr(name),
                Field: allFields,
        }
}

func init() {
        gogoproto.RegisterFile("tenites/tsc/v1/tx.proto", getTxFileDescGzip())
}

func (m *MsgMint) Descriptor() ([]byte, []int)                          { return getTxFileDescGzip(), []int{0} }
func (m *MsgMintResponse) Descriptor() ([]byte, []int)                  { return getTxFileDescGzip(), []int{1} }
func (m *MsgBurn) Descriptor() ([]byte, []int)                          { return getTxFileDescGzip(), []int{2} }
func (m *MsgBurnResponse) Descriptor() ([]byte, []int)                  { return getTxFileDescGzip(), []int{3} }
func (m *MsgTransfer) Descriptor() ([]byte, []int)                      { return getTxFileDescGzip(), []int{4} }
func (m *MsgTransferResponse) Descriptor() ([]byte, []int)              { return getTxFileDescGzip(), []int{5} }
func (m *MsgAddMintAuthority) Descriptor() ([]byte, []int)              { return getTxFileDescGzip(), []int{6} }
func (m *MsgAddMintAuthorityResponse) Descriptor() ([]byte, []int)      { return getTxFileDescGzip(), []int{7} }
func (m *MsgRemoveMintAuthority) Descriptor() ([]byte, []int)           { return getTxFileDescGzip(), []int{8} }
func (m *MsgRemoveMintAuthorityResponse) Descriptor() ([]byte, []int)   { return getTxFileDescGzip(), []int{9} }
func (m *MsgAddBurnAuthority) Descriptor() ([]byte, []int)              { return getTxFileDescGzip(), []int{10} }
func (m *MsgAddBurnAuthorityResponse) Descriptor() ([]byte, []int)      { return getTxFileDescGzip(), []int{11} }
func (m *MsgRemoveBurnAuthority) Descriptor() ([]byte, []int)           { return getTxFileDescGzip(), []int{12} }
func (m *MsgRemoveBurnAuthorityResponse) Descriptor() ([]byte, []int)   { return getTxFileDescGzip(), []int{13} }

