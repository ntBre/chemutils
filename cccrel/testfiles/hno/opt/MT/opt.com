*** AlONC, CCSD(T)/MT QFF
memory, 995, m;
gthresh,energy=1.d-12,zero=1.d-22,oneint=1.d-22,twoint=1.d-22;
gthresh,optgrad=1.d-8,optstep=1.d-8;
nocompress;
geometry={

N
O 1 no
H 1 nh 2 hno
}

 no=1.2 ANG
 nh=0.8 ANG
 hno=105.0 DEG

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
}
  {hf,maxit=500;wf,charge=0,spin=0;accu,20;}
  {ccsd(t),maxit=250;wf,charge=0,spin=0;orbital,IGNORE_ERROR;}
 {optg,grms=1.d-8,srms=1.d-8}

