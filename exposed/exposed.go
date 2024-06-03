package exposed

import (
	"os"

	"github.com/consensys/gnark/backend/groth16/bn254/mpcsetup"
)

func P2v(inputPath string, originPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	input := &mpcsetup.Phase2{}
	input.ReadFrom(inputFile)

	originFile, err := os.Open(originPath)
	if err != nil {
		return err
	}
	origin := &mpcsetup.Phase2{}
	origin.ReadFrom(originFile)

	err = mpcsetup.VerifyPhase2(origin, input)
	if err != nil {
		return err
	}

	return nil
}
