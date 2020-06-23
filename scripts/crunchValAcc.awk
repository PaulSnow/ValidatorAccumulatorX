
/Elapsed/ {
   gsub(",","",$0)
   gsub("s","",$1)
   split($1,t,"m",seps)
   tps = $11
}

/Running/ {
   gsub(",","",$0)
   chains = $12
   entries = $8
   tps =16
}

/Recorded/ {
   gsub(",","",$0)
  print "Entries,", entries, " ,chains, ",chains, ", tps Limit, ", tps, ",time, ","00:" t[1] ":" t[2], " ,tps, " tps
}
