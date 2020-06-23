$1+0 > 0 && $2+0 > 0 {
   if (maxCPU < $1) {
      maxCPU = $1
   }
   if (minCPU == 0 || minCPU > $1) {
      minCPU = $1
   }
   sumCPU += $1
   cntCPU++

   if (maxMEM < $2) {
      maxMEM = $2
   }
   if (minMEM == 0 || minMEM > $2) {
      minMEM = $2
   }
   sumMEM += $2
   cntMEM++
}

/no such process/ {

   if (sumCPU > 0) {
      printf("CPU, %6.2f,%6.2f,%6.2f, MEM, %6.2f,%6.2f,%6.2f\n",
        maxCPU/800, minCPU/800, sumCPU/cntCPU/800,
        maxMEM, minMEM, sumMEM/cntMEM)
      maxCPU = minCPU = sumCPU = cntCPU = 0
      maxMEM = minMEM = sumMEM = cntMEM = 0
   }
}
