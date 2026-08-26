package model

type SensorCatalog struct {
	Kind        string
	DisplayName string
	DefaultUnit string
	MinValue    float64
	MaxValue    float64
	Precision   int
}

var Catalog = map[string]SensorCatalog{
	"conductivity": {Kind: "conductivity", DisplayName: "电导率", DefaultUnit: "uS/cm", MinValue: 0, MaxValue: 5000, Precision: 1},
	"outlet_temp":  {Kind: "outlet_temp", DisplayName: "出水温度", DefaultUnit: "C", MinValue: -20, MaxValue: 90, Precision: 2},
	"inlet_temp":   {Kind: "inlet_temp", DisplayName: "进水温度", DefaultUnit: "C", MinValue: -20, MaxValue: 90, Precision: 2},
	"flow":         {Kind: "flow", DisplayName: "循环流量", DefaultUnit: "m3/h", MinValue: 0, MaxValue: 10000, Precision: 2},
	"vibration":    {Kind: "vibration", DisplayName: "振动幅度", DefaultUnit: "mm/s", MinValue: 0, MaxValue: 100, Precision: 3},
}

func LookupSensor(kind string) (SensorCatalog, bool) {
	v, ok := Catalog[kind]
	return v, ok
}

func SensorKinds() []string {
	out := make([]string, 0, len(Catalog))
	for k := range Catalog {
		out = append(out, k)
	}
	return out
}

func IsPlausible(kind string, value float64) bool {
	v, ok := LookupSensor(kind)
	if !ok {
		return false
	}
	return value >= v.MinValue && value <= v.MaxValue
}
