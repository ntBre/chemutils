*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        72        72
N        0.0000000000       -0.1570468069       -1.1444214411
H        0.0000000000        1.6670758796       -1.9349557185
O        0.0000000000        0.0317387173        1.1287875686
}
 
basis=av5z
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
