package main

import (
	"fmt"
	"os"
)

func initConst() {
	if *plain && *tex {
		fmt.Fprintf(os.Stderr,
			"conflicting format options %q and %q specified, aborting\n",
			"plain", "tex")
		os.Exit(1)
	} else if *plain {
		DeltaOrder = []string{
			"Delta_J ",
			"Delta_K ",
			"Delta_JK",
			"delta_J ",
			"delta_K ",
		}
		PhiOrder = []string{
			"Phi_J ",
			"Phi_K ",
			"Phi_JK",
			"Phi_KJ",
			"phi_j ",
			"phi_jk",
			"phi_k ",
		}
	} else if *tex {
		DeltaOrder = []string{
			"$\\Delta_{J }$",
			"$\\Delta_{K }$",
			"$\\Delta_{JK}$",
			"$\\delta_{J }$",
			"$\\delta_{K }$",
		}
		PhiOrder = []string{
			"$\\Phi{_J }$",
			"$\\Phi{_K }$",
			"$\\Phi{_JK}$",
			"$\\Phi{_KJ}$",
			"$\\phi{_j }$",
			"$\\phi{_jk}$",
			"$\\phi{_k }$",
		}
	} else {
		var (
			upperDelta = "\u0394"
			lowerDelta = "\u03B4"
			upperPhi   = "\u03A6"
			lowerPhi   = "\u03C6"
		)
		DeltaOrder = []string{
			upperDelta + "_J ",
			upperDelta + "_K ",
			upperDelta + "_JK",
			lowerDelta + "_J ",
			lowerDelta + "_K ",
		}
		PhiOrder = []string{
			upperPhi + "_J ",
			upperPhi + "_K ",
			upperPhi + "_JK",
			upperPhi + "_KJ",
			lowerPhi + "_j ",
			lowerPhi + "_jk",
			lowerPhi + "_k ",
		}
	}
}
