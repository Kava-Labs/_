// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kava/community/v1beta1/params.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_types "github.com/gogo/protobuf/types"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	io "io"
	math "math"
	math_bits "math/bits"
	time "time"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf
var _ = time.Kitchen

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Params defines the parameters of the community module.
type Params struct {
	UpgradeTimeDisableInflation time.Time             `protobuf:"bytes,1,opt,name=upgrade_time_disable_inflation,json=upgradeTimeDisableInflation,proto3,stdtime" json:"upgrade_time_disable_inflation"`
	StakingRewardsPerSecond     cosmossdk_io_math.Int `protobuf:"bytes,2,opt,name=staking_rewards_per_second,json=stakingRewardsPerSecond,proto3,customtype=cosmossdk.io/math.Int" json:"staking_rewards_per_second"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_0a48475520900507, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetUpgradeTimeDisableInflation() time.Time {
	if m != nil {
		return m.UpgradeTimeDisableInflation
	}
	return time.Time{}
}

func init() {
	proto.RegisterType((*Params)(nil), "kava.community.v1beta1.Params")
}

func init() {
	proto.RegisterFile("kava/community/v1beta1/params.proto", fileDescriptor_0a48475520900507)
}

var fileDescriptor_0a48475520900507 = []byte{
	// 336 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xcf, 0x4e, 0x2a, 0x31,
	0x14, 0x87, 0xa7, 0x77, 0x41, 0xae, 0xe3, 0x8e, 0xf8, 0x07, 0xc7, 0x64, 0x86, 0xe8, 0x86, 0x84,
	0xd0, 0x06, 0x7d, 0x03, 0xa2, 0x0b, 0x76, 0x04, 0x5d, 0xb9, 0x69, 0x3a, 0x4c, 0x29, 0x0d, 0xd3,
	0x76, 0xd2, 0x76, 0x50, 0xde, 0x82, 0x87, 0xf1, 0x21, 0x58, 0x12, 0x57, 0xc6, 0x44, 0x34, 0xf0,
	0x22, 0xa6, 0xd3, 0xc1, 0xb8, 0x3b, 0xe7, 0xf4, 0xeb, 0xd7, 0xfc, 0x4e, 0xc3, 0xeb, 0x39, 0x59,
	0x10, 0x34, 0x51, 0x42, 0x94, 0x92, 0xdb, 0x25, 0x5a, 0xf4, 0x53, 0x6a, 0x49, 0x1f, 0x15, 0x44,
	0x13, 0x61, 0x60, 0xa1, 0x95, 0x55, 0xcd, 0x33, 0x07, 0xc1, 0x5f, 0x08, 0xd6, 0x50, 0x74, 0x31,
	0x51, 0x46, 0x28, 0x83, 0x2b, 0x0a, 0xf9, 0xc6, 0x5f, 0x89, 0x4e, 0x98, 0x62, 0xca, 0xcf, 0x5d,
	0x55, 0x4f, 0x13, 0xa6, 0x14, 0xcb, 0x29, 0xaa, 0xba, 0xb4, 0x9c, 0x22, 0xcb, 0x05, 0x35, 0x96,
	0x88, 0xc2, 0x03, 0x57, 0x9f, 0x20, 0x6c, 0x8c, 0xaa, 0xa7, 0x9b, 0x3c, 0x8c, 0xcb, 0x82, 0x69,
	0x92, 0x51, 0xec, 0x28, 0x9c, 0x71, 0x43, 0xd2, 0x9c, 0x62, 0x2e, 0xa7, 0x39, 0xb1, 0x5c, 0xc9,
	0x16, 0x68, 0x83, 0xce, 0xf1, 0x4d, 0x04, 0xbd, 0x14, 0x1e, 0xa4, 0xf0, 0xf1, 0x20, 0x1d, 0xfc,
	0x5f, 0x6f, 0x93, 0x60, 0xf5, 0x95, 0x80, 0xf1, 0x65, 0xed, 0x72, 0x67, 0x77, 0xde, 0x34, 0x3c,
	0x88, 0x9a, 0xb3, 0x30, 0x32, 0x96, 0xcc, 0xb9, 0x64, 0x58, 0xd3, 0x67, 0xa2, 0x33, 0x83, 0x0b,
	0xaa, 0xb1, 0xa1, 0x13, 0x25, 0xb3, 0xd6, 0xbf, 0x36, 0xe8, 0x1c, 0x0d, 0xba, 0x4e, 0xf5, 0xb1,
	0x4d, 0x4e, 0x7d, 0x4c, 0x93, 0xcd, 0x21, 0x57, 0x48, 0x10, 0x3b, 0x83, 0x43, 0x69, 0xdf, 0x5e,
	0x7b, 0x61, 0x9d, 0x7f, 0x28, 0xed, 0xf8, 0xbc, 0xd6, 0x8d, 0xbd, 0x6d, 0x44, 0xf5, 0x43, 0xe5,
	0x1a, 0xdc, 0xaf, 0x77, 0x31, 0xd8, 0xec, 0x62, 0xf0, 0xbd, 0x8b, 0xc1, 0x6a, 0x1f, 0x07, 0x9b,
	0x7d, 0x1c, 0xbc, 0xef, 0xe3, 0xe0, 0xa9, 0xcb, 0xb8, 0x9d, 0x95, 0xa9, 0xdb, 0x32, 0x72, 0xeb,
	0xee, 0xe5, 0x24, 0x35, 0x55, 0x85, 0x5e, 0xfe, 0xfc, 0x8f, 0x5d, 0x16, 0xd4, 0xa4, 0x8d, 0x2a,
	0xeb, 0xed, 0x4f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xc5, 0xad, 0xf1, 0xc6, 0xbe, 0x01, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.StakingRewardsPerSecond.Size()
		i -= size
		if _, err := m.StakingRewardsPerSecond.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	n1, err1 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpgradeTimeDisableInflation, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.UpgradeTimeDisableInflation):])
	if err1 != nil {
		return 0, err1
	}
	i -= n1
	i = encodeVarintParams(dAtA, i, uint64(n1))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintParams(dAtA []byte, offset int, v uint64) int {
	offset -= sovParams(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = github_com_gogo_protobuf_types.SizeOfStdTime(m.UpgradeTimeDisableInflation)
	n += 1 + l + sovParams(uint64(l))
	l = m.StakingRewardsPerSecond.Size()
	n += 1 + l + sovParams(uint64(l))
	return n
}

func sovParams(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozParams(x uint64) (n int) {
	return sovParams(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowParams
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UpgradeTimeDisableInflation", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := github_com_gogo_protobuf_types.StdTimeUnmarshal(&m.UpgradeTimeDisableInflation, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field StakingRewardsPerSecond", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowParams
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthParams
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthParams
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.StakingRewardsPerSecond.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipParams(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthParams
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipParams(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowParams
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowParams
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowParams
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthParams
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupParams
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthParams
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthParams        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowParams          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupParams = fmt.Errorf("proto: unexpected end of group")
)