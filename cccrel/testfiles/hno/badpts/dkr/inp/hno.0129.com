*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
       129       129
N        0.0000000000       -0.1824487036       -1.1404926860
H        0.0000000000        1.6833874597       -1.9296233002
O        0.0000000000        0.0408290340        1.1195263952
}
 
basis=vtz-dk
dkroll=1
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
