# #!/bin/bash

if [ -z "$1" ]; then
    echo "Error: Ptau file must be provided as the first argument"
    exit 1
fi

ptau_file="$1"
if [ ! -f "$ptau_file" ]; then
    echo "Error: The file '$ptau_file' does not exist"
    exit 1
fi

#

if [ -z "$2" ]; then
    echo "Error: Last contribution must be provided as the second argument"
    exit 1
fi

latest_contribution="$2"
if [ ! -f "$latest_contribution" ]; then
    echo "Error: The file '$latest_contribution' does not exist"
    exit 1
fi

#

if [ -z "$3" ]; then
    echo "Error: phase2 evaluations file must be provided as the third argument"
    exit 1
fi

phase2Evaluations="$3"
if [ ! -f "$phase2Evaluations" ]; then
    echo "Error: The file '$phase2Evaluations' does not exist"
    exit 1
fi

#

if [ -z "$4" ]; then
    echo "Error: r1cs file must be provided as the fourth argument"
    exit 1
fi

r1cs="$4"
if [ ! -f "$r1cs" ]; then
    echo "Error: The file '$r1cs' does not exist"
    exit 1
fi

if [ ! -f "setup" ]; then
    echo "Pulling and building semaphore-phase2-setup"
    git clone https://github.com/irfanbozkurt/semaphore-phase2-setup.git
    cd semaphore-phase2-setup
    go build
    cd ..
    mv semaphore-phase2-setup/semaphore-phase2-setup ./setup
    rm -rf semaphore-phase2-setup
    echo ""
fi

echo "Extracting pkey and vkey"
./setup extract-keys $ptau_file $latest_contribution $phase2Evaluations $r1cs
echo ""

rm setup
