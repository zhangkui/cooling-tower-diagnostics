package ingest

import "testing"

func TestDecode(t *testing.T) {
	r, e := Decoder{}.Decode([]byte(`{"tower_id":"t","sensor":" Conductivity ","value":2}`))
	if e != nil || r.Sensor != "conductivity" {
		t.Fatalf("%#v %v", r, e)
	}
}
