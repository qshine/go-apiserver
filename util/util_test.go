package util

import (
	"testing"
)

// go会自动执行以 Test/Benchmark/Example 开头的测试函数
// go test 命令默认不会执行性能测试函数, 需要指定 -test.bench 参数来运行. 比如 go test -test.bench=".*" 表示执行所有性能测试函数

func TestGenShortId(t *testing.T) {
	shortId, err := GenShortId()
	if shortId == "" || err != nil {
		t.Error("GenShortId failed!")
	}

	t.Log("GenShortId test pass")
}

func BenchmarkGenShortId(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenShortId()
	}
}

func BenchmarkGenShortIdTimeConsuming(b *testing.B) {
	b.StopTimer() //调用该函数停止压力测试的时间计数

	shortId, err := GenShortId()
	if shortId == "" || err != nil {
		b.Error(err)
	}

	b.StartTimer() //重新开始时间

	for i := 0; i < b.N; i++ {
		GenShortId()
	}
}
