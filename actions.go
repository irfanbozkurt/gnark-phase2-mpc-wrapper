package main

import (
	"errors"
	"strconv"

	"github.com/bnbchain/zkbnb-setup/phase1"
	"github.com/bnbchain/zkbnb-setup/phase2"
	"github.com/bnbchain/zkbnb-setup/keys"
	"github.com/urfave/cli/v2"
)

func p1n(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	powerStr := cCtx.Args().Get(0)
	power, err := strconv.Atoi(powerStr)
	if err != nil {
		return err
	}
	if power > 26 {
		return errors.New("can't support powers larger than 26")
	}
	outputPath := cCtx.Args().Get(1)
	err = phase1.Initialize(byte(power), outputPath)
	return err
}

func p1c(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	outputPath := cCtx.Args().Get(1)
	err := phase1.Contribute(inputPath, outputPath)
	return err
}

func p1v(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 1 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	err := phase1.Verify(inputPath)
	return err
}

func p2n(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 4 {
		return errors.New("please provide the correct arguments")
	}
	
	phase1Path := cCtx.Args().Get(0)
	r1csPath := cCtx.Args().Get(1)
	phase2Path := cCtx.Args().Get(2)
	evalsPath := cCtx.Args().Get(3)
	err := phase2.Initialize(phase1Path, r1csPath, phase2Path, evalsPath)
	return err
}

func p2c(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	outputPath := cCtx.Args().Get(1)
	err := phase2.Contribute(inputPath, outputPath)
	return err
}

func p2v(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	originPath := cCtx.Args().Get(1)
	err := phase2.Verify(inputPath, originPath)
	return err
}

func extract(cCtx *cli.Context) error {
	// sanity check
	if cCtx.Args().Len() != 2 {
		return errors.New("please provide the correct arguments")
	}
	inputPath := cCtx.Args().Get(0)
	evalsPath := cCtx.Args().Get(1)
	err := keys.ExtractKeys(inputPath, evalsPath)
	return err
}