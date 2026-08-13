package ingest

import("fmt";"strconv";"time";"github.com/softrenzu/Prometheus/internal/engine")
func nanoToMilli(s string)int64{if s==""{return time.Now().UnixMilli()};n,err:=strconv.ParseInt(s,10,64);if err!=nil{return time.Now().UnixMilli()};return n/1_000_000}
func attrs(in []otlpAttr)map[string]string{out:=map[string]string{};for _,a:=range in{switch{case a.Value.StringValue!=nil:out[a.Key]=*a.Value.StringValue;case a.Value.IntValue!=nil:out[a.Key]=*a.Value.IntValue;case a.Value.DoubleValue!=nil:out[a.Key]=fmt.Sprintf("%g",*a.Value.DoubleValue);case a.Value.BoolValue!=nil:out[a.Key]=strconv.FormatBool(*a.Value.BoolValue)}};return out}
func merge(a,b map[string]string)map[string]string{out:=engine.CloneLabels(a);for k,v:=range b{out[k]=v};return out}
