This tool allows Gnark 0.9.2 users to run a 2nd-phase MPC ceremony for generating the proving and verifying keys for the Groth16 protocol as presented in [BGM17](https://eprint.iacr.org/2017/1050.pdf). It does not include the beacon contribution since it was proved in [KMSV21](https://eprint.iacr.org/2021/219.pdf) that the security of the generated SRS still holds without it.

## Reasoning behind a custom trusted setup

Each groth16 proof of a circuit requires a trusted setup that has 2 parts: a phase 1 which is also known as a "Powers of Tau ceremony" which is universal (the same one can be used for any circuit) and a phase2 which is circuit-specific, meaning that you need to a separate phase2 for every single circuit.

For phase 1, you can reuse the setup done by the joint effort of many community members - it is a powers of tau ceremony with 54 different contributions ([more info here](https://github.com/privacy-scaling-explorations/perpetualpowersoftau)). A list of downloadable `.ptau` files can be found [here](https://github.com/iden3/snarkjs/blob/master/README.md#7-prepare-phase-2).

## Notice

Phase 2 is circuit-specific, so if you have `n` circuits, then you need to run this phase `n` times.

Remember that you need sufficiently high powers of tau ceremony to generate a proof for a circuit with a given amount of constraints ($2^{POW_OF_TAU} >= CIRCUIT_CONSTRAINTS$)

## Pre-requisites

1. Git
2. Go
3. Minimum RAM requirement is 16GB

# Usage

1. Initialize phase2

After exporting your circuit as in gnark's r1cs format, place it in the same directory as [init_phase2.sh](./bash/init_phase2.sh). Then:

```bash
$ chmod 777 init_phase2.sh
$ ./init_phase2.sh <powerOfTau>
```

Where 'powerOfTau' is the smallest power of 2 covering your number of constraints.
This will download corresponding ptau file, initialize phase2 ceremony, and perform 2 contributions.

Output of this step will be a 'output_2.ph2' file, and a 'phase2Evaluations' file. Output will be used by the next contributor to perform their contribution, and 'phase2Evaluations' will be used in the end when extracting vkey and pkey.

2. Contributions

Contributors will use the output of the last contribution to perform their own:

```bash
$ chmod 777 contribute.sh
$ ./contribute.sh <outputOfPrevContribution> <outputFileName>
```

3. Contribution verification

The coordinator needs to verify every contribution (w.r.t. previous contribution)

```bash
$ chmod 777 verify.sh
$ ./verify.sh <outputOfContributionToBeVerified> <outputOfPrevContribution>
```

4. Extracting the keys

After repeating steps 2 and 3 as much as necessary, anyone can extract pkey and vkey containing randomness created by a group of people. These keys finalize the trusted setup and they're intended to be used in production.

```bash
$ chmod 777 extract_keys.sh
$ ./extract_keys.sh <ptau_file> <latest_contribution> <phase2Evaluations> <r1cs>
```

## Keys Extraction

At the end of the ceremony, the coordinator runs `semaphore-phase2-setup key <lastPhase2Contribution.ph2>` which will output **Groth16 bn254 curve** `pk` and `vk` files

## Acknowledgements

This repository is a fork of the [semaphore-mtb](https://github.com/worldcoin/semaphore-mtb) repository. We would like to thank the authors of the original repository for their work as this project is a slight tweak of the original work to fit our needs.

We appreciate the community efforts to generate [a good universal SRS](https://github.com/privacy-scaling-explorations/perpetualpowersoftau) for everyone's benefit to use and for the [iden3 team for building [snarkjs](https://github.com/iden3/snarkjs).
