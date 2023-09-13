// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kava/community/v1beta1/params.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/types"
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
	UpgradeTimeDisableInflation time.Time  `protobuf:"bytes,1,opt,name=upgrade_time_disable_inflation,json=upgradeTimeDisableInflation,proto3,stdtime" json:"upgrade_time_disable_inflation"`
	RewardsPerSecond            types.Coin `protobuf:"bytes,2,opt,name=rewards_per_second,json=rewardsPerSecond,proto3" json:"rewards_per_second"`
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

func (m *Params) GetRewardsPerSecond() types.Coin {
	if m != nil {
		return m.RewardsPerSecond
	}
	return types.Coin{}
}

func init() {
	proto.RegisterType((*Params)(nil), "kava.community.v1beta1.Params")
}

func init() {
	proto.RegisterFile("kava/community/v1beta1/params.proto", fileDescriptor_0a48475520900507)
}

var fileDescriptor_0a48475520900507 = []byte{
	// 314 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x91, 0xc1, 0x4a, 0xc3, 0x30,
	0x1c, 0xc6, 0x1b, 0x91, 0x21, 0xf5, 0x22, 0x45, 0x64, 0x4e, 0xc8, 0x44, 0x2f, 0x82, 0x98, 0x30,
	0x7d, 0x83, 0xa9, 0x07, 0x0f, 0xc2, 0x98, 0x9e, 0xbc, 0x94, 0xa4, 0xcd, 0x6a, 0xb0, 0xed, 0x3f,
	0x24, 0xe9, 0x74, 0x6f, 0xb1, 0x67, 0xf2, 0xb4, 0xe3, 0x8e, 0x9e, 0x54, 0xd6, 0x17, 0x91, 0xa4,
	0xed, 0xf0, 0xf6, 0x25, 0xf9, 0xbe, 0x5f, 0xf8, 0xfe, 0xff, 0xf0, 0xfc, 0x8d, 0xcd, 0x19, 0x4d,
	0xa0, 0x28, 0xaa, 0x52, 0xda, 0x05, 0x9d, 0x8f, 0xb8, 0xb0, 0x6c, 0x44, 0x15, 0xd3, 0xac, 0x30,
	0x44, 0x69, 0xb0, 0x10, 0x1d, 0x39, 0x13, 0xd9, 0x9a, 0x48, 0x6b, 0x1a, 0xe0, 0x04, 0x4c, 0x01,
	0x86, 0x72, 0x66, 0xc4, 0x36, 0x99, 0x80, 0x2c, 0x9b, 0xdc, 0xe0, 0x30, 0x83, 0x0c, 0xbc, 0xa4,
	0x4e, 0xb5, 0xb7, 0xc3, 0x0c, 0x20, 0xcb, 0x05, 0xf5, 0x27, 0x5e, 0xcd, 0xa8, 0x95, 0x85, 0x30,
	0x96, 0x15, 0xaa, 0x31, 0x9c, 0x7d, 0xa2, 0xb0, 0x37, 0xf1, 0xff, 0x47, 0x32, 0xc4, 0x95, 0xca,
	0x34, 0x4b, 0x45, 0xec, 0x5c, 0x71, 0x2a, 0x0d, 0xe3, 0xb9, 0x88, 0x65, 0x39, 0xcb, 0x99, 0x95,
	0x50, 0xf6, 0xd1, 0x29, 0xba, 0xd8, 0xbf, 0x1e, 0x90, 0x06, 0x4a, 0x3a, 0x28, 0x79, 0xee, 0xa0,
	0xe3, 0xbd, 0xd5, 0xf7, 0x30, 0x58, 0xfe, 0x0c, 0xd1, 0xf4, 0xa4, 0x65, 0xb9, 0xb7, 0xbb, 0x86,
	0xf4, 0xd0, 0x81, 0xa2, 0xc7, 0x30, 0xd2, 0xe2, 0x9d, 0xe9, 0xd4, 0xc4, 0x4a, 0xe8, 0xd8, 0x88,
	0x04, 0xca, 0xb4, 0xbf, 0xe3, 0xf1, 0xc7, 0xa4, 0x69, 0x4a, 0x5c, 0xd3, 0xae, 0x3e, 0xb9, 0x05,
	0x59, 0x8e, 0x77, 0x1d, 0x7d, 0x7a, 0xd0, 0x46, 0x27, 0x42, 0x3f, 0xf9, 0xe0, 0xf8, 0x7e, 0xb5,
	0xc1, 0x68, 0xbd, 0xc1, 0xe8, 0x77, 0x83, 0xd1, 0xb2, 0xc6, 0xc1, 0xba, 0xc6, 0xc1, 0x57, 0x8d,
	0x83, 0x97, 0xcb, 0x4c, 0xda, 0xd7, 0x8a, 0xbb, 0x79, 0x52, 0x37, 0xd8, 0xab, 0x9c, 0x71, 0xe3,
	0x15, 0xfd, 0xf8, 0xb7, 0x09, 0xbb, 0x50, 0xc2, 0xf0, 0x9e, 0x2f, 0x74, 0xf3, 0x17, 0x00, 0x00,
	0xff, 0xff, 0x11, 0x8a, 0x61, 0x75, 0xa8, 0x01, 0x00, 0x00,
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
		size, err := m.RewardsPerSecond.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintParams(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	n2, err2 := github_com_gogo_protobuf_types.StdTimeMarshalTo(m.UpgradeTimeDisableInflation, dAtA[i-github_com_gogo_protobuf_types.SizeOfStdTime(m.UpgradeTimeDisableInflation):])
	if err2 != nil {
		return 0, err2
	}
	i -= n2
	i = encodeVarintParams(dAtA, i, uint64(n2))
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
	l = m.RewardsPerSecond.Size()
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
				return fmt.Errorf("proto: wrong wireType = %d for field RewardsPerSecond", wireType)
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
			if err := m.RewardsPerSecond.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
