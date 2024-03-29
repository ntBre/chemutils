*** hno CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        68        68
H        0.0000000000        1.6674498676       -1.9398038131
N        0.0000000000       -0.1539508152       -1.1430180377
O        0.0000000000        0.0282687377        1.1322322598
}
 
basis=avtz
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;}
