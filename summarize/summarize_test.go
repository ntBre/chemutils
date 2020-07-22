package summarize

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestSpectro(t *testing.T) {
	tests := []struct {
		file   string
		nfreq  int
		zpt    float64
		harm   []float64
		fund   []float64
		corr   []float64
		rotABC [][]float64
		deltas []float64
		phis   []float64
		requil []float64
		ralpha []float64
		rhead  []string
		fermi  []string
	}{
		{
			file:  "testfiles/spectro.out",
			nfreq: 6,
			zpt:   4682.2527,
			harm: []float64{3811.360, 2337.700, 1267.577,
				1086.351, 496.788, 437.756},
			fund: []float64{3623.015, 2299.805, 1231.309,
				1081.661, 513.228, 454.579},
			corr: []float64{3623.0149, 2299.8053, 1231.3094,
				1081.6611, 513.2276, 454.5787},
			rotABC: [][]float64{
				{0.3533242, 0.3473852, 22.5883184},
				{0.3531433, 0.3469946, 21.5758850},
				{0.3508629, 0.3449969, 22.5509605},
				{0.3536392, 0.3472748, 23.9984685},
				{0.3517191, 0.3456623, 22.5514979},
				{0.3538316, 0.3484810, 53.7297798},
				{0.3547570, 0.3480413, -8.6483579},
			},
			deltas: []float64{
				0.0041596072, 276.2016104107,
				0.6722227103, 0.0000596455,
				0.3035637199,
			},
			phis: []float64{
				-0.4946484183E-03, 0.2374310264E+06,
				0.3252182153E+01, -0.2896689993E+04,
				0.4401912504E-04, 0.1940400502E+01,
				0.6140585605E+04,
			},
			requil: []float64{
				1.1573804, 1.2990232, 0.9626153,
				109.6544076, 41.1568277,
			},
			ralpha: []float64{
				1.1586186, 1.3025842, 0.9724295,
				109.7346599, 40.9623102,
			},
			rhead: []string{
				"r(N1-C2)", "r(C2-O3)", "r(O3-H4)",
				"<(O3-C2-H4)", "<(H4-O3-C2)",
			},
			fermi: []string{
				"2v_4 = v_4+v_3 = v_2",
				"2v_5 = v_4",
			},
		},
		{
			file:  "testfiles/bn.spectro.out",
			nfreq: 18,
			zpt:   15183.1168,
			harm: []float64{
				3589.676, 3589.544, 3474.134,
				2503.147, 2503.062, 2445.411,
				1678.230, 1678.002, 1339.900,
				1200.972, 1200.887, 1199.250,
				1075.609, 1075.586, 678.492,
				647.259, 647.219, 259.360,
			},
			fund: []float64{
				3415.154, 3413.695, 3307.093,
				2392.656, 2391.767, 2337.922,
				1620.244, 1624.482, 1286.615,
				1174.141, 1167.334, 1176.662,
				1042.560, 1045.505, 631.003,
				658.558, 638.619, 299.717,
			},
			corr: []float64{
				3415.1544, 3413.6946, 3333.2029,
				2395.8613, 2405.2258, 2274.5058,
				1608.5485, 1611.7732, 1286.6151,
				1172.2846, 1165.4915, 1173.5409,
				1044.9090, 1048.7250, 643.5182,
				658.5578, 638.6192, 299.7165,
			},
			rotABC: [][]float64{
				{0.5811210, 0.5810957, 2.4335266},
				{0.5809882, 0.5809551, 2.7705867},
				{0.5809804, 0.5809642, 2.0858824},
				{0.5808512, 0.5808255, 2.4235207},
				{0.5818595, 0.5817493, 3.8384803},
				{0.5817740, 0.5818355, 1.0061189},
				{0.5814376, 0.5814113, 2.4190953},
				{0.5814054, 0.5830474, 2.4245754},
				{0.5830710, 0.5813788, 2.4245757},
				{0.5774825, 0.5774588, 2.4337605},
				{0.5812362, 0.5816348, 18.1744941},
				{0.5812952, 0.5808612, math.NaN()},
				{0.5793261, 0.5792690, 2.4418074},
				{0.5783859, 0.5776643, 2.4387711},
				{0.5776692, 0.5783597, 2.4387689},
				{0.5713277, 0.5713028, 2.4304535},
				{0.5796127, 0.5777998, 26.7853636},
				{0.5778180, 0.5795809, math.NaN()},
				{0.5775190, 0.5774985, 2.4322943},
			},
			deltas: []float64{
				0.0469588942, 0.2228671420,
				0.0594222941, 0.0000002512,
				0.2123725433,
			},
			phis: []float64{
				-0.1633979197E+00, -0.1188042331E+04,
				-0.5090171025E+03, 0.1700920469E+04,
				-0.3176907240E-04, 0.4680895674E+00,
				-0.4161588292E+09,
			},
			requil: []float64{
				1.0135437, 1.0135427, 1.0135427,
				1.6495880, 1.2092633, 1.2092631,
				1.2092631, 110.9864784, 110.9777107,
				110.9777107, 104.8600368, 104.8602644,
				104.8602644,
			},
			ralpha: []float64{
				1.0139841, 1.0139881, 1.0139932,
				1.6730848, 1.2168612, 1.2168475,
				1.2168590, 110.9798368, 110.9703946,
				110.9730406, 104.7957009, 104.7957520,
				104.7942464,
			},
			rhead: []string{
				"r(N1-H2)", "r(N1-H3)", "r(N1-H4)",
				"r(N1-B5)", "r(B5-H6)", "r(B5-H7)",
				"r(B5-H8)",
				"<(N1-H2-B5)", "<(N1-H3-B5)", "<(N1-H4-B5)",
				"<(B5-H6-N1)", "<(B5-H7-N1)", "<(B5-H8-N1)",
			},
			fermi: []string{
				"v_16+v_15 = v_17+v_16 = v_10",
				"2v_16 = 2v_17 = v_17+v_15 = v_11",
				"2v_15 = v_12",
				"v_18+v_17 = v_13",
				"v_18+v_16 = v_14",
				"2v_18 = v_15",
				"2v_7 = 2v_8 = v_3",
				"2v_10 = 2v_11 = v_12+v_11 = v_4",
				"v_11+v_10 = v_12+v_10 = v_5",
				"2v_10 = 2v_11 = 2v_12 = v_13+v_10 = v_14+v_11 = v_6",
				"v_15+v_13 = v_16+v_14 = v_17+v_13 = v_7",
				"v_15+v_14 = v_16+v_13 = v_17+v_14 = v_8",
				"2v_15 = 2v_16 = 2v_17 = v_9",
			},
		},
		{
			file:  "testfiles/jax.prob.out",
			nfreq: 9,
			zpt:   6847.4520,
			harm: []float64{
				3819.706, 3478.907, 2243.281,
				1272.576, 1069.500, 622.746,
				537.802, 401.606, 368.267,
			},
			fund: []float64{
				3628.429, 3345.921, 2193.109,
				1235.865, 1057.343, 629.176,
				548.446, 417.098, 386.093,
			},
			corr: []float64{
				3628.4295, 3347.8029, 2203.3286,
				1235.8645, 1052.7773, 629.1764,
				548.4458, 417.0980, 386.0932,
			},
			rotABC: [][]float64{
				{0.3228534, 0.3178922, 22.4479576},
				{0.3227327, 0.3176009, 21.4728179},
				{0.3220040, 0.3170685, 22.4445304},
				{0.3209536, 0.3160462, 22.4205485},
				{0.3230913, 0.3177488, 23.7354379},
				{0.3215155, 0.3164633, 22.4077827},
				{0.3233114, 0.3179611, 23.3635186},
				{0.3230655, 0.3185377, 22.3819206},
				{0.3233493, 0.3189513, 21.3237610},
				{0.3241649, 0.3184727, 22.6271116},
			},
			deltas: []float64{
				0.0032166547, 253.0816116693,
				0.6221602357, 0.0000416412,
				0.2759818046,
			},
			phis: []float64{
				-0.3148755392E-03, 0.1450201513E+06,
				0.3478976071E+01, -0.5278691516E+04,
				0.2327941810E-04, 0.1937446213E+01,
				0.6871318174E+04,
			},
			requil: []float64{
				1.0609763, 1.2034593, 1.3145119,
				0.9620542, 41.6792532,
			},
			ralpha: []float64{
				1.0492227, 1.2063767, 1.3169088,
				0.9718421, 41.4652508,
			},
			rhead: []string{
				"r(H1-C2)", "r(C2-C3)", "r(C3-O4)",
				"r(O4-H5)", "<(H5-O4-C3)",
			},
			fermi: []string{
				"v_5+v_3 = v_2",
				"2v_5 = v_5+v_4 = v_3",
				"2v_6 = 2v_7 = v_5",
			},
		},
	}
	for _, test := range tests {
		gzpt, gharm, gfund, gcorr,
			rotABC, deltas, phis,
			requil, ralpha, rhead,
			fermi := Spectro(test.file, test.nfreq)
		if gzpt != test.zpt {
			t.Errorf("Spectro(%s, %d): got %f, wanted %f\n",
				test.file, test.nfreq, gzpt, test.zpt)
		}
		if !reflect.DeepEqual(gharm, test.harm) {
			t.Errorf("Spectro(%s, %d): got %v, wanted %v\n",
				test.file, test.nfreq, gharm, test.harm)
		}
		if !reflect.DeepEqual(gfund, test.fund) {
			t.Errorf("Spectro(%s, %d): got\n%v, wanted\n%v\n",
				test.file, test.nfreq, gfund, test.fund)
		}
		if !reflect.DeepEqual(gcorr, test.corr) {
			t.Errorf("Spectro(%s, %d): got\n%v,\nwanted\n%v\n",
				test.file, test.nfreq, gcorr, test.corr)
		}
		if !reflect.DeepEqual(rotABC, test.rotABC) {
			if len(rotABC) != len(test.rotABC) {
				t.Errorf("Spectro(%s, %d): length mismatch on ABCs, got %d, wanted %d\n",
					test.file, test.nfreq, len(rotABC), len(test.rotABC))
			}
			for i := range rotABC {
				if !reflect.DeepEqual(rotABC[i], test.rotABC[i]) {
					for j := range rotABC[i] {
						if (rotABC[i][j] != test.rotABC[i][j]) &&
							(!math.IsNaN(rotABC[i][j]) ||
								!math.IsNaN(test.rotABC[i][j])) {
							fmt.Println(rotABC[i][j], test.rotABC[i][j])
							t.Errorf("Spectro(%s, %d): different ABCs\n",
								test.file, test.nfreq)
							fmt.Printf("i: %d\ng: %v\nw: %v\n",
								i, rotABC[i], test.rotABC[i])
						}
					}
				}
			}
		}
		if !reflect.DeepEqual(deltas, test.deltas) {
			t.Errorf("Spectro(%s, %d): got %v, wanted %v\n",
				test.file, test.nfreq, deltas, test.deltas)
		}
		if !reflect.DeepEqual(phis, test.phis) {
			t.Errorf("Spectro(%s, %d): got %v, wanted %v\n",
				test.file, test.nfreq, phis, test.phis)
		}
		if !reflect.DeepEqual(requil, test.requil) {
			t.Errorf("Spectro(%s, %d): got %v, wanted %v\n",
				test.file, test.nfreq, requil, test.requil)
		}
		if !reflect.DeepEqual(ralpha, test.ralpha) {
			t.Errorf("Spectro(%s, %d): got %v, wanted %v\n",
				test.file, test.nfreq, ralpha, test.ralpha)
		}
		if !reflect.DeepEqual(rhead, test.rhead) {
			t.Errorf("Spectro(%s, %d): got %v, wanted %v\n",
				test.file, test.nfreq, rhead, test.rhead)
		}
		if !reflect.DeepEqual(fermi, test.fermi) {
			t.Errorf("Spectro(%s, %d): got %q, wanted %q\n",
				test.file, test.nfreq, fermi, test.fermi)
		}
	}
}
