package main

import (
	"errors"
	"fmt"
	"os"

	groth16 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/backend/groth16/bn254/mpcsetup"
	cs "github.com/consensys/gnark/constraint/bn254"
	"github.com/urfave/cli/v2"
	deserializer "github.com/worldcoin/ptau-deserializer/deserialize"
)

func readPtauFileAsPh1(ptauFilePath string) (*mpcsetup.Phase1, error) {
	fmt.Fprintln(os.Stdout, []any{"\n################# Reading ptau file...\n"}...)
	ptau, err := deserializer.ReadPtau(ptauFilePath)
	if err != nil {
		return nil, err
	}

	fmt.Printf("len(ptau.PTauPubKey.TauG1): %d\n", len(ptau.PTauPubKey.TauG1))
	fmt.Printf("len(ptau.PTauPubKey.AlphaTauG1): %d\n", len(ptau.PTauPubKey.AlphaTauG1))
	fmt.Printf("len(ptau.PTauPubKey.BetaTauG1): %d\n", len(ptau.PTauPubKey.BetaTauG1))

	fmt.Printf("len(ptau.PTauPubKey.TauG2): %d\n", len(ptau.PTauPubKey.TauG2))
	fmt.Printf("len(ptau.PTauPubKey.BetaG2): %v\n", ptau.PTauPubKey.BetaG2)

	fmt.Fprintln(os.Stdout, []any{"\n################# Converting ptau file into a phase1 object...\n"}...)
	_phase1, err := deserializer.ConvertPtauToPhase1(ptau)
	if err != nil {
		return nil, err
	}

	phase1 := &mpcsetup.Phase1{}

	phase1.Parameters.G1.Tau = _phase1.GetTauG1()
	phase1.Parameters.G1.AlphaTau = _phase1.GetAlphaTauG1()
	phase1.Parameters.G1.BetaTau = _phase1.GetBetaTauG1()
	phase1.Parameters.G2.Tau = _phase1.GetTauG2()
	phase1.Parameters.G2.Beta = _phase1.GetBetaG2()

	return phase1, nil
}

func p2n(cCtx *cli.Context) error {
	if cCtx.Args().Len() != 3 {
		return errors.New("please provide the correct arguments")
	}

	ptauFilePath := cCtx.Args().Get(0)
	r1csPath := cCtx.Args().Get(1)
	phase2Path := cCtx.Args().Get(2)
	phase2EvaluationsPath := "phase2Evaluations"

	phase1, err := readPtauFileAsPh1(ptauFilePath)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, []any{"\n################# Reading the r1cs file...\n"}...)
	r1csFile, err := os.Open(r1csPath)
	if err != nil {
		return err
	}
	r1cs := &cs.R1CS{}
	r1cs.ReadFrom(r1csFile)

	fmt.Fprintln(os.Stdout, []any{"\n################# Initializing phase2\n"}...)
	phase2, phase2Evaluations := mpcsetup.InitPhase2(r1cs, phase1)
	fmt.Println("phase2.Hash: ", phase2.Hash)

	fmt.Println("\n################# Writing to file: " + phase2Path)
	phase2File, err := os.Create(phase2Path)
	if err != nil {
		return err
	}
	phase2.WriteTo(phase2File)

	fmt.Println("\n################# Writing to file: " + phase2EvaluationsPath)
	phase2EvaluationsFile, err := os.Create(phase2EvaluationsPath)
	if err != nil {
		return err
	}
	phase2Evaluations.WriteTo(phase2EvaluationsFile)

	return nil
}

func p2c(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	outputPath := cCtx.Args().Get(1)

	inputFile, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	phase2 := &mpcsetup.Phase2{}
	phase2.ReadFrom(inputFile)

	phase2.Contribute()

	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	phase2.WriteTo(outputFile)

	fmt.Println("Contribution done successfully. Written to " + outputPath)

	return nil
}

func p2v(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	originPath := cCtx.Args().Get(1)

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

	fmt.Println("Phase 2 contributions verified successfully")

	return nil
}

func extractKeys(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 4 {
		return errors.New("please provide the correct arguments")
	}

	ptauPath := cCtx.Args().Get(0)
	phase1, err := readPtauFileAsPh1(ptauPath)
	if err != nil {
		return err
	}

	phase2Path := cCtx.Args().Get(1)
	phase2 := &mpcsetup.Phase2{}
	phase2File, err := os.Open(phase2Path)
	if err != nil {
		return err
	}
	phase2.ReadFrom(phase2File)

	evalsPath := cCtx.Args().Get(2)
	evals := &mpcsetup.Phase2Evaluations{}
	evalsFile, err := os.Open(evalsPath)
	if err != nil {
		return err
	}
	evals.ReadFrom(evalsFile)

	fmt.Println("len(evals.G1.VKK)", len(evals.G1.VKK))

	r1csPath := cCtx.Args().Get(3)
	r1cs := &cs.R1CS{}
	r1csFile, err := os.Open(r1csPath)
	if err != nil {
		return err
	}
	r1cs.ReadFrom(r1csFile)

	pk, vk := mpcsetup.ExtractKeys(phase1, phase2, evals, r1cs.NbConstraints)

	pkFile, err := os.Create("pk")
	if err != nil {
		return err
	}
	pk.WriteTo(pkFile)

	vkFile, err := os.Create("vk")
	if err != nil {
		return err
	}
	vk.WriteTo(vkFile)

	return nil
}

func exportSol(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 1 {
		return errors.New("please provide the correct arguments")
	}

	vkPath := cCtx.Args().Get(0)
	vk := &groth16.VerifyingKey{}
	vkFile, err := os.Open(vkPath)
	if err != nil {
		return err
	}
	vk.ReadFrom(vkFile)

	solFile, err := os.Create("verifier.sol")
	if err != nil {
		return err
	}

	err = vk.ExportSolidity(solFile)
	return err
}
