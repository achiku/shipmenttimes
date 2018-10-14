package main

import (
	"bytes"
	"os"
	"testing"
)

func TestParseBaseCSV(t *testing.T) {
	f, err := os.Open("./testdata/base-input.csv")
	if err != nil {
		t.Fatal(err)
	}
	ls, err := ParseBaseCSV(f)
	if err != nil {
		t.Fatal(err)
	}
	for _, l := range ls {
		t.Logf("%+v", l)
	}
}

func TestTransformBaseOrder(t *testing.T) {
	f, err := os.Open("./testdata/base-input.csv")
	if err != nil {
		t.Fatal(err)
	}
	ls, err := ParseBaseCSV(f)
	if err != nil {
		t.Fatal(err)
	}
	ods, err := TransformBaseOrder(ls)
	if err != nil {
		t.Fatal(err)
	}
	for _, o := range ods {
		t.Logf("%+v", o)
	}
}

func TestWriteSummaryFormat(t *testing.T) {
	f, err := os.Open("./testdata/base-input.csv")
	if err != nil {
		t.Fatal(err)
	}
	ls, err := ParseBaseCSV(f)
	if err != nil {
		t.Fatal(err)
	}
	ods, err := TransformBaseOrder(ls)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := WriteSummaryFormat(&buf, ods, "mac"); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", buf.String())
}

func TestWriteClickpostFormat(t *testing.T) {
	f, err := os.Open("./testdata/base-input.csv")
	if err != nil {
		t.Fatal(err)
	}
	ls, err := ParseBaseCSV(f)
	if err != nil {
		t.Fatal(err)
	}
	ods, err := TransformBaseOrder(ls)
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	if err := WriteClickpostFormat(&buf, ods, "mac"); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", buf.String())
}

func TestQuantityFilter(t *testing.T) {
	f, err := os.Open("./testdata/base-input.csv")
	if err != nil {
		t.Fatal(err)
	}
	ls, err := ParseBaseCSV(f)
	if err != nil {
		t.Fatal(err)
	}
	ods, err := TransformBaseOrder(ls)
	if err != nil {
		t.Fatal(err)
	}
	a, b := QuantityFilter(ods, 4)
	t.Logf("A: %+v", a)
	t.Logf("B: %+v", b)
}
