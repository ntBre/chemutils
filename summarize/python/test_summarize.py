import unittest
import numpy as np
import pandas as pd
import summarize as summ

pd.options.display.max_columns = None

# TODO abstract this series of tests, I need the same ones in degen


class TestInit(unittest.TestCase):
    got = summ.Spectro("h2o.json")
    want = summ.new_spectro(
        zpt=4656.4369,
        harms=[3943.69, 3833.702, 1650.933],
        funds=[3753.167, 3656.538, 1598.509],
        corrs=[3753.1674, 3656.5383, 1598.5086],
        rots=[
            [27.6557795, 14.5045043, 9.2632043],
            [26.4993144, 14.4058166, 9.1202295],
            [26.9667697, 14.2851894, 9.0861663],
            [30.2541123, 14.6663643, 9.1168119],
        ],
        deltas=[
            34.9642012678,
            761.2706215801,
            -149.5172903186,
            13.9815699759,
            10.9619056097,
        ],
        phis=[
            13799.0093,
            2086925.836,
            -88998.91746,
            -225200.3817,
            6835.370567,
            -17508.2843,
            338318.6858,
        ],
        rhead=["r(H1-O2)", "r(O2-H3)", "<(O2-H1-H3)"],
        ralpha=[0.9733425, 0.9733425, 104.2809485],
        requil=[0.9586139, 0.9586139, 104.401023],
        fermi=None,
        be=[27.28099, 14.57684, 9.50051],
        lin=False,
        imag=False,
        lx=[3943.69, 3833.7, 1650.93, 30.35, 29.74, 0.09, 0.17, 0.32, 29.61],
    )

    def aslist(self, df, x):
        return df[x].tolist()

    def compDF(self, a, b: pd.DataFrame, eps: float):
        "compare two purely numerical DataFrames"
        if np.linalg.norm(a.subtract(b)) > eps:
            self.fail(f"{a} != {b}")

    def test_ZPT(self):
        self.assertEqual(self.got.ZPT, self.want.ZPT)

    def test_freqs(self):
        self.compDF(self.got.freqs, self.want.freqs, 1e-7)

    def test_rots(self):
        self.compDF(self.got.rots, self.want.rots, 1e-7)

    def test_deltas(self):
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
            self.want.Rhead,
        )
        self.assertEqual(
            self.got.Ralpha,
            self.want.Ralpha,
        )
        self.assertEqual(
            self.got.Requil,
            self.want.Requil,
        )
        self.assertEqual(
            self.got.fermi,
            self.want.fermi,
        )
        self.assertEqual(
            self.got.Be,
            self.want.Be,
        )
        self.assertEqual(self.got.Lin, self.want.Lin)
        self.assertEqual(self.got.Imag, self.want.Imag)
        self.assertEqual(
            self.got.LX,
            self.want.LX,
        )


class TestDegenerate(unittest.TestCase):
    got = summ.Spectro(
        "drane.json",
        [
            [0],
            [1, 2],
            [4],
            [5, 6],
            [8, 9],
            [11, 12],
            [13, 14],
            [16, 17],
        ],
    )
    want = summ.Spectro
