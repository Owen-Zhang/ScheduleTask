package main

import (
	"bytes"
	"testing"
)

func BenchmarkAlgorithmOne(b *testing.B) {
	var output bytes.Buffer
	in := []byte("dddddd5454545454545454elviselviselviselviselviselviselviselvisdddddddddd")
	find := []byte("elvis")
	repl := []byte("Elvis")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		output.Reset()
		algOne(in, find, repl, &output)
	}
}

func TestGetArea(t *testing.T) {
	area := GetArea(40, 50)
	if area != 2000 {
		t.Error("测试失败")
	}
}

/*
func TestGetArea2(t *testing.T) {
	area := GetArea(40, 50)
	if area != 3000 {
		t.Error("测试失败")
	}
}
*/

func BenchmarkGetArea(t *testing.B) {

	for i := 0; i < t.N; i++ {
		GetArea(40, 50)
	}
}