package main

import (
	"testing"

	res "k8s.io/apimachinery/pkg/api/resource"
)

// func Test_Quantity(t *testing.T) {
// 	q := res.NewScaledQuantity(res.MaxMilliValue, res.Nano)
// 	t.Log(q)
// 	q2, precise := q.AsScale(res.Kilo)
// 	t.Log(q2, precise)
// }

func Test_CmpWTF(t *testing.T) {
	for _, tt := range []struct {
		x, y string
	}{
		{"7816Mi", "7872Mi"},
		{"15624744Ki", "15522Mi"},
		{"15500m", "15580m"},
	} {
		x := res.MustParse(tt.x)
		y := res.MustParse(tt.y)
		if !((&x).Cmp(y) < 0) {
			t.Errorf("%s should be less than %s\n%v\n%v", &x, &y, x, y)
		}
	}
}

func Test_Summarize(t *testing.T) {
	for _, tt := range []struct {
		val, cap string
		pct      int64
	}{
		{"7170m", "7910m", 91},
		{"7353m", "7910m", 93},
		{"7558m", "7910m", 96},
		{"7870m", "7910m", 99},

		{"6408Mi", "27219332Ki", 24},
		{"15624744Ki", "27219332Ki", 57},
		{"15522Mi", "27219332Ki", 58},
		{"3830Mi", "27219332Ki", 14},
	} {
		val := res.MustParse(tt.val)
		cap := res.MustParse(tt.cap)
		if res := percent(&val, &cap); tt.pct != res {
			t.Errorf("%v not equal to %v\n%v\n%v", res, tt.pct, val, cap)
		}
	}
}
