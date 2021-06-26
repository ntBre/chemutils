#!/usr/bin/python

#Readinter - Py Version
#OCO2 CCSD(T)-F12/cc-pVTZ
# Takes the compound name and then the order of the atoms as arguments
# For example, for Al2O3, the command line looks like
# $ python brent.py al2o3 Al Al O O O 
import os
import argparse
parser = argparse.ArgumentParser()
parser.add_argument("compname", type=str)
parser.add_argument("atoms", nargs="*", type=str)
args = parser.parse_args()
print args.compname
print args.atoms

os.mkdir("inp")
    
compname = args.compname
atoms = args.atoms
chunk = len(atoms)
infile = open("./file07", 'r')
lines = infile.readlines()

count = 0
line = 0
disp = 0
while line < len(lines):
    if '#' in lines[line]:
        disp += 1
        filenum = compname + "." + str(disp).zfill(4)
        filename = "inp/" + filenum +  ".com"
        outfile = open(filename, 'w')
        outfile.writelines(["*** ", str(compname) , " CCSD(T)/TZ-F12 QFF\n",
                            "memory, 995, m;\n",
        "gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;\n",
                            "gthresh,optgrad=1.d-8,optstep=1.d-8;\n",
                            "nocompress;\n",
                            "geomtyp=xyz\n",
                            "bohr\n",
                            "geometry={\n",
                            str(chunk), "\n",
                            "%10d%10d\n" % (disp,disp)])
        for i in range(chunk+1):
            if i < len(atoms):
                outfile.writelines([atoms[i], lines[line+1]])
            line += 1
        outfile.writelines(["}\n",
                            " \n",
                            "basis=vtz-dk\n",
			    "dkroll=1\n",
                            "  {hf,maxit=500;accu,20;}\n",
                    "{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}\n"])
        pbsname = "inp/" + filenum + ".pbs"
        pbsfile = open(pbsname, 'w')
        pbsfile.writelines(["#PBS -S /bin/sh\n",
                            "#PBS -j oe\n",
                            "#PBS -W umask=022\n",
                            "#PBS -l walltime=5000:00:00\n",
                            "#PBS -l ncpus=1\n",
                            "#PBS -l mem=32gb\n",
                            "\n",
                            "module load intel\n",
                            "module load mvapich2\n",
                            "module load pbspro\n",
			    "export PATH=/usr/local/apps/molpro/2015.1.35/bin:$PATH\n",
                            "\n",
                            "export WORKDIR=$PBS_O_WORKDIR\n",
                            "export TMPDIR=/tmp/$USER/$PBS_JOBID\n",
                            "cd $WORKDIR\n",
                            "mkdir -p $TMPDIR\n",
                            "\n",
                            "date\n",
                            "molpro -t 1 " + filenum + ".com" + "\n",
                            "date\n",
                            "\n",
                            "rm -rf $TMPDIR\n"])
        submit = open("inp/submit", 'a')
        submit.write("qsub " + filenum + ".pbs\n")

os.chmod("inp/submit", 0755)
