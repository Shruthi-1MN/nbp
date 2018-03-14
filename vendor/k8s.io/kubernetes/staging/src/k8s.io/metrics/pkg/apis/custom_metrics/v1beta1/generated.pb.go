/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by protoc-gen-gogo.
// source: k8s.io/kubernetes/vendor/k8s.io/metrics/pkg/apis/custom_metrics/v1beta1/generated.proto
// DO NOT EDIT!

/*
	Package v1beta1 is a generated protocol buffer package.

	It is generated from these files:
		k8s.io/kubernetes/vendor/k8s.io/metrics/pkg/apis/custom_metrics/v1beta1/generated.proto

	It has these top-level messages:
		MetricValue
		MetricValueList
*/
package v1beta1

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"

import strings "strings"
import reflect "reflect"

import io "io"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

func (m *MetricValue) Reset()                    { *m = MetricValue{} }
func (*MetricValue) ProtoMessage()               {}
func (*MetricValue) Descriptor() ([]byte, []int) { return fileDescriptorGenerated, []int{0} }

func (m *MetricValueList) Reset()                    { *m = MetricValueList{} }
func (*MetricValueList) ProtoMessage()               {}
func (*MetricValueList) Descriptor() ([]byte, []int) { return fileDescriptorGenerated, []int{1} }

func init() {
	proto.RegisterType((*MetricValue)(nil), "k8s.io.metrics.pkg.apis.custom_metrics.v1beta1.MetricValue")
	proto.RegisterType((*MetricValueList)(nil), "k8s.io.metrics.pkg.apis.custom_metrics.v1beta1.MetricValueList")
}
func (m *MetricValue) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetricValue) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintGenerated(dAtA, i, uint64(m.DescribedObject.Size()))
	n1, err := m.DescribedObject.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n1
	dAtA[i] = 0x12
	i++
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.MetricName)))
	i += copy(dAtA[i:], m.MetricName)
	dAtA[i] = 0x1a
	i++
	i = encodeVarintGenerated(dAtA, i, uint64(m.Timestamp.Size()))
	n2, err := m.Timestamp.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n2
	if m.WindowSeconds != nil {
		dAtA[i] = 0x20
		i++
		i = encodeVarintGenerated(dAtA, i, uint64(*m.WindowSeconds))
	}
	dAtA[i] = 0x2a
	i++
	i = encodeVarintGenerated(dAtA, i, uint64(m.Value.Size()))
	n3, err := m.Value.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n3
	return i, nil
}

func (m *MetricValueList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MetricValueList) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	dAtA[i] = 0xa
	i++
	i = encodeVarintGenerated(dAtA, i, uint64(m.ListMeta.Size()))
	n4, err := m.ListMeta.MarshalTo(dAtA[i:])
	if err != nil {
		return 0, err
	}
	i += n4
	if len(m.Items) > 0 {
		for _, msg := range m.Items {
			dAtA[i] = 0x12
			i++
			i = encodeVarintGenerated(dAtA, i, uint64(msg.Size()))
			n, err := msg.MarshalTo(dAtA[i:])
			if err != nil {
				return 0, err
			}
			i += n
		}
	}
	return i, nil
}

func encodeFixed64Generated(dAtA []byte, offset int, v uint64) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	dAtA[offset+4] = uint8(v >> 32)
	dAtA[offset+5] = uint8(v >> 40)
	dAtA[offset+6] = uint8(v >> 48)
	dAtA[offset+7] = uint8(v >> 56)
	return offset + 8
}
func encodeFixed32Generated(dAtA []byte, offset int, v uint32) int {
	dAtA[offset] = uint8(v)
	dAtA[offset+1] = uint8(v >> 8)
	dAtA[offset+2] = uint8(v >> 16)
	dAtA[offset+3] = uint8(v >> 24)
	return offset + 4
}
func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *MetricValue) Size() (n int) {
	var l int
	_ = l
	l = m.DescribedObject.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.MetricName)
	n += 1 + l + sovGenerated(uint64(l))
	l = m.Timestamp.Size()
	n += 1 + l + sovGenerated(uint64(l))
	if m.WindowSeconds != nil {
		n += 1 + sovGenerated(uint64(*m.WindowSeconds))
	}
	l = m.Value.Size()
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *MetricValueList) Size() (n int) {
	var l int
	_ = l
	l = m.ListMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovGenerated(uint64(l))
		}
	}
	return n
}

func sovGenerated(x uint64) (n int) {
	for {
		n++
		x >>= 7
		if x == 0 {
			break
		}
	}
	return n
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *MetricValue) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MetricValue{`,
		`DescribedObject:` + strings.Replace(strings.Replace(this.DescribedObject.String(), "ObjectReference", "k8s_io_api_core_v1.ObjectReference", 1), `&`, ``, 1) + `,`,
		`MetricName:` + fmt.Sprintf("%v", this.MetricName) + `,`,
		`Timestamp:` + strings.Replace(strings.Replace(this.Timestamp.String(), "Time", "k8s_io_apimachinery_pkg_apis_meta_v1.Time", 1), `&`, ``, 1) + `,`,
		`WindowSeconds:` + valueToStringGenerated(this.WindowSeconds) + `,`,
		`Value:` + strings.Replace(strings.Replace(this.Value.String(), "Quantity", "k8s_io_apimachinery_pkg_api_resource.Quantity", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *MetricValueList) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&MetricValueList{`,
		`ListMeta:` + strings.Replace(strings.Replace(this.ListMeta.String(), "ListMeta", "k8s_io_apimachinery_pkg_apis_meta_v1.ListMeta", 1), `&`, ``, 1) + `,`,
		`Items:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Items), "MetricValue", "MetricValue", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *MetricValue) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetricValue: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetricValue: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DescribedObject", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.DescribedObject.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MetricName", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= (uint64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MetricName = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Timestamp", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Timestamp.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field WindowSeconds", wireType)
			}
			var v int64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= (int64(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.WindowSeconds = &v
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Value", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Value.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGenerated
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
func (m *MetricValueList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: MetricValueList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MetricValueList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ListMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ListMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, MetricValue{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthGenerated
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
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
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
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
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
			iNdEx += length
			if length < 0 {
				return 0, ErrInvalidLengthGenerated
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowGenerated
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipGenerated(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthGenerated = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated   = fmt.Errorf("proto: integer overflow")
)

func init() {
	proto.RegisterFile("k8s.io/kubernetes/vendor/k8s.io/metrics/pkg/apis/custom_metrics/v1beta1/generated.proto", fileDescriptorGenerated)
}

var fileDescriptorGenerated = []byte{
	// 546 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x93, 0x41, 0x6f, 0xd3, 0x3e,
	0x18, 0xc6, 0x9b, 0xf5, 0xdf, 0x3f, 0x9b, 0xab, 0xa9, 0x2c, 0x17, 0xa2, 0x1e, 0xd2, 0xaa, 0x5c,
	0x0a, 0xd2, 0x6c, 0xb5, 0x20, 0x84, 0xc4, 0x2d, 0xe2, 0x82, 0x44, 0x41, 0x64, 0x13, 0x93, 0x00,
	0x09, 0x9c, 0xe4, 0x6d, 0x6a, 0xba, 0xc4, 0x91, 0xed, 0x74, 0xda, 0x8d, 0x13, 0x67, 0x3e, 0x56,
	0x8f, 0xe3, 0xb6, 0x53, 0x45, 0xc3, 0x17, 0x41, 0x49, 0x9c, 0xb6, 0x5b, 0x19, 0xb0, 0x5b, 0x6c,
	0xbf, 0xcf, 0xcf, 0xcf, 0xf3, 0xbe, 0x0e, 0x3a, 0x99, 0x3e, 0x95, 0x98, 0x71, 0x32, 0x4d, 0x3d,
	0x10, 0x31, 0x28, 0x90, 0x64, 0x06, 0x71, 0xc0, 0x05, 0xd1, 0x07, 0x11, 0x28, 0xc1, 0x7c, 0x49,
	0x92, 0x69, 0x48, 0x68, 0xc2, 0x24, 0xf1, 0x53, 0xa9, 0x78, 0xf4, 0xb1, 0xda, 0x9f, 0x0d, 0x3c,
	0x50, 0x74, 0x40, 0x42, 0x88, 0x41, 0x50, 0x05, 0x01, 0x4e, 0x04, 0x57, 0xdc, 0xc4, 0xa5, 0x1e,
	0xeb, 0x3a, 0x9c, 0x4c, 0x43, 0x9c, 0xeb, 0xf1, 0x55, 0x3d, 0xd6, 0xfa, 0xf6, 0x61, 0xc8, 0xd4,
	0x24, 0xf5, 0xb0, 0xcf, 0x23, 0x12, 0xf2, 0x90, 0x93, 0x02, 0xe3, 0xa5, 0xe3, 0x62, 0x55, 0x2c,
	0x8a, 0xaf, 0x12, 0xdf, 0xee, 0x69, 0x7b, 0x34, 0x61, 0xc4, 0xe7, 0x02, 0xc8, 0x6c, 0xcb, 0x42,
	0xfb, 0xf1, 0xba, 0x26, 0xa2, 0xfe, 0x84, 0xc5, 0x20, 0xce, 0xab, 0x1c, 0x44, 0x80, 0xe4, 0xa9,
	0xf0, 0xe1, 0x56, 0x2a, 0x99, 0xb7, 0x83, 0xfe, 0xee, 0x2e, 0x72, 0x93, 0x4a, 0xa4, 0xb1, 0x62,
	0xd1, 0xf6, 0x35, 0x4f, 0xfe, 0x26, 0x90, 0xfe, 0x04, 0x22, 0xba, 0xa5, 0x7b, 0x74, 0x93, 0x2e,
	0x55, 0xec, 0x94, 0xb0, 0x58, 0x49, 0x25, 0xae, 0x8b, 0x7a, 0x5f, 0xeb, 0xa8, 0x39, 0x2a, 0x1a,
	0xfe, 0x96, 0x9e, 0xa6, 0x60, 0x8e, 0x51, 0x2b, 0x00, 0xe9, 0x0b, 0xe6, 0x41, 0xf0, 0xda, 0xfb,
	0x0c, 0xbe, 0xb2, 0x8c, 0xae, 0xd1, 0x6f, 0x0e, 0xef, 0x57, 0x63, 0xa3, 0x09, 0xc3, 0x79, 0x5f,
	0xf1, 0x6c, 0x80, 0xcb, 0x0a, 0x17, 0xc6, 0x20, 0x20, 0xf6, 0xc1, 0xb9, 0x37, 0x5f, 0x74, 0x6a,
	0xd9, 0xa2, 0xd3, 0x7a, 0x7e, 0x95, 0xe1, 0x5e, 0x87, 0x9a, 0x43, 0x84, 0xca, 0x39, 0xbf, 0xa2,
	0x11, 0x58, 0x3b, 0x5d, 0xa3, 0xbf, 0xe7, 0x98, 0x5a, 0x8d, 0x46, 0xab, 0x13, 0x77, 0xa3, 0xca,
	0x7c, 0x8f, 0xf6, 0xf2, 0xfc, 0x52, 0xd1, 0x28, 0xb1, 0xea, 0x85, 0xab, 0x87, 0x1b, 0xae, 0x56,
	0xa1, 0xd7, 0x2f, 0x2a, 0x9f, 0x49, 0xee, 0xf3, 0x98, 0x45, 0xe0, 0x1c, 0x68, 0xfc, 0xde, 0x71,
	0x05, 0x71, 0xd7, 0x3c, 0xf3, 0x01, 0xfa, 0xff, 0x8c, 0xc5, 0x01, 0x3f, 0xb3, 0xfe, 0xeb, 0x1a,
	0xfd, 0xba, 0x73, 0x90, 0x2d, 0x3a, 0xfb, 0x27, 0xc5, 0xce, 0x11, 0xf8, 0x3c, 0x0e, 0xa4, 0xab,
	0x0b, 0xcc, 0x23, 0xd4, 0x98, 0xe5, 0xcd, 0xb2, 0x1a, 0x85, 0x07, 0xfc, 0x27, 0x0f, 0xb8, 0x7a,
	0x4d, 0xf8, 0x4d, 0x4a, 0x63, 0xc5, 0xd4, 0xb9, 0xb3, 0xaf, 0x7d, 0x34, 0x8a, 0x8e, 0xbb, 0x25,
	0xab, 0xf7, 0xdd, 0x40, 0xad, 0x8d, 0x41, 0xbc, 0x64, 0x52, 0x99, 0x1f, 0xd0, 0x6e, 0x9e, 0x20,
	0xa0, 0x8a, 0xea, 0x29, 0xe0, 0x7f, 0xcb, 0x9b, 0xab, 0x47, 0xa0, 0xa8, 0x73, 0x57, 0xdf, 0xb5,
	0x5b, 0xed, 0xb8, 0x2b, 0xa2, 0xf9, 0x09, 0x35, 0x98, 0x82, 0x48, 0x5a, 0x3b, 0xdd, 0x7a, 0xbf,
	0x39, 0x7c, 0x76, 0xcb, 0xff, 0x12, 0x6f, 0xb8, 0x5d, 0x67, 0x7a, 0x91, 0x13, 0xdd, 0x12, 0xec,
	0x1c, 0xce, 0x97, 0x76, 0xed, 0x62, 0x69, 0xd7, 0x2e, 0x97, 0x76, 0xed, 0x4b, 0x66, 0x1b, 0xf3,
	0xcc, 0x36, 0x2e, 0x32, 0xdb, 0xb8, 0xcc, 0x6c, 0xe3, 0x47, 0x66, 0x1b, 0xdf, 0x7e, 0xda, 0xb5,
	0x77, 0x77, 0x34, 0xf0, 0x57, 0x00, 0x00, 0x00, 0xff, 0xff, 0x64, 0x6b, 0xb8, 0x76, 0x72, 0x04,
	0x00, 0x00,
}
