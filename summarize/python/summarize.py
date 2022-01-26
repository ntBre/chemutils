import json
import pandas as pd
import numpy as np


class Spectro:
    def __init__(self, filename, deg_modes=None):
        """load a Spectro from a JSON file. deg_modes is a list of
        lists where each element in the list is a list of modes in the
        original representation corresponding to the new modes in that
        position.

        For example,

        deg_modes = [
            [0],
            [1, 2],
            [4],
            [5, 6],
            [8, 9],
            [11, 12],
            [13, 14],
            [16, 17],
        ]

        corresponds to the first mode mapping to itself, the second
        and third modes combining to give the new second mode, and so
        on. Mode indices 3, 7, and 10 are omitted from the new
        representation. The resulting modes are averages of the
        composing modes.

        """
        with open(filename) as f:
            js = json.load(f)
        self.ZPT = js["ZPT"]
        self.freqs = pd.DataFrame(
            {
                "harms": pd.Series(js["Harm"]),
                "funds": pd.Series(js["Fund"]),
                "corrs": pd.Series(js["Corr"]),
            }
        )
        rots = [sorted(x, reverse=True) for x in js["Rots"]]
        self.rots = pd.DataFrame(rots, columns=["A", "B", "C"])
        self.Deltas = js["Deltas"]
        self.Phis = js["Phis"]
        self.Rhead = js["Rhead"]
        self.Ralpha = js["Ralpha"]
        self.Requil = js["Requil"]
        self.Fermi = js["Fermi"]
        self.Be = js["Be"]
        self.Lin = js["Lin"]
        self.Imag = js["Imag"]
        self.LX = js["LX"]
        if deg_modes is not None:
            self.degenerate(deg_modes)

    def __repr__(self):
        return f"""{{
  "ZPT": {self.ZPT},
  "freqs": {self.freqs},
  "Rots": {self.Rots},
  "Deltas": {self.Deltas},
  "Phis": {self.Phis},
  "Rhead": {self.Rhead},
  "Ralpha": {self.Ralpha},
  "Requil": {self.Requil},
  "Fermi": {self.Fermi},
  "Be": {self.Be},
  "Lin": {self.Lin},
  "Imag": {self.Imag},
  "LX": {self.LX},
}}"""

    def degenerate(self, deg_modes):
        """average all applicable spectroscopic constants across
        degenerate modes. See the __init__ documentation for more
        details about the form of deg_modes"""
        new_freqs = pd.DataFrame()
        for _type in ["harms", "funds", "corrs"]:
            tmp = []
            for modes in deg_modes:
                tmp.append(
                    np.average(
                        self.freqs[_type][modes],
                    )
                )
            new_freqs.insert(len(new_freqs.columns), _type, tmp)
        self.freqs = new_freqs

        new_rots = pd.DataFrame()
        for _type in ["A", "B", "C"]:
            tmp = []
            rot_modes = deg_modes.copy()
            # always keep the 0th element since A_0 is a real thing,
            # unlike v_0
            rot_modes.insert(0, [-1])
            for modes in rot_modes:
                tmp.append(
                    np.average(
                        self.rots[_type][[x+1 for x in modes]],
                    )
                )
            new_rots.insert(len(new_rots.columns), _type, tmp)
        self.rots = new_rots

    def freq_table(self):
        """output the harmonic and resonance-corrected frequencies as
        a LaTeX table"""
        print("\\begin{tabular}{ll}")
        print("%11s & %7s \\\\" % ("Mode", "Freq"))
        print("\\hline")
        # TODO different case when symmetries and/or descriptions are
        # available
        for i, v in enumerate(self.freqs["harms"]):
            print("\\omega_{%2d} & %7.1f \\\\" % (i + 1, v))
        for i, v in enumerate(self.freqs["funds"]):
            print("   \\nu_{%2d} & %7.1f \\\\" % (i + 1, v))
        print("\\end{tabular}")
