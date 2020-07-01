#!/usr/bin/env bash

echo > time.out

for j in `seq 1 1000000`;
do
		echo "Timer", $j
		pid="$(ps ax | grep "[0-9] ValAcc" | gawk "{print \$1}")"
		[[ !  -z  $pid  ]] && echo "pid=" $pid
		[[ !  -z  $pid  ]] && python log.py $pid >> time.out
		[[ -z  $pid  ]] && echo "waiting"  
		sleep 10
done

