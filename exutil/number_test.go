package exutil

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func TestMinMax(t *testing.T) {
	if MaxUint64(0, 0) != 0 {
		t.Fail()
	}
	if MinUint64(^uint64(0), ^uint64(0)) != ^uint64(0) {
		t.Fail()
	}
	if MaxUint64(^uint64(0), 0) != ^uint64(0) {
		t.Fail()
	}
	if MinUint64(^uint64(0), 0) != 0 {
		t.Fail()
	}
	if MaxUint64(1234567890, 0) != 1234567890 {
		t.Fail()
	}
	if MinUint64(^uint64(0), 1234567890) != 1234567890 {
		t.Fail()
	}
	if MinUint64(0, 1234567890) != 0 {
		t.Fail()
	}
}

func TestAddBigInts(t *testing.T) {
	// 10^80 > 2 ^256 ~= 2^77
	res, err := AddBigInt("12345678901234567890123456789012345678901234567890123456789012345678901234567890", "9876543210987654321098765432109876543210987654321098765432109876543210987654321")
	log.Println("Got", res.Text(10))
	if err != nil {
		log.Println(err)
		t.Fail()
	}
	if res.Text(10) != "22222222112222222211222222221122222222112222222211222222221122222222112222222211" {
		t.Fail()
	}
}

func BenchmarkAddBigInts(b *testing.B) {
	AddBigInt("12345678901234567890123456789012345678901234567890123456789012345678901234567890", "9876543210987654321098765432109876543210987654321098765432109876543210987654321")
}

func TestImprFloat_Float(t *testing.T) {
	imf := new(ImprFloat)
	tests := []struct {
		in  float64
		dec int
		out float64
	}{
		{1000, -2, 10},
		{0.1, 2, 10},
		{1000, -10, 0.0000001},
		{0.00001, 7, 100},
		{0, -7, 0},
	}
	for i, test := range tests {
		out := imf.FromFloat(test.in).Shift(test.dec).ToFloat()
		if out != test.out {
			t.Errorf("case %v failed, get %v, expect %v", i, out, test.out)
		}
	}
}
func TestImprFloat_IntString(t *testing.T) {
	imf := new(ImprFloat)
	tests := []struct {
		in  string
		dec int
		out string
		err bool
	}{
		{"0.001", 10, "10000000", false},
		{"100.1", -1, "10", false},
		{"100.1", -5, "0", false},
		{"0.001", 1, "0", false},
		{"1234.567890", 4, "12345678", false},
		{"-1234.567890", 4, "0", true},
	}
	for i, test := range tests {
		out := imf.FromString(test.in).Shift(test.dec).ToIntString()
		if out != test.out {
			t.Errorf("case %v failed, get %v, expect %v", i, out, test.out)
		}
		if (imf.GetLastErr() != nil) != test.err {
			t.Errorf("case %v failed, get error %v, expect error %v", i, imf.GetLastErr(), test.err)
		}
	}
}

func TestImprFloat_TrimDecimal(t *testing.T) {
	imf := new(ImprFloat)
	tests := []struct {
		in  float64
		dec int
		out float64
	}{
		{0, 0, 0},
		{0.1, 2, 10},
		{1000, -10, 0},
		{1000010000.00001, -4, 100001},
		{16.1233434457445645745, 7, 161233434},
	}
	for i, test := range tests {
		out := imf.FromFloat(test.in).Shift(test.dec).TrimDecimal().ToFloat()
		if out != test.out {
			t.Errorf("case %v failed, get %v, expect %v", i, out, test.out)
		}
	}
}

func TestValidFloatStr(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"normal float1", args{"1.1"}, true},
		{"normal float2", args{"1.1000"}, true},
		{"normal float3", args{"0001.1"}, true},
		{"normal float4", args{"-1.1"}, true},
		{"normal float5", args{"-0001.1"}, true},
		{"normal float6", args{"-111.1000"}, true},
		{"normal float7", args{"-111"}, true},
		{"normal float8", args{"000111000"}, true},
		{"normal big float", args{"-000123412412341264675678789670723412314324532463456457345645111.10001232513451432412342341234123643654640"}, true},
		{"wrong float1", args{"abc"}, false},
		{"wrong float2", args{"abc1.1"}, false},
		{"wrong float3", args{"1.1abc"}, false},
		{"wrong float4", args{"1E2"}, false},
		{"wrong float5", args{"-3.4e56"}, false},
		{"wrong float6", args{"abc-56"}, false},
		{"wrong float7", args{"-56abc"}, false},
		{"wrong float8", args{"null"}, false},
		{"wrong float9", args{"NaN"}, false},
		{"wrong float10", args{"FFEEFFFF"}, false},
		{"wrong float11", args{"CDCDCDCDCD"}, false},
		{"wrong float12", args{"烫烫烫烫烫"}, false},
		{"wrong float13", args{".111"}, false},
		{"wrong float14", args{"111."}, false},
		{"wrong float15", args{""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidFloatStr(tt.args.val); got != tt.want {
				t.Errorf("ValidFloatStr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidPositiveFloatStr(t *testing.T) {
	type args struct {
		val string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{"normal float1", args{"1.1"}, true},
		{"normal float2", args{"1.1000"}, true},
		{"normal float3", args{"0001.1"}, true},
		{"normal float4", args{"-1.1"}, false},
		{"normal float5", args{"-0001.1"}, false},
		{"normal float6", args{"-111.1000"}, false},
		{"normal float7", args{"-111"}, false},
		{"normal float8", args{"000111000"}, true},
		{"normal big float", args{"-000123412412341264675678789670723412314324532463456457345645111.10001232513451432412342341234123643654640"}, false},
		{"wrong float1", args{"abc"}, false},
		{"wrong float2", args{"abc1.1"}, false},
		{"wrong float3", args{"1.1abc"}, false},
		{"wrong float4", args{"1E2"}, false},
		{"wrong float5", args{"-3.4e56"}, false},
		{"wrong float6", args{"abc-56"}, false},
		{"wrong float7", args{"-56abc"}, false},
		{"wrong float8", args{"null"}, false},
		{"wrong float9", args{"NaN"}, false},
		{"wrong float10", args{"FFEEFFFF"}, false},
		{"wrong float11", args{"CDCDCDCDCD"}, false},
		{"wrong float12", args{"烫烫烫烫烫"}, false},
		{"wrong float13", args{".111"}, false},
		{"wrong float14", args{"111."}, false},
		{"wrong float15", args{""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidPositiveFloatStr(tt.args.val); got != tt.want {
				t.Errorf("ValidFloatStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestImprFloat_GetLastError(t *testing.T) {
	imf := new(ImprFloat)
	imf.FromString("-1.5").Shift(2).ToIntString()
	imf.FromString("1.5").Shift(2).ToIntString()
	err := imf.GetLastErr()
	if err != nil {
		t.Error(err)
	}
	imf.FromString("-1.5").Shift(2).ToIntString()
	err = imf.GetLastErr()
	if err == nil {
		t.Error("expect error")
	}
	err = imf.GetLastErr()
	if err != nil {
		t.Error("GetLastError not clear err")
	}
}
