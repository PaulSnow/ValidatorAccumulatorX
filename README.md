# Simple

This code base was used to measure just the overhead of creating merkle trees by adding elements one at
a time.  The goal was originally to spread out the overhead of creating merkle trees over time, so there
would not be a huge cpu demand at the end of a block.

Turns out the algorithm is very simple and flexible.  The current code gives us a rough idea of the
performance of building the blockchain infrastructure under a stream of ordered, validated hashes from
some set of sources.

# ValidatorAccumulator

The Validator/Accumulator (ValAcc) architecture allows for multiple applications to use the same facility
to collect and order data.  These applications can happily run along side each other, or interact
with each other.

Cryptographic proofs for sets of data can be exported outside the ValAcc, and used by outside processes
without access to all the data held within the ValAcc itself.

Use cases include digital identities, IoT security, Supply Chain, Digital Rights Management, Document Management,
Loan Origination, Tokenization, Smart Contracts, and pretty much anyting else that can be done on a blockchain.

The ValAcc is intended to be run on top of the Factom Protocol, and used to extend the use cases that can be
built on the Factom Protocol.

The ValAcc can be run from the commandline to do various tests on performance on various architectures.
```
Usage of ValAcc:
  -a int
    	the number of accumulator instances used in this test (default 1)
  -c int
    	The number of chains updated while processing this test (default 1000)
  -e int
    	the number of entries to be processed in this test (default 1000000)
  -t int
    	the tps limit of data generated to run this test. if t < 0, no limit (default -1)
```

So running at the commandline:
```
     ValAcc -a 8 -e 30000000 -c 10000 -t -1
```     
yields on my machine: 
```
=========================
 -e <number of entries>
 -c <number of chains>
 -t <tps limit ( -1 is none)>
 -a <number of accumulators>
=========================
Entry limit of          30,000,000
Chain limit of             100,000
TPS limit of                  none
# of Accumulators                8
=========================

badger 2020/07/02 12:48:07 INFO: All 0 tables opened in 0s
badger 2020/07/02 12:48:07 INFO: All 0 tables opened in 0s
badger 2020/07/02 12:48:07 INFO: All 0 tables opened in 0s
badger 2020/07/02 12:48:07 INFO: All 0 tables opened in 0s
badger 2020/07/02 12:48:07 INFO: All 0 tables opened in 0s
badger 2020/07/02 12:48:07 INFO: All 0 tables opened in 0s
badger 2020/07/02 12:48:07 INFO: All 0 tables opened in 0s
badger 2020/07/02 12:48:07 INFO: All 0 tables opened in 0s
EOB 1
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Merkle DAG Root hash for 0 is 9eb418888d9fe9604ad57190c82e9f4c4bd14efa664ba19df3736d0a25bbc7b3
Merkle DAG Root hash for 1 is d4a7fbce7b9b28a7d86c2284361b65bea38ea0eccee96787803f5041d57758bc
Merkle DAG Root hash for 2 is 2f553533270b8c41be2ebed389073cd0c82c95e76b37ed96256a50368b10d451
Merkle DAG Root hash for 3 is 44e866a2ec78dfb9fc0bc9c87b71f0f7635e8d73212507d80409187dcf6cf3be
Merkle DAG Root hash for 4 is 3c4dda8c4093242111e32dfa2e4447ea172024c2e6bd6394c426c84d34039a9b
Merkle DAG Root hash for 5 is c31268bd85f6ec97afb05170084b2c800e2dcc330f08e49666331455db481fd2
Merkle DAG Root hash for 6 is 7b857cf271f4f000c55d915096c4c9137e6d61e98e0964ef5c9e78be47806b73
Merkle DAG Root hash for 7 is e2a700e490fa4c839d4dfcca91a036c2c087104c6f59104dbe3f415c2afeff34
Total Entries Written 6,398,061 to 100,000 chains, @ 426,537 tps
EOB 2
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Merkle DAG Root hash for 0 is c9c799ae4f9b0f311317a582b370d9da8c75ba81a82a2e6c748261f272fbc07a
Merkle DAG Root hash for 1 is a1e10bbc5c92351d0fea29fc4b6fe5b781cdca4a3d06b625b3ce2b7935e3a70f
Merkle DAG Root hash for 2 is 1de70daded87a39cb037b8ef77baa9ece2966730c96a9c154d5452cfc6d11548
Merkle DAG Root hash for 3 is d710fef08bfd2929b3b98f32ec272e0146df87aa880b32f667a1637f2fd643f8
Merkle DAG Root hash for 4 is 87edb6c8840c2cdb1c126d5664fca4a02b752c00fe7acedc22d0141d9d9638f1
Merkle DAG Root hash for 5 is abfa95dd62dab309ff06f179663ffb7f835c19b818ddcbb92410fdb8ecb53b65
Merkle DAG Root hash for 6 is b1085d37b070b76f6cb688e124556381e011e9ad4e85c3bd97427503505bf811
Merkle DAG Root hash for 7 is 75ff6dc9b2d0d74aa10a920e91a78991a74c5b156884b55f2abce196340d0e4e
Total Entries Written 13,065,742 to 100,000 chains, @ 483,916 tps
EOB 3
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Merkle DAG Root hash for 0 is 97a7efc7518c2dcb713d07abc9f05d45209f75abee2d29eeafc16b1cd042e1ae
Merkle DAG Root hash for 1 is e51bd442dc162006d0bef374518b3f6f374cdad97f14b6042f792e3ee7d6745b
Merkle DAG Root hash for 2 is d5adc0eb150f6a860d26970570e7000603d588b9a824cc60277bb49f25ce9fe1
Merkle DAG Root hash for 3 is b9549d126fc4f03b370fe2b92dac770e91af07ed358c68bea5678a813e5b34dc
Merkle DAG Root hash for 4 is d6ee702a40bb863439195645febc51ab8b0e918327e6fe0a3bf20f41ecdc9b4e
Merkle DAG Root hash for 5 is 12e92af6e5cdbaa6d28111d88049a0a43df9ddde1402a3adbba290b8fe72ae60
Merkle DAG Root hash for 6 is 7771395520df446d71c24c2e8a30dd1f687d528f74654374c7c876f0c2a650d3
Merkle DAG Root hash for 7 is 9d4c68b2fc09ae693462850f855c79b410b46228af3f22cbfabefbc66d6b6635
Total Entries Written 19,764,407 to 100,000 chains, @ 494,110 tps
EOB 4
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Processing EOB
Merkle DAG Root hash for 0 is ac4b203a66820fe74cc645d0364a60ea36dd148867796519f7255712a11d8c4b
Merkle DAG Root hash for 1 is c25b69d7c9dbb0295fc4fbb8adb2eb8567a8e8fb28fcd1b44923652bf2bd6ba0
Merkle DAG Root hash for 2 is fd40125952c6c4f74e9623a1f58de6d42df92d639b3ab3befd1cee9754019514
Merkle DAG Root hash for 3 is 7d122c9e2789248093696fa0c5cfd92fd1f3671dd9425ab6ab148c01f612cb1b
Merkle DAG Root hash for 4 is 596a30b40957decb18cc7430bb059dc29fb71af26c43b5dd7484b219341d6be8
Merkle DAG Root hash for 5 is cf083d74bc13cf9c659d2c1e8eda3e64b70990b9a64277b5a77ceb1411562d6c
Merkle DAG Root hash for 6 is 42d165fc2dd59501f3a93615c9c352a6ea58e7238f9d5f96fdfd0edffd5736bb
Merkle DAG Root hash for 7 is 14fd29783f2bc8ef28ffa1506d1135fbb94a886eaaa9d9edc1e6b5a912a99b5e
Total Entries Written 26,823,216 to 100,000 chains, @ 515,831 tps

====================
Recorded 30,000,000 Entries in 0 Blocks
Test complete.
```
