package ingest

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	"github.com/softrenzu/Prometheus/internal/engine"
)
func ParsePrometheusText(r io.Reader,tenant string)([]engine.Sample,error){scanner:=bufio.NewScanner(r);scanner.Buffer(make([]byte,64*1024),8*1024*1024);out:=make([]engine.Sample,0,128);now:=time.Now().UnixMilli();lineNo:=0;for scanner.Scan(){lineNo++;line:=strings.TrimSpace(scanner.Text());if line==""||strings.HasPrefix(line,"#"){continue};nameLabels,rest,ok:=strings.Cut(line," ");if !ok{return nil,fmt.Errorf("line %d: missing value",lineNo)};fields:=strings.Fields(rest);if len(fields)==0{return nil,fmt.Errorf("line %d: missing value",lineNo)};value,err:=strconv.ParseFloat(fields[0],64);if err!=nil{return nil,fmt.Errorf("line %d: invalid value: %w",lineNo,err)};ts:=now;if len(fields)>1{if parsed,err:=strconv.ParseInt(fields[1],10,64);err==nil{ts=parsed}};metric,labels,err:=parseMetricLabels(nameLabels);if err!=nil{return nil,fmt.Errorf("line %d: %w",lineNo,err)};out=append(out,engine.Sample{Tenant:tenant,Metric:metric,Labels:labels,Timestamp:ts,Value:value})};return out,scanner.Err()}
func parseMetricLabels(s string)(string,map[string]string,error){labels:=map[string]string{};i:=strings.IndexByte(s,'{');if i<0{return s,labels,nil};if !strings.HasSuffix(s,"}"){return "",nil,fmt.Errorf("invalid label set")};metric:=s[:i];body:=s[i+1:len(s)-1];var parts []string;start:=0;quoted:=false;escaped:=false;for j:=0;j<len(body);j++{c:=body[j];if escaped{escaped=false;continue};if c=='\\'{escaped=true;continue};if c=='"'{quoted=!quoted;continue};if c==','&&!quoted{parts=append(parts,body[start:j]);start=j+1}};if body!=""{parts=append(parts,body[start:])};for _,p:=range parts{kv:=strings.SplitN(strings.TrimSpace(p),"=",2);if len(kv)!=2{return "",nil,fmt.Errorf("invalid label")};labels[strings.TrimSpace(kv[0])]=strings.Trim(strings.TrimSpace(kv[1]),"\"")};return metric,labels,nil}
