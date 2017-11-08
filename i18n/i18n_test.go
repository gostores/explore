/*=================================
* Copyright(c)2017 gostores
* From: github.com/Unknwon/i18n
*=================================*/
package i18n

import (
	"testing"
)

func Benchmark_Tr(b *testing.B) {
	SetMessage("en-US", "testdata/en-US.ini")
	for i := 0; i < b.N; i++ {
		Tr("en-US", "NAME")
	}
}

func Benchmark_TrWithFormat(b *testing.B) {
	SetMessage("en-US", "testdata/en-US.ini")
	for i := 0; i < b.N; i++ {
		Tr("en-US", "SECONDS", 10)
	}
}
