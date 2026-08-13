package ingest

import("encoding/binary";"errors";"fmt")
func consumeBytes(b []byte)([]byte,[]byte,error){l,n:=binary.Uvarint(b);if n<=0{return nil,nil,errors.New("invalid length")};b=b[n:];if l>uint64(len(b)){return nil,nil,errors.New("truncated length-delimited field")};return b[:int(l)],b[int(l):],nil}
func skipProto(b []byte,wire int)([]byte,error){switch wire{case 0:_,n:=binary.Uvarint(b);if n<=0{return nil,errors.New("invalid varint")};return b[n:],nil;case 1:if len(b)<8{return nil,errors.New("short fixed64")};return b[8:],nil;case 2:_,rest,err:=consumeBytes(b);return rest,err;case 5:if len(b)<4{return nil,errors.New("short fixed32")};return b[4:],nil;default:return nil,fmt.Errorf("unsupported protobuf wire type %d",wire)}}
