package api

import("github.com/softrenzu/Prometheus/internal/engine";"strconv")
func instantData(series []engine.Series)[]map[string]any{out:=[]map[string]any{};for _,ss:=range series{if len(ss.Points)==0{continue};p:=ss.Points[len(ss.Points)-1];m:=engine.CloneLabels(ss.Labels);m["__name__"]=ss.Metric;out=append(out,map[string]any{"metric":m,"value":[]any{float64(p.Timestamp)/1000,strconv.FormatFloat(p.Value,'g',-1,64)}})};return out}
func matrixData(series []engine.Series)[]map[string]any{out:=[]map[string]any{};for _,ss:=range series{m:=engine.CloneLabels(ss.Labels);m["__name__"]=ss.Metric;vals:=[][]any{};for _,p:=range ss.Points{vals=append(vals,[]any{float64(p.Timestamp)/1000,strconv.FormatFloat(p.Value,'g',-1,64)})};out=append(out,map[string]any{"metric":m,"values":vals})};return out}
