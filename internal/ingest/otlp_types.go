package ingest

type otlpRequest struct{ResourceMetrics []otlpResourceMetrics `json:"resourceMetrics"`}
type otlpResourceMetrics struct{Resource otlpResource `json:"resource"`;ScopeMetrics []otlpScopeMetrics `json:"scopeMetrics"`}
type otlpResource struct{Attributes []otlpAttr `json:"attributes"`}
type otlpScopeMetrics struct{Metrics []otlpMetric `json:"metrics"`}
type otlpMetric struct{Name string `json:"name"`;Gauge *otlpPoints `json:"gauge"`;Sum *otlpPoints `json:"sum"`;Histogram *otlpHistogram `json:"histogram"`}
type otlpPoints struct{DataPoints []otlpPoint `json:"dataPoints"`}
type otlpHistogram struct{DataPoints []otlpHistPoint `json:"dataPoints"`}
type otlpPoint struct{Attributes []otlpAttr `json:"attributes"`;TimeUnixNano string `json:"timeUnixNano"`;AsDouble *float64 `json:"asDouble"`;AsInt *string `json:"asInt"`}
type otlpHistPoint struct{Attributes []otlpAttr `json:"attributes"`;TimeUnixNano string `json:"timeUnixNano"`;Count string `json:"count"`;Sum *float64 `json:"sum"`}
type otlpAttr struct{Key string `json:"key"`;Value otlpValue `json:"value"`}
type otlpValue struct{StringValue *string `json:"stringValue"`;IntValue *string `json:"intValue"`;DoubleValue *float64 `json:"doubleValue"`;BoolValue *bool `json:"boolValue"`}
