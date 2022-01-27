import json
import pandas as pd
import numpy as np
import re
import collections

DEBUG = False


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
        self.fermi = js["Fermi"]
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
  "rots": {self.rots},
  "Deltas": {self.Deltas},
  "Phis": {self.Phis},
  "Rhead": {self.Rhead},
  "Ralpha": {self.Ralpha},
  "Requil": {self.Requil},
  "fermi": {self.fermi},
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
                        self.rots[_type][[x + 1 for x in modes]],
                    )
                )
            new_rots.insert(len(new_rots.columns), _type, tmp)
        self.rots = new_rots

        # build fermi resonance patterns
        res = []
        for i, d in enumerate(deg_modes):
            for v in d:
                res.append([f"(.*)v_{v+1}$", f"v_{i+1}"])

        # loop over fermi resonance strings, first pass to update
        # values
        new_fermi = []
        for fs in self.fermi:
            # break apart equivalences
            fs = fs.split("=")
            row = []
            for f in fs:
                # break apart type 2 resonances to prevent double
                # substitution
                tmp = []
                sp = f.split("+")
                for s in sp:
                    matched = False
                    for rx in res:
                        match = re.match(rx[0], s)
                        if match:
                            if DEBUG:
                                print(
                                    f"match {rx[0]} in {s}, "
                                    + "replacing with {match.group(1) + rx[1]}"
                                )
                            tmp.append(re.sub(rx[0], match.group(1) + rx[1], s))
                            matched = True
                            break
                    # if part of it wasn't matched, the pair is nonsense, so
                    # remove it if it was pushed partially and break
                    if not matched and len(tmp) > 0:
                        tmp.pop()
                        break
                j = "+".join(tmp)
                if j != "":
                    row.append(j)
            new_fermi.append("=".join(row))
        # second pass to deduplicate within groups
        new_new = collections.OrderedDict()
        for fermi in new_fermi:
            sp = fermi.split("=")
            od = collections.OrderedDict()
            for s in sp:
                od[s] = True
            if len(od) > 1:
                new_new["=".join(od.keys())] = True
        self.fermi = list(new_new.keys())

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
