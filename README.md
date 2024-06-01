This tool is a wrapper around Gnark to perform a 2nd-phase MPC ceremony for generating the proving and verifying keys for the Groth16 protocol as presented in [BGM17](https://eprint.iacr.org/2017/1050.pdf). It does not include the beacon contribution since it was proved in [KMSV21](https://eprint.iacr.org/2021/219.pdf) that the security of the generated SRS still holds without it.

## Pre-requisites

1. Git
2. Go

## Usage

```bash
# Initialize a phase2
go run main.go p2n path_to.ptau path_to.r1cs output_0.ph2
### Outputs output_0.ph2 and a phase2Evaluations file
### output_0.ph2 will be used for the first contribution
### phase2Evaluations file will be used in the last step

# Contribute to phase2
go run main.go p2c path_to_input.ph2 output_<contribution_number>.ph2
### Outputs output_<contribution_number>.ph2 file

# Verify a phase2 contribution
go run main.go p2v path_to_be_verified.ph2 path_to_previous_conribution.ph2
### Outputs true or false to the console

# Extract pkey and vkey from the final ph2 file
go run main.go extract-keys path_to.ptau path_to_latest_conribution.ph2
### Outputs proving and verification keys
```
