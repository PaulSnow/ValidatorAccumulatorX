
/Elapsed/ {
   gsub(",","",$0)
   gsub("s","",$1)
   split($1,t,"m",seps)
   tps = $11
}

/Entry limit/ {
   gsub(",","",$0)
   entrylimit = $4
}

/Chain limit/ {
   gsub(",","",$0)
   chainlimit = $4
}

/TPS limit/ {
   gsub(",","",$0)
   tpslimit = $4
}

/Total Entries/ {
   gsub(",","",$0)
   tps = $9
}

/of Accumulators/ {
   gsub(",","",$0)
   accs = $9
}
 

/Recorded/ {
   gsub(",","",$0)
  print "Entries,", entrylimit, " ,chains, ",chainlimit, ", tps Limit, ", tpslimit, ",tps, ", tps, ",accs,", accs
}
