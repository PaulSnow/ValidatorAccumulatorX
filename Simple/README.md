# LoadTest

LoadTest represents a simulation of a system of accumulators and test 
transaction generators to explore how much data can be collected and 
anchored into a blockchain like Factom.  The goal is not to just create 
the merkle roots for proof of existance, but to do so using a system of 
chains similar to what is used in Factom to order and segregate the data 
into chains.

The main purpose of the LoadTest is to get estimate of what throughput is 
possible with a Accumulator/Validator architecture running on top of Factom.  
To collect some numbers, we used Go to:

1. Create an accumulator that accepts hashes and builds iteratively a Merkle Tree
2. Create a test generator that creates and validates simple token transactions
3. Create a router that starts accumulators per chain to collect data 
on a per chain ID basis

##Results

Actual performance depends on many variables, but the main conclusion is clear.  
Such an architecture can collect and secure data on Factom at unbounded rates.  
With the data organized into chains, using a typical single I7 processor, this 
architecture can handle over 500 tps with as many as 100k chains.

##TBD

The following need to be done:

1. Code push both data and proofs into a distributed key value store  
2. Produce receipts of data from the Accumulator/Validator database
3. Provide access to data in chains to mirror the APIs provided by Factom
