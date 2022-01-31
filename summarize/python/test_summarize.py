import unittest
import numpy as np
import pandas as pd
import summarize as summ


class TestMethods(unittest.TestCase):
    def compDF(self, a, b: pd.DataFrame, eps: float):
        if np.linalg.norm(a.subtract(b)) > eps:
            self.fail(f"{a} != {b}")

    def test_init(self):
        want = summ.Spectro
        want.ZPT = 4656.4369
        want.freqs = pd.DataFrame(
            {
                "harms": pd.Series([3943.69, 3833.702, 1650.933]),
                "funds": pd.Series([3753.167, 3656.538, 1598.509]),
                "corrs": pd.Series([3753.1674, 3656.5383, 1598.5086]),
            }
        )
        want.rots = pd.DataFrame(
            [
                [27.6557795, 14.5045043, 9.2632043],
                [26.4993144, 14.4058166, 9.1202295],
                [26.9667697, 14.2851894, 9.0861663],
                [30.2541123, 14.6663643, 9.1168119],
            ],
            columns=["A", "B", "C"],
        )
        #   ],
        #   "Deltas": [
        #     34.9642012678,
        #     761.2706215801,
        #     -149.5172903186,
        #     13.9815699759,
        #     10.9619056097
        #   ],
        #   "Phis": [
        #     13799.0093,
        #     2086925.836,
        #     -88998.91746,
        #     -225200.3817,
        #     6835.370567,
        #     -17508.2843,
        #     338318.6858
        #   ],
        #   "Rhead": [
        #     "r(H1-O2)",
        #     "r(O2-H3)",
        #     "<(O2-H1-H3)"
        #   ],
        #   "Ralpha": [
        #     0.9733425,
        #     0.9733425,
        #     104.2809485
        #   ],
        #   "Requil": [
        #     0.9586139,
        #     0.9586139,
        #     104.401023
        #   ],
        #   "Fermi": null,
        #   "Be": [
        #     27.28099,
        #     14.57684,
        #     9.50051
        #   ],
        #   "Lin": false,
        #   "Imag": false,
        #   "LX": [
        #     3943.69,
        #     3833.7,
        #     1650.93,
        #     30.35,
        #     29.74,
        #     0.09,
        #     0.17,
        #     0.32,
        #     29.61
        #   ]
        # }
        got = summ.Spectro("h2o.json")
        self.assertEqual(got.ZPT, want.ZPT)
        self.compDF(got.freqs, want.freqs, 1e-7)
