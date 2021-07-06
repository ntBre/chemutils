*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        47        47
H        0.0000000000        1.6642187737       -1.9191672325
N        0.0000000000       -0.1639827206       -1.1381119749
O        0.0000000000        0.0415317369        1.1066896164
}
 
basis=avtz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
