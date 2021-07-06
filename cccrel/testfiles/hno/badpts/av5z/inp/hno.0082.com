*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        82        82
N        0.0000000000       -0.1546325333       -1.1492399406
H        0.0000000000        1.6686154288       -1.9417895529
O        0.0000000000        0.0277848945        1.1404399024
}
 
basis=av5z
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
