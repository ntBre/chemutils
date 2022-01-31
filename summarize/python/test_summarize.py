import unittest
import numpy as np
import pandas as pd
import summarize as summ

pd.options.display.max_columns = None


class TestInit(unittest.TestCase):
    got = summ.Spectro("h2o.json")
    want = summ.Spectro

    def aslist(self, df, x):
        return df[x].tolist()

    def compDF(self, a, b: pd.DataFrame, eps: float):
        "compare two purely numerical DataFrames"
        if np.linalg.norm(a.subtract(b)) > eps:
            self.fail(f"{a} != {b}")

    def test_ZPT(self):
        self.want.ZPT = 4656.4369
        self.assertEqual(self.got.ZPT, self.want.ZPT)

    def test_freqs(self):
        self.want.freqs = pd.DataFrame(
            {
                "harms": pd.Series([3943.69, 3833.702, 1650.933]),
                "funds": pd.Series([3753.167, 3656.538, 1598.509]),
                "corrs": pd.Series([3753.1674, 3656.5383, 1598.5086]),
            }
        )
        self.compDF(self.got.freqs, self.want.freqs, 1e-7)

    def test_rots(self):
        self.want.rots = pd.DataFrame(
            [
                [27.6557795, 14.5045043, 9.2632043],
                [26.4993144, 14.4058166, 9.1202295],
                [26.9667697, 14.2851894, 9.0861663],
                [30.2541123, 14.6663643, 9.1168119],
            ],
            columns=["A", "B", "C"],
        )
        self.compDF(self.got.rots, self.want.rots, 1e-7)

    def test_deltas(self):
        self.want.deltas = pd.DataFrame(
            {
                "Const.": [
                    "$\\Delta_{J}$",
                    "$\\Delta_{K}$",
                    "$\\Delta_{JK}$",
                    "$\\delta_{J}$",
                    "$\\delta_{K}$",
                ],
                "Value": [
                    34.9642012678,
                    761.2706215801,
                    -149.5172903186,
                    13.9815699759,
                    10.9619056097,
                ],
                "Units": ["MHz"] * 5,
            }
        )
        self.assertEqual(
            self.aslist(self.got.deltas, "Value"),
            self.aslist(self.want.deltas, "Value"),
        )
        self.assertEqual(
            self.aslist(self.got.deltas, "Const."),
            self.aslist(self.want.deltas, "Const."),
        )
        self.assertEqual(
            self.aslist(self.got.deltas, "Units"),
            self.aslist(self.want.deltas, "Units"),
        )

    def test_phis(self):
        self.want.phis = pd.DataFrame(
            {
                "Const.": [
                    "$\\Phi_{J}$",
                    "$\\Phi_{K}$",
                    "$\\Phi_{JK}$",
                    "$\\Phi_{KJ}$",
                    "$\\phi_{j}$",
                    "$\\phi_{jk}$",
                    "$\\phi_{k}$",
                ],
                "Value": [
                    13799.0093,
                    2086925.836,
                    -88998.91746,
                    -225200.3817,
                    6835.370567,
                    -17508.2843,
                    338318.6858,
                ],
                "Units": ["Hz"] * 7,
            }
        )
        self.assertEqual(
            self.aslist(self.got.phis, "Value"),
            self.aslist(self.want.phis, "Value"),
        )
        self.assertEqual(
            self.aslist(self.got.phis, "Const."),
            self.aslist(self.want.phis, "Const."),
        )
        self.assertEqual(
            self.aslist(self.got.phis, "Units"),
            self.aslist(self.want.phis, "Units"),
        )

    def test_rest(self):
        self.assertEqual(
            self.got.Rhead,
            ["r(H1-O2)", "r(O2-H3)", "<(O2-H1-H3)"],
        )
        self.assertEqual(
            self.got.Ralpha,
            [0.9733425, 0.9733425, 104.2809485],
        )
        self.assertEqual(
            self.got.Requil,
            [0.9586139, 0.9586139, 104.401023],
        )
        self.assertEqual(self.got.fermi, None)
        self.assertEqual(
            self.got.Be,
            [27.28099, 14.57684, 9.50051],
        )
        self.assertEqual(self.got.Lin, False)
        self.assertEqual(self.got.Imag, False)
        self.assertEqual(
            self.got.LX,
            [3943.69, 3833.7, 1650.93, 30.35, 29.74, 0.09, 0.17, 0.32, 29.61],
        )
