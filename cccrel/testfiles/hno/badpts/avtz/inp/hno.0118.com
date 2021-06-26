*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
       118       118
N        0.0000000000       -0.1747064527       -1.1498217351
H        0.0000000000        1.6756440532       -1.9270384139
O        0.0000000000        0.0408301894        1.1262705579
}
 
basis=avtz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
