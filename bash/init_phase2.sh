# #!/bin/bash

find_r1cs_file() {
    local r1cs_file=$(find . -maxdepth 1 -type f -name "*.r1cs" | head -n 1)
    echo "$r1cs_file"
}

if ! [[ "$1" =~ ^[0-9]+$ ]] || (( $1 < 10 )) || (($1 > 28)); then
    echo "Error: Parameter \$1 must be a number >= 10 and <= 28"
    exit 1
fi

r1cs_file=$(find_r1cs_file)
if [ -z "$r1cs_file" ]; then
    echo "Error: No .r1cs file found in the current directory"
    exit 1
fi

# Check if the phase1 ceremony file exists
if [ ! -f "$1.ptau" ]; then
    echo "Downloading phase1 ceremony for 2^$1"
    wget https://hermez.s3-eu-west-1.amazonaws.com/powersOfTau28_hez_final_$1.ptau 
    mv powersOfTau28_hez_final_$1.ptau $1.ptau
    echo ""
fi

# Build phase2 executable
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

echo "Initializing phase2 ceremony"
./setup p2n $1.ptau $r1cs_file output_0.ph2
echo ""

echo "Performing the first contribution"
./setup p2c output_0.ph2 output_1.ph2
echo ""

echo "Verifying the first contribution"
./setup p2v output_1.ph2 output_0.ph2
echo ""

echo "Performing second contribution"
./setup p2c output_1.ph2 output_2.ph2
echo ""

echo "Verifying second contribution"
./setup p2v output_2.ph2 output_1.ph2
echo ""


echo "Performing third contribution"
./setup p2c output_2.ph2 output_3.ph2
echo ""

echo "Verifying third contribution"
./setup p2v output_3.ph2 output_2.ph2
echo ""

rm output_0.ph2
rm output_1.ph2
rm output_2.ph2
rm setup
