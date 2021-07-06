*** hcn CCSD(T)/TZ-F12 QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geomtyp=xyz
bohr
geometry={
3
        52        52
N        0.0000000000        0.0000000000        1.1404491076
C        0.0000000000        0.0000000000       -1.0667464493
H        0.0000000000        0.0000000000       -3.0707726403
}
 
basis={
default=aug-cc-pvtz
s,C,8236.0,1235.0,280.8, 79.27,25.59, 8.997,3.319
s,C,0.9059,0.3643,0.1285000
p,C,56.0,18.71,4.133,0.2827,0.3827,0.1209
d,C,30.0,10.0,3.3,1.097,0.318
f,C,7.0,2.3,0.7610
s,N,11420.0,1712.0,389.3,110.0,35.57,12.54,4.644
s,N,1.293,0.5118,0.1787
p,N,79.89,26.63,5.948,1.742,0.555,0.1725
d,N,45.0,15.0,5.0,1.654,0.469
f,N,9.9,3.3,1.093
s,O,15330.0,2299.0,522.4,147.3,47.55,16.76,6.207
s,O,1.752,0.6882,0.2384
p,O,103.5,34.46,7.749,2.28,0.7156,0.214
d,O,63.0,21.0,7.0,2.314,0.645
f,O,12.9,4.3,1.428
s,mg,164900.0,24710.0,5628.0,1596.0,521.0;
s,mg,188.0,73.01,29.90,12.54,4.306,1.826;
s,mg,0.7417,0.0761,0.145,0.033,0.0129;
p,mg,950.70,316.90,74.86,23.72,8.669,3.363;
p,mg,1.310,0.4911,0.2364,0.08733,0.03237,0.00745;
d,mg,1.601,0.686,0.126,0.294,0.0468;
f,mg,1.372,0.588,0.094,0.252;
}
  {hf,maxit=500;accu,20;}
{ccsd(t),nocheck,maxit=250;orbital,IGNORE_ERROR;core}
