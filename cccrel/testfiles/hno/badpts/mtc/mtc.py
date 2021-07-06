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
                            "basis={\n",
				"default=aug-cc-pvtz\n"
				"s,C,8236.0,1235.0,280.8, 79.27,25.59, 8.997,3.319\n",
				"s,C,0.9059,0.3643,0.1285000\n",
				"p,C,56.0,18.71,4.133,0.2827,0.3827,0.1209\n",
				"d,C,30.0,10.0,3.3,1.097,0.318\n",
				"f,C,7.0,2.3,0.7610\n",
				"s,N,11420.0,1712.0,389.3,110.0,35.57,12.54,4.644\n",
				"s,N,1.293,0.5118,0.1787\n",
				"p,N,79.89,26.63,5.948,1.742,0.555,0.1725\n",
				"d,N,45.0,15.0,5.0,1.654,0.469\n",
				"f,N,9.9,3.3,1.093\n",
				"s,O,15330.0,2299.0,522.4,147.3,47.55,16.76,6.207\n",
				"s,O,1.752,0.6882,0.2384\n",
				"p,O,103.5,34.46,7.749,2.28,0.7156,0.214\n",
				"d,O,63.0,21.0,7.0,2.314,0.645\n",
				"f,O,12.9,4.3,1.428\n",
				"s,mg,164900.0,24710.0,5628.0,1596.0,521.0;\n",
				"s,mg,188.0,73.01,29.90,12.54,4.306,1.826;\n",
				"s,mg,0.7417,0.0761,0.145,0.033,0.0129;\n",
				"p,mg,950.70,316.90,74.86,23.72,8.669,3.363;\n",
				"p,mg,1.310,0.4911,0.2364,0.08733,0.03237,0.00745;\n",
				"d,mg,1.601,0.686,0.126,0.294,0.0468;\n",
				"f,mg,1.372,0.588,0.094,0.252;\n}\n",
                            "  {hf,maxit=500;accu,20;}\n",
                    "{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;core}\n"])
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
