*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        39        39
N        0.0000000000       -0.1546592722       -1.1519403594
H        0.0000000000        1.6634878931       -1.9324718219
O        0.0000000000        0.0329391691        1.1338225903
}
 
basis=avtz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
