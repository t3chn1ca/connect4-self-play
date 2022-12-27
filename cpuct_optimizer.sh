#!/bin/bash

for i in {1..20..2}
do
	#echo "Cpuct = 0.$i"
	cpuct=`echo $i/10|bc -l|cut -c1,2,3`
	echo "Cpuct = "$cpuct
	./nn-vs-mcts $cpuct > logs/nn-vs-mcts-gen4-cpuct-eq-0$cpuct.log

done	


