package ingest

import("encoding/binary";"math";"strings";"testing")
func TestPrometheusText(t *testing.T){s,err:=ParsePrometheusText(strings.NewReader("cpu_usage{host=\"a\"} 1.5 1234\n"),"t");if err!=nil{t.Fatal(err)};if len(s)!=1||s[0].Metric!="cpu_usage"||s[0].Labels["host"]!="a"||s[0].Value!=1.5{t.Fatalf("bad sample %+v",s)}}
func field(n,wire uint64)[]byte{return appendVarint(nil,(n<<3)|wire)}
func appendVarint(dst []byte,v uint64)[]byte{var b [10]byte;n:=binary.PutUvarint(b[:],v);return append(dst,b[:n]...)}
func bytesField(n uint64,b []byte)[]byte{out:=field(n,2);out=appendVarint(out,uint64(len(b)));return append(out,b...)}
func label(k,v string)[]byte{b:=bytesField(1,[]byte(k));return append(b,bytesField(2,[]byte(v))...)}
func sample(v float64,ts int64)[]byte{b:=field(1,1);var x [8]byte;binary.LittleEndian.PutUint64(x[:],math.Float64bits(v));b=append(b,x[:]...);b=append(b,field(2,0)...);return appendVarint(b,uint64(ts))}
func literalSnappy(raw []byte)[]byte{out:=appendVarint(nil,uint64(len(raw)));n:=len(raw);if n<=60{out=append(out,byte((n-1)<<2))}else{t:=n-1;out=append(out,byte(60<<2),byte(t))};return append(out,raw...)}
func TestRemoteWriteV1(t *testing.T){ts:=bytesField(1,label("__name__","cpu"));ts=append(ts,bytesField(1,label("host","a"))...);ts=append(ts,bytesField(2,sample(3.5,1234))...);wr:=bytesField(1,ts);got,err:=DecodeRemoteWriteV1(literalSnappy(wr),"t");if err!=nil{t.Fatal(err)};if len(got)!=1||got[0].Metric!="cpu"||got[0].Labels["host"]!="a"||got[0].Value!=3.5||got[0].Timestamp!=1234{t.Fatalf("bad %+v",got)}}
