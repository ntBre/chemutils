*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        26        26
H        0.0000000000        1.6612745377       -1.9198808225
N        0.0000000000       -0.1603243845       -1.1474393392
O        0.0000000000        0.0408176368        1.1167305707
}
 
basis=av5z
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
