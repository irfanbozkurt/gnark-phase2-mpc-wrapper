# #!/bin/bash

if [ -z "$1" ]; then
    echo "Error: No input .ph2 file provided"
    exit 1
fi

to_be_verified="$1"
if [ ! -f "$to_be_verified" ]; then
    echo "Error: The file '$to_be_verified' does not exist"
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

if [ -n "$2" ]; then
    output_of_prev_contribution="$2"
else
    output_of_prev_contribution="output_0.ph2"
fi

echo "Verifying given contribution"
./setup p2v $to_be_verified $output_of_prev_contribution
echo ""

rm setup
