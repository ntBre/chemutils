import json
import warnings
import numpy as np
import re
import collections

# ignore warning on pandas to_latex()
warnings.simplefilter(action="ignore", category=FutureWarning)
import pandas as pd

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

        These degeneracies are applied to all three frequency types
        contained in the DataFrame freqs field, the DataFrame of rots
        in that field, and to the list of string Fermi resonance
        relationships in the fermi field.

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

    def harms(self):
        return self.freqs["harms"]

    def corrs(self):
        return self.freqs["corrs"]

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


def freq_table(spec: Spectro, symms=None, descs=None, ints=None):
    """output the harmonic and resonance-corrected frequencies of spec
    as a LaTeX table. symms are symmetry labels for the modes. It
    should be the same length as the number of vibrational
    modes. Similarly, descs is a textual description of the modes and
    should have the same length as well.

    The harmonic frequencies are numbered from 1 and prefaced with
    \\omega, while the fundamental frequencies are prefaced with \\nu

    The default output is a table like

    Mode & Freq \\
    \\hline

    but a full output looks like

    Mode & Symm. & Desc. & Freq. \\
    \\hline

    ints is a list of intensities. If the length is equal to one set
    of frequencies, they are assumed to be harmonic intensities and
    are placed in parens after the corresponding harmonic
    frequency. If the length is twice the frequencies, the first half
    is assumed to be harmonic and the second half anharmonic.

    """
    tab = pd.DataFrame()
    labels = [f"$\\omega_{{{i+1}}}$" for i in range(len(spec.freqs["harms"]))]
    labels.extend([f"$\\nu_{{{i+1}}}$" for i in range(len(spec.freqs["corrs"]))])
    mode, freq, symm, desc, _int = "Mode", "Freq.", "Symm.", "Desc.", "Int."
    tab[mode] = labels
    tab[freq] = list(spec.harms()) + list(spec.corrs())
    col_names = [mode, freq]
    if descs is not None:
        tab[desc] = 2 * descs
        col_names.insert(1, desc)
    if symms is not None:
        tab[symm] = 2 * symms
        col_names.insert(1, symm)
    if ints is not None:
        tmp = ["%.0f" % x for x in ints]
        if len(tmp) == len(tab[freq])//2:
            tmp.extend([""] * len(tmp))
        tab[_int] = tmp
        col_names.append(_int)
    tab = tab.reindex(columns=col_names)
    rx = re.compile(r"^(\\).*rule$", re.MULTILINE)
    print(
        re.sub(
            rx, r"\\hline", tab.to_latex(escape=False, float_format="%.1f", index=False)
        )
    )
