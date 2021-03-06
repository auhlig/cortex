package util

import (
	"fmt"
	"reflect"
	"sort"
	"testing"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/storage/metric"
)

func TestWriteRequest(t *testing.T) {
	want := []*model.Sample{}
	for i := 0; i < 10; i++ {
		want = append(want, &model.Sample{
			Metric: model.Metric{
				model.MetricNameLabel: model.LabelValue(fmt.Sprintf("testmetric_%d", i)),
				model.JobLabel:        "testjob",
			},
			Timestamp: model.Time(i),
			Value:     model.SampleValue(float64(i)),
		})
	}

	have := FromWriteRequest(ToWriteRequest(want))

	if !reflect.DeepEqual(have, want) {
		t.Fatalf("Bad FromWriteRequest(ToWriteRequest) round trip")
	}
}

func TestQueryRequest(t *testing.T) {
	from, to := model.Time(int64(0)), model.Time(int64(10))
	matchers := []*metric.LabelMatcher{}
	matcher1, err := metric.NewLabelMatcher(metric.Equal, model.LabelName("foo"), model.LabelValue("1"))
	if err != nil {
		t.Fatal(err)
	}
	matchers = append(matchers, matcher1)

	matcher2, err := metric.NewLabelMatcher(metric.NotEqual, model.LabelName("bar"), model.LabelValue("2"))
	if err != nil {
		t.Fatal(err)
	}
	matchers = append(matchers, matcher2)

	matcher3, err := metric.NewLabelMatcher(metric.RegexMatch, model.LabelName("baz"), model.LabelValue("3"))
	if err != nil {
		t.Fatal(err)
	}
	matchers = append(matchers, matcher3)

	matcher4, err := metric.NewLabelMatcher(metric.RegexNoMatch, model.LabelName("bop"), model.LabelValue("4"))
	if err != nil {
		t.Fatal(err)
	}
	matchers = append(matchers, matcher4)

	req, err := ToQueryRequest(from, to, matchers...)
	if err != nil {
		t.Fatal(err)
	}

	haveFrom, haveTo, haveMatchers, err := FromQueryRequest(req)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(haveFrom, from) {
		t.Fatalf("Bad from FromQueryRequest(ToQueryRequest) round trip")
	}
	if !reflect.DeepEqual(haveTo, to) {
		t.Fatalf("Bad to FromQueryRequest(ToQueryRequest) round trip")
	}
	if !reflect.DeepEqual(haveMatchers, matchers) {
		t.Fatalf("Bad have FromQueryRequest(ToQueryRequest) round trip - %v != %v", haveMatchers, matchers)
	}
}

func buildTestMatrix(numSeries int, samplesPerSeries int, offset int) model.Matrix {
	m := make(model.Matrix, 0, numSeries)
	for i := 0; i < numSeries; i++ {
		ss := model.SampleStream{
			Metric: model.Metric{
				model.MetricNameLabel: model.LabelValue(fmt.Sprintf("testmetric_%d", i)),
				model.JobLabel:        "testjob",
			},
			Values: make([]model.SamplePair, 0, samplesPerSeries),
		}
		for j := 0; j < samplesPerSeries; j++ {
			ss.Values = append(ss.Values, model.SamplePair{
				Timestamp: model.Time(i + j + offset),
				Value:     model.SampleValue(i + j + offset),
			})
		}
		m = append(m, &ss)
	}
	sort.Sort(m)
	return m
}

func TestQueryResponse(t *testing.T) {
	want := buildTestMatrix(10, 10, 10)
	have := FromQueryResponse(ToQueryResponse(want))
	if !reflect.DeepEqual(have, want) {
		t.Fatalf("Bad FromQueryResponse(ToQueryResponse) round trip")
	}

}
