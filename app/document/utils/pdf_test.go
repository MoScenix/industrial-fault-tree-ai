package utils

import "testing"

func TestMeasurePDFTextQualityChineseText(t *testing.T) {
	text := "通用型驱动系统故障数据\n电机过热可能由过载、散热不良或轴承故障导致。"
	quality := measurePDFTextQuality(text)
	if !quality.usable() {
		t.Fatalf("expected chinese text to be usable, got %+v", quality)
	}
}

func TestMeasurePDFTextQualityGarbledText(t *testing.T) {
	text := "M3^\u00d0p\u000f\u00a3\u00e9($pl\u00d9\u00e0\n^\u00ed\u00d1\u00c7T\u00a4\u00aa^\u00d0\u000ep23p($6k\u00a4e\u00de\u00a6\u00f7"
	quality := measurePDFTextQuality(text)
	if quality.usable() {
		t.Fatalf("expected garbled text to be rejected, got %+v", quality)
	}
}

func TestNormalizePDFText(t *testing.T) {
	text := " 第一行 \n\n第二\t行 \f 第三行 "
	got := normalizePDFText(text)
	want := "第一行\n第二 行\n第三行"
	if got != want {
		t.Fatalf("normalizePDFText() = %q, want %q", got, want)
	}
}
