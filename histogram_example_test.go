package metrics_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/VictoriaMetrics/metrics"
)

func TestExampleHistogram(t *testing.T) {
	// Define a histogram in global scope.
	var h = &metrics.Histogram{}
	var m = &metrics.Histogram{}

	h.Update(10)
	m.Update(1000)
	m.Update(99)

	buf := new(bytes.Buffer)

	h.MarshalTo("h", buf)
	t.Log(buf.String())

	buf.Reset()
	m.MarshalTo("m", buf)
	t.Log(buf.String())

	h.Merge(m)
	buf.Reset()
	h.MarshalTo("hm{path=\"test\"}", buf)
	t.Log(buf.String())
	t.FailNow()
}

func TestSerDe(t *testing.T) {
	var h = &metrics.Histogram{}
	h.Update(10)

	buf := new(bytes.Buffer)
	err := h.Serialize(buf)
	if err != nil {
		t.Logf("Can't serialize due: %s", err)
		t.FailNow()
	}

	var d = &metrics.Histogram{}
	err = d.Deserialize(buf)
	if err != nil {
		t.Logf("Can't deserialize due: %s", err)
		t.FailNow()
	}

	out := new(bytes.Buffer)
	d.MarshalToGraphite("d;bver=1.2.3", out)
	t.Log(out.String())
	t.FailNow()
}

func ExampleHistogram_vec() {
	for i := 0; i < 3; i++ {
		// Dynamically construct metric name and pass it to GetOrCreateHistogram.
		name := fmt.Sprintf(`response_size_bytes{path=%q}`, "/foo/bar")
		response := processRequest()
		metrics.GetOrCreateHistogram(name).Update(float64(len(response)))
	}
}
