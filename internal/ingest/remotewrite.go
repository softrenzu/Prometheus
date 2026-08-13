package ingest

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/softrenzu/Prometheus/internal/engine"
)
func DecodeRemoteWriteV1(snappyPayload []byte,tenant string)([]engine.Sample,error){raw,err:=decodeSnappy(snappyPayload);if err!=nil{return nil,fmt.Errorf("snappy: %w",err)};out:=[]engine.Sample{};for len(raw)>0{key,n:=binary.Uvarint(raw);if n<=0{return nil,errors.New("invalid protobuf key")};raw=raw[n:];field,wire:=int(key>>3),int(key&7);if field==1&&wire==2{msg,rest,err:=consumeBytes(raw);if err!=nil{return nil,err};raw=rest;samples,err:=decodeTimeSeries(msg,tenant);if err!=nil{return nil,err};out=append(out,samples...)}else{raw,err=skipProto(raw,wire);if err!=nil{return nil,err}}};return out,nil}
