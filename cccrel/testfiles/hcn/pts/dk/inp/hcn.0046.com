*** hcn CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        46        46
N        0.0000000000        0.0000000000        1.1310004776
C        0.0000000000        0.0000000000       -1.0667464493
H        0.0000000000        0.0000000000       -3.0613240104
}
 
basis=vtz-dk
dkroll=0
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
