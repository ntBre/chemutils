package main

import (
	"fmt"
	"os"

	"github.com/ntBre/chemutils/summarize"
)

func initConst() {
	if *plain && *tex {
		fmt.Fprintf(os.Stderr,
			"conflicting format options %q and %q specified, aborting\n",
			"plain", "tex")
		os.Exit(1)
	} else if *plain {
		summarize.DeltaOrder = []string{
			"Delta_J ",
			"Delta_K ",
			"Delta_JK",
			"delta_J ",
			"delta_K ",
		}
		summarize.PhiOrder = []string{
			"Phi_J ",
			"Phi_K ",
			"Phi_JK",
			"Phi_KJ",
			"phi_j ",
			"phi_jk",
			"phi_k ",
		}
	} else if *tex {
		summarize.DeltaOrder = []string{
			"$\\Delta_{J }$",
			"$\\Delta_{K }$",
			"$\\Delta_{JK}$",
			"$\\delta_{J }$",
			"$\\delta_{K }$",
		}
		summarize.PhiOrder = []string{
			"$\\Phi{_J }$",
			"$\\Phi{_K }$",
			"$\\Phi{_JK}$",
			"$\\Phi{_KJ}$",
			"$\\phi{_j }$",
			"$\\phi{_jk}$",
			"$\\phi{_k }$",
		}
	}
}
