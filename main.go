package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	groth16 "github.com/consensys/gnark/backend/groth16/bn254"
	"github.com/consensys/gnark/backend/groth16/bn254/mpcsetup"
	cs "github.com/consensys/gnark/constraint/bn254"
	deserializer "github.com/irfanbozkurt/ptau-deserializer/deserialize"
	"github.com/urfave/cli/v2"
)

func readPtauFileAsPh1(ptauFilePath string) (*mpcsetup.Phase1, error) {
	fmt.Fprintln(os.Stdout, []any{"\n################# Reading ptau file...\n"}...)
	ptau, err := deserializer.ReadPtau(ptauFilePath)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stdout, []any{"\n################# Converting ptau file into a phase1 object...\n"}...)
	_phase1, err := deserializer.ConvertPtauToPhase1(ptau)
	if err != nil {
		return nil, err
	}

	phase1 := &mpcsetup.Phase1{}

	// Implement the getters for Phase1 from ptau-deserializer
	phase1.Parameters.G1.Tau = _phase1.TauG1()
	phase1.Parameters.G1.AlphaTau = _phase1.AlphaTauG1()
	phase1.Parameters.G1.BetaTau = _phase1.BetaTauG1()
	phase1.Parameters.G2.Tau = _phase1.TauG2()
	phase1.Parameters.G2.Beta = _phase1.BetaG2()

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

	fmt.Println("\n################# Writing to file: " + phase2Path)
	phase2File, err := os.Create(phase2Path)
	if err != nil {
		return err
	}
	defer phase2File.Close()
	if _, err = phase2.WriteTo(phase2File); err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	phase2File.Close()

	fmt.Println("\n################# Writing to file: " + phase2EvaluationsPath)
	phase2EvaluationsFile, err := os.Create(phase2EvaluationsPath)
	if err != nil {
		return err
	}
	defer phase2EvaluationsFile.Close()
	if _, err := phase2Evaluations.WriteTo(phase2EvaluationsFile); err != nil {
		fmt.Println("Error writing to file:", err)
		return err
	}
	phase2EvaluationsFile.Close()

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
	evals := mpcsetup.Phase2Evaluations{}
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

	pk, vk := mpcsetup.ExtractKeys(phase1, phase2, &evals, r1cs.NbConstraints)
	fmt.Println("len(vk.G1.K)", len(vk.G1.K))

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

func main() {
	app := &cli.App{
		Name:      "setup",
		Usage:     "Use this tool to generate parameters of Groth16 via MPC",
		UsageText: "setup command [arguments...]",
		Commands: []*cli.Command{
			/* --------------------------- Phase 2 Initialize --------------------------- */
			{
				Name:        "p2n",
				Usage:       "p2n <ptauPath> <r1csPath> <phase2Path>",
				Description: "initialize phase 2 for the given circuit",
				Action:      p2n,
			},
			/* --------------------------- Phase 2 Contribute --------------------------- */
			{
				Name:        "p2c",
				Usage:       "p2c <inputPath> <outputPath>",
				Description: "contribute phase 2 randomness for Groth16",
				Action:      p2c,
			},
			/* ----------------------------- Phase 2 Verify ----------------------------- */
			{
				Name:        "p2v",
				Usage:       "p2v <inputPath> <originPath>",
				Description: "verify phase 2 contributions for Groth16",
				Action:      p2v,
			},
			/* ----------------------------- Keys Extraction ---------------------------- */
			{
				Name:        "extract-keys",
				Usage:       "extract-keys <phase1Path> <phase2Path> <phase2EvalsPath> <r1csPath>",
				Description: "extract proving and verifying keys",
				Action:      extractKeys,
			},
			{
				Name:        "sol",
				Usage:       "sol <verifyingKey>",
				Description: "export verifier smart contract from verifying key",
				Action:      exportSol,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
