# chemutils

Utility libraries and executables for chemistry problems. They are
primarily tailored to my group's research on rovibrational
spectroscopy via quartic force fields.

## atom

atom is an exploratory symmetry library.

## cccrel

cccrel is an executable that reads CcCR output files and prints a helpful summary, along
with the final relative energies needed to pass into ANPASS.

### TODO
* Implement -mono flag for reading component energies from a single
  file

## summarize

summarize is a library and an executable that reads output from
FORTRAN programs and summarizes/formats them nicely. Currently only
SPECTRO is supported.

### TODO
* Implement -tex flag
  * Still need document header and footer
    * Flag to disable?
  * Still need delimiters
* Implement -org flag
* Support ANPASS
* Support INTDER
* Make tables from multiple input files
