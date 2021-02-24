package summarize

import (
	"fmt"
	"io/ioutil"
	"math"
	"reflect"
	"testing"
)

func TestSpectro(t *testing.T) {
	tests := []struct {
		file string
		res  Result
	}{
		{
			file: "testfiles/long.spectro.out",
			res: Result{
				ZPT: 11096.158700,
				Harm: []float64{3248.349, 3222.482, 3155.97,
					3140.984, 1672.172, 1476.484,
					1368.466, 1247.682, 1049.892,
					963.528, 950.665, 824.895},
				Fund: []float64{3106.262, 3092.351, 3009.299,
					2740.651, 1647.52, 1451.969,
					1370.597, 1241.679, 1049.265,
					1058.938, 921.156, 916.467},
				Corr: []float64{3106.2618, 3092.3506, 3009.2991,
					2740.6511, 1647.5204, 1451.9686,
					1370.5967, 1241.6791, 1049.265,
					1058.938, 921.1556, 916.4675},
				Rots: [][]float64{
					{0.9968289, 0.8246127, 4.8478362},
					{0.9945645, 0.8226369, 4.8227033},
					{0.9947128, 0.8229659, 4.8215817},
					{0.9953656, 0.8224301, 4.8059433},
					{0.9956388, 0.8226594, 4.8129777},
					{0.9939572, 0.8198069, 4.8612817},
					{1.0032013, 0.8230419, 4.9123915},
					{0.9977388, 0.8255303, 4.8597252},
					{0.998367, 0.8195055, 4.9850938},
					{1.0086962, 0.8246204, 4.7581942},
					{0.9913069, 0.8259753, 4.9864229},
					{0.9912475, 0.825292, 4.7204069},
					{0.9819411, 0.8230453, 4.7241748},
				},
				Deltas: []float64{
					0.0432422353, 2.4620967217, 0.3039591196,
					0.0081676141, 0.2863179597,
				},
				Phis: []float64{
					0.06856144277, 169.7072549, 5.2465456,
					-12.08390707, 0.03225995931, 3.043895151,
					98.29430518,
				},
				Requil: []float64{
					1.0826038, 1.0826038, 1.0826038,
					1.0826038, 1.333685, 25.9538557,
					25.9538557, 25.9538557, 25.9538557,
				},
				Ralpha: []float64{
					1.0901207, 1.0901058, 1.0901174,
					1.090135, 1.340629, 25.9776209,
					25.9990334, 25.9802257, 25.996323,
				},
				Rhead: []string{
					"r(C1-H3)", "r(C1-H4)", "r(C2-H5)",
					"r(C2-H6)", "r(C1-C2)", "<(C1-C2-H5)",
					"<(C1-C2-H6)", "<(C2-C1-H3)", "<(C2-C1-H4)",
				},
				Fermi: []string{
					"2v_5=v_7+v_5=v_3",
					"v_6+v_5=v_4",
					"2v_12=v_5"},
			},
		},
		{
			file: "testfiles/spectro.out",
			res: Result{
				ZPT: 4682.2527,
				Harm: []float64{
					3811.360, 2337.700, 1267.577,
					1086.351, 496.788, 437.756,
				},
				Fund: []float64{
					3623.015, 2299.805, 1231.309,
					1081.661, 513.228, 454.579,
				},
				Corr: []float64{
					3623.0149, 2299.8053, 1231.3094,
					1081.6611, 513.2276, 454.5787,
				},
				Rots: [][]float64{
					{0.3533242, 0.3473852, 22.5883184},
					{0.3531433, 0.3469946, 21.5758850},
					{0.3508629, 0.3449969, 22.5509605},
					{0.3536392, 0.3472748, 23.9984685},
					{0.3517191, 0.3456623, 22.5514979},
					{0.3538316, 0.3484810, 53.7297798},
					{0.3547570, 0.3480413, -8.6483579},
				},
				Deltas: []float64{
					0.0041596072, 276.2016104107,
					0.6722227103, 0.0000596455,
					0.3035637199,
				},
				Phis: []float64{
					-0.4946484183E-03, 0.2374310264E+06,
					0.3252182153E+01, -0.2896689993E+04,
					0.4401912504E-04, 0.1940400502E+01,
					0.6140585605E+04,
				},
				Requil: []float64{
					1.1573804, 1.2990232, 0.9626153,
					109.6544076, 41.1568277,
				},
				Ralpha: []float64{
					1.1586186, 1.3025842, 0.9724295,
					109.7346599, 40.9623102,
				},
				Rhead: []string{
					"r(N1-C2)", "r(C2-O3)", "r(O3-H4)",
					"<(O3-C2-H4)", "<(H4-O3-C2)",
				},
				Fermi: []string{
					"2v_4=v_4+v_3=v_2",
					"2v_5=v_4",
				},
			},
		},
		{
			file: "testfiles/bn.spectro.out",
			res: Result{
				ZPT: 15183.1168,
				Harm: []float64{
					3589.676, 3589.544, 3474.134,
					2503.147, 2503.062, 2445.411,
					1678.230, 1678.002, 1339.900,
					1200.972, 1200.887, 1199.250,
					1075.609, 1075.586, 678.492,
					647.259, 647.219, 259.360,
				},
				Fund: []float64{
					3415.154, 3413.695, 3307.093,
					2392.656, 2391.767, 2337.922,
					1620.244, 1624.482, 1286.615,
					1174.141, 1167.334, 1176.662,
					1042.560, 1045.505, 631.003,
					658.558, 638.619, 299.717,
				},
				Corr: []float64{
					3415.1544, 3413.6946, 3333.2029,
					2395.8613, 2405.2258, 2274.5058,
					1608.5485, 1611.7732, 1286.6151,
					1172.2846, 1165.4915, 1173.5409,
					1044.9090, 1048.7250, 643.5182,
					658.5578, 638.6192, 299.7165,
				},
				Rots: [][]float64{
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
				Deltas: []float64{
					0.0469588942, 0.2228671420,
					0.0594222941, 0.0000002512,
					0.2123725433,
				},
				Phis: []float64{
					-0.1633979197E+00, -0.1188042331E+04,
					-0.5090171025E+03, 0.1700920469E+04,
					-0.3176907240E-04, 0.4680895674E+00,
					-0.4161588292E+09,
				},
				Requil: []float64{
					1.0135437, 1.0135427, 1.0135427,
					1.6495880, 1.2092633, 1.2092631,
					1.2092631, 110.9864784, 110.9777107,
					110.9777107, 104.8600368, 104.8602644,
					104.8602644,
				},
				Ralpha: []float64{
					1.0139841, 1.0139881, 1.0139932,
					1.6730848, 1.2168612, 1.2168475,
					1.2168590, 110.9798368, 110.9703946,
					110.9730406, 104.7957009, 104.7957520,
					104.7942464,
				},
				Rhead: []string{
					"r(N1-H2)", "r(N1-H3)", "r(N1-H4)",
					"r(N1-B5)", "r(B5-H6)", "r(B5-H7)",
					"r(B5-H8)",
					"<(N1-H2-B5)", "<(N1-H3-B5)", "<(N1-H4-B5)",
					"<(B5-H6-N1)", "<(B5-H7-N1)", "<(B5-H8-N1)",
				},
				Fermi: []string{
					"v_16+v_15=v_17+v_16=v_10",
					"2v_16=2v_17=v_17+v_15=v_11",
					"2v_15=v_12",
					"v_18+v_17=v_13",
					"v_18+v_16=v_14",
					"2v_18=v_15",
					"2v_7=2v_8=v_3",
					"2v_10=2v_11=v_12+v_11=v_4",
					"v_11+v_10=v_12+v_10=v_5",
					"2v_10=2v_11=2v_12=v_13+v_10=v_14+v_11=v_6",
					"v_15+v_13=v_16+v_14=v_17+v_13=v_7",
					"v_15+v_14=v_16+v_13=v_17+v_14=v_8",
					"2v_15=2v_16=2v_17=v_9",
				},
			},
		},
		{
			file: "testfiles/jax.prob.out",
			res: Result{
				ZPT: 6847.4520,
				Harm: []float64{
					3819.706, 3478.907, 2243.281,
					1272.576, 1069.500, 622.746,
					537.802, 401.606, 368.267,
				},
				Fund: []float64{
					3628.429, 3345.921, 2193.109,
					1235.865, 1057.343, 629.176,
					548.446, 417.098, 386.093,
				},
				Corr: []float64{
					3628.4295, 3347.8029, 2203.3286,
					1235.8645, 1052.7773, 629.1764,
					548.4458, 417.0980, 386.0932,
				},
				Rots: [][]float64{
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
				Deltas: []float64{
					0.0032166547, 253.0816116693,
					0.6221602357, 0.0000416412,
					0.2759818046,
				},
				Phis: []float64{
					-0.3148755392E-03, 0.1450201513E+06,
					0.3478976071E+01, -0.5278691516E+04,
					0.2327941810E-04, 0.1937446213E+01,
					0.6871318174E+04,
				},
				Requil: []float64{
					1.0609763, 1.2034593, 1.3145119,
					0.9620542, 41.6792532,
				},
				Ralpha: []float64{
					1.0492227, 1.2063767, 1.3169088,
					0.9718421, 41.4652508,
				},
				Rhead: []string{
					"r(H1-C2)", "r(C2-C3)", "r(C3-O4)",
					"r(O4-H5)", "<(H5-O4-C3)",
				},
				Fermi: []string{
					"v_5+v_3=v_2",
					"2v_5=v_5+v_4=v_3",
					"2v_6=2v_7=v_5",
				},
			},
		},
		{
			file: "testfiles/nofermi2.out",
			res: Result{
				ZPT: 6974.5686,
				Harm: []float64{
					3281.244, 3247.542, 1623.324,
					1307.596, 1090.695, 992.978,
					908.490, 901.527, 785.379,
				},
				Fund: []float64{
					3140.115, 3113.252, 1589.153,
					1273.128, 1059.592, 967.523,
					887.009, 845.919, 769.511,
				},
				Corr: []float64{
					3128.0166, 3113.2520, 1590.4451,
					1273.1281, 1059.5924, 967.5230,
					887.0092, 845.9192, 769.5109,
				},
				Rots: [][]float64{
					[]float64{1.0938807, 0.5645509, 1.1728388},
					[]float64{1.1019479, 0.5662932, 1.1708971},
					[]float64{1.1015505, 0.5662303, 1.1709113},
					[]float64{1.1082936, 0.5680195, 1.1731129},
					[]float64{1.0962916, 0.5658264, 1.1697352},
					[]float64{1.0971252, 0.5605153, 1.1683363},
					[]float64{1.0890950, 0.5648179, 1.1737299},
					[]float64{1.0942300, 0.5618767, 1.1711949},
					[]float64{1.0935530, 0.5641348, 1.1748108},
					[]float64{1.0877676, 0.5637609, 1.1683857},
				},
				Deltas: []float64{
					0.0452918791, -0.0314863134,
					0.1330975467, 0.0140192418,
					0.0763005051,
				},
				Phis: []float64{
					0.2976480237E+00, 0.5011295115E+01, 0.1221389490E+01,
					-0.5534684999E+01, 0.2248273880E+00, 0.1403699085E+01,
					-0.7512823561E-01,
				},
				Requil: []float64{
					1.3215820, 1.4160153, 1.4160153,
					1.0748590, 1.0748590, 147.7726162,
					147.7726162,
				},
				Ralpha: []float64{
					1.3151069, 1.4204888, 1.4204890,
					1.0743545, 1.0743537, 147.6296447,
					147.6296238,
				},
				Rhead: []string{
					"r(C2-C3)", "r(C1-C2)", "r(C1-C3)",
					"r(C2-H4)", "r(C3-H5)", "<(C2-H4-C1)",
					"<(C3-H5-C1)",
				},
				Fermi: []string{
					"2v_3=v_1",
					"2v_7=2v_8=2v_9=v_3",
				},
			},
		},
		{
			file: "testfiles/degen.out",
			res: Result{
				ZPT: 5772.6827,
				Harm: []float64{
					3503.373, 3411.176, 2008.488,
					749.015, 616.381,
				},
				Fund: []float64{
					3366.417, 3281.873, 1972.166,
					736.705, 614.136,
				},
				Corr: []float64{
					3366.4175, 3281.8730, 1972.1663,
					736.7052, 614.1360,
				},
				Rots: [][]float64{
					[]float64{-0.0058963, 0.0000000, -0.0058963},
					[]float64{-0.0127827, 0.0000000, -0.0127827},
					[]float64{-0.0117556, 0.0000000, -0.0117556},
					[]float64{-0.0119290, 0.0000000, -0.0119290},
					[]float64{-0.0037384, 0.0000000, -0.0037384},
					[]float64{-0.0045613, 0.0000000, -0.0045613},
				},
				Requil: []float64{
					1.0631185, 1.2054552, 1.0631185,
				},
				Ralpha: []float64{
					1.0578407, 1.2113454, 1.0578407,
				},
				Rhead: []string{
					"r(H1-C2)", "r(C2-C3)", "r(C3-H4)",
				},
			},
		},
	}
	for _, test := range tests {
		got := *Spectro(test.file)
		if !reflect.DeepEqual(got, test.res) {
			if !reflect.DeepEqual(got.ZPT, test.res.ZPT) {
				t.Errorf("Spectro(%s): got %v, wanted %v\n",
					test.file, got.ZPT, test.res.ZPT)
			}
			if !reflect.DeepEqual(got.Harm, test.res.Harm) {
				t.Errorf("got %v, wanted %v\n", got.Harm, test.res.Harm)
			}
			if !reflect.DeepEqual(got.Fund, test.res.Fund) {
				t.Error("fund")
			}
			if !reflect.DeepEqual(got.Corr, test.res.Corr) {
				t.Errorf("got %v, wanted %v\n", got.Corr, test.res.Corr)
			}
			if !reflect.DeepEqual(got.Rots, test.res.Rots) {
				// check manually to handle NaN
				if len(got.Rots) != len(test.res.Rots) {
					t.Errorf("got %v, wanted %v\n", got.Rots, test.res.Rots)
				} else {
					for i := range got.Rots {
						for j := range got.Rots[i] {
							// not equal and both are non-NaN
							if got.Rots[i][j] !=
								test.res.Rots[i][j] &&
								!(math.IsNaN(got.Rots[i][j]) &&
									math.IsNaN(test.res.
										Rots[i][j])) {
								t.Errorf("got %v, wanted %v\n",
									got.Rots, test.res.Rots)
							}
						}
					}
				}
			}
			if !reflect.DeepEqual(got.Deltas, test.res.Deltas) {
				t.Error("deltas")
			}
			if !reflect.DeepEqual(got.Phis, test.res.Phis) {
				t.Error("phis")
			}
			if !reflect.DeepEqual(got.Rhead, test.res.Rhead) {
				t.Error("rhead")
			}
			if !reflect.DeepEqual(got.Ralpha, test.res.Ralpha) {
				t.Error("ralpha")
			}
			if !reflect.DeepEqual(got.Requil, test.res.Requil) {
				t.Error("requil")
			}
			if !reflect.DeepEqual(got.Fermi, test.res.Fermi) {
				t.Errorf("fermi:\ngot %v, wanted %v\n",
					got.Fermi, test.res.Fermi)
			}
		}
	}
}

func TestIntder(t *testing.T) {
	tests := []struct {
		infile string
		want   Intder
	}{
		{
			infile: "testfiles/intder.out",
			want: Intder{
				Geom: []Atom{
					{"C", 0.0000000000, 0.0000000000, -0.8888094004},
					{"C", 0.0000000000, 0.6626968171, 0.3682892206},
					{"C", 0.0000000000, -0.6626968171, 0.3682892206},
					{"H", 0.0000000000, 1.5951938489, 0.9069605214},
					{"H", 0.0000000000, -1.5951938489, 0.9069605214},
				},
				SiIC: [][]int{
					{1, 2, -1, -1, 0},
					{0, 1, -1, -1, 0},
					{0, 2, -1, -1, 0},
					{1, 3, -1, -1, 0},
					{2, 4, -1, -1, 0},
					{3, 1, 0, -1, 1},
					{4, 2, 0, -1, 1},
					{3, 1, 0, 2, 2},
					{4, 2, 0, 1, 2},
				},
				SyIC: [][]int{
					{0},
					{1, 2},
					{3, 4},
					{5, 6},
					{1, -2},
					{3, -4},
					{5, -6},
					{7, -8},
					{7, 8},
				},
				Freq: []float64{
					785.1, 901.7, 908.6,
					992.8, 1090.6, 1307.4,
					1623.6, 3247.6, 3281.4,
				},
				Vibs: `1.000S_{8}
0.809S_{4}+0.133S_{2}-0.057S_{1}
0.810S_{7}+0.189S_{5}
1.000S_{9}
0.809S_{5}-0.190S_{7}
0.686S_{2}-0.187S_{4}-0.126S_{1}
0.793S_{1}+0.174S_{2}
0.998S_{6}
0.971S_{3}
`,
			},
		},
		{
			infile: "testfiles/mason.out",
			want: Intder{
				Geom: []Atom{
					{"He", 0.0000000000, 0.0000000000, -1.7679827066},
					{"H", 0.0000000000, 0.0000000000, -0.5307856662},
					{"H", 0.0000000000, 0.0000000000, 0.5307856678},
					{"He", 0.0000000000, 0.0000000000, 1.7679827049},
				},
				Dumm: []Atom{
					{"X", 1.1111111111, 0.0000000000, -1.0030394677},
					{"X", 0.0000000000, 1.1111111111, -1.0030394677},
					{"X", 1.1111111111, 0.0000000000, 1.0030394709},
					{"X", 0.0000000000, 1.1111111111, 1.0030394709},
				},
				SiIC: [][]int{
					{0, 1, -1, -1, STRE},
					{1, 2, -1, -1, STRE},
					{2, 3, -1, -1, STRE},
					{0, 1, 2, 4, LIN1},
					{3, 2, 1, 6, LIN1},
					{0, 1, 2, 5, LIN1},
					{3, 2, 1, 7, LIN1},
				},
				SyIC: [][]int{
					{0},
					{1},
					{2},
					{3},
					{4},
					{5},
					{6},
				},
				Freq: []float64{
					228.4, 228.4, 280.1,
					398.6, 705.7, 705.7,
					2313.8,
				},
				Vibs: `0.500S_{4}-0.500S_{5}
0.500S_{6}-0.500S_{7}
0.500S_{3}-0.500S_{1}
0.455S_{1}+0.455S_{3}+0.090S_{2}
0.500S_{7}+0.500S_{6}
0.500S_{5}+0.500S_{4}
0.910S_{2}-0.045S_{1}-0.045S_{3}
`,
			},
		},
	}
	for _, test := range tests {
		got := ReadIntder(test.infile)
		if !reflect.DeepEqual(got, test.want) {
			gfile := "/tmp/got.txt"
			wfile := "/tmp/want.txt"
			ioutil.WriteFile(gfile, []byte(got.String()), 0755)
			ioutil.WriteFile(wfile, []byte(test.want.String()), 0755)
			fmt.Printf("(diff %q %q)\n", gfile, wfile)
			if !reflect.DeepEqual(got.Geom, test.want.Geom) {
				t.Errorf("got\n%v, test.wanted\n%v\n", got.Geom, test.want.Geom)
			}
			if !reflect.DeepEqual(got.Dumm, test.want.Dumm) {
				t.Errorf("got\n%v, test.wanted\n%v\n", got.Dumm, test.want.Dumm)
			}
			if !reflect.DeepEqual(got.SiIC, test.want.SiIC) {
				t.Errorf("got\n%v, test.wanted\n%v\n", got.SiIC, test.want.SiIC)
			}
			if !reflect.DeepEqual(got.SyIC, test.want.SyIC) {
				t.Errorf("got\n%v, test.wanted\n%v\n", got.SyIC, test.want.SyIC)
			}
			if !reflect.DeepEqual(got.Freq, test.want.Freq) {
				t.Errorf("got\n%v, test.wanted\n%v\n", got.Freq, test.want.Freq)
			}
			if !reflect.DeepEqual(got.Vibs, test.want.Vibs) {
				t.Errorf("got\n%#+v, test.wanted\n%#+v\n", got.Vibs,
					test.want.Vibs)
			}
		}
	}
}
