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
