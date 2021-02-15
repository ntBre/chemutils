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

## diagram

diagram is a tool for adding labels to images using imagemagick

### Usage

Obviously it requires imagemagick to be installed with pango support.

```
$ diagram [-grid h,v] captions image
```

Will put the captions in the `captions` file onto the PNG image
`image`. The format of the caption file is

```
Text Size  xpos,ypos
```

Where `Text` is a single "word" (no spaces allowed), `Size` is an
integer points for the font size, and `xpos,ypos` is a comma-separated
pair denoting the x- and y-position of the label. The `-grid` flag
takes a comma-separated pair of integers denoting the number of
`h`orizontal and `v`ertical gridlines to draw.

## polecat

polecat is a tool for drawing simple images of molecules with their
associated dipole and rotational axes.

### Usage

```
$ polecat dipole.out
```

where `dipole.out` is a Molpro output file from a dipole calculation.

## spectro

spectro is a library for running SPECTRO and a standalone executable
wrapping the Fortran version of SPECTRO.

### Usage

spectro takes a SPECTRO intput file as an argument, runs the Fortran
version of SPECTRO on it, parses the output to identify Fermi,
Coriolis, and Darling-Dennison resonances, then writes spectro2.in and
reruns SPECTRO. The command line looks like:

```
$ spectro spectro.in
```

The `-cmd` flag can be used to specify an alternative SPECTRO
executable.

## summarize

summarize is a library for parsing SPECTRO output and a standalone
program that formats it nicely.

### Usage

summarize takes a SPECTRO output file as its only argument. For example,

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
