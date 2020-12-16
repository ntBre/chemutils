# chemutils

Utility libraries and executables for chemistry problems. They are
primarily tailored to my group's research on rovibrational
spectroscopy via quartic force fields.

## atom

atom is an exploratory symmetry library.

## cccrel

cccrel is an executable that reads CcCR output files and prints a helpful summary, along
with the final relative energies needed to pass into ANPASS.

### Usage

If you have your CcCR points files organized in a directory called
`pts`, just run `$ cccrel` in the directory above `pts`. Otherwise,
use the `-help` flag to see how to direct it to the right input files.

### TODO
* Implement -mono flag for reading component energies from a single
  file

## summarize

summarize is a library for parsing SPECTRO output and a standalone
program that formats it nicely.

### Usage

summarize takes the SPECTRO output file as its first argument. For the
example in testfiles:

```
$ summarize spectro.out
```

will write the text output to stdout. The `-tex` flag will convert the
format to LaTeX, and you can use the `-nohead` flag to disable
printing of LaTeX header/footer information in case you just want to
include the output in an existing file. The `-plain` flag can be used
to disable the Unicode characters in the Delta and Phi output.

### TODO
* Implement -org flag
* Make tables from multiple input files

## spectro

spectro is a library for actually running SPECTRO. Eventually it will
also be a standalone executable wrapping the Fortran version of
SPECTRO.


