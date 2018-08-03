package main

import (
	"flag"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/pkg/errors"
)

// Holds the paths requested by the user for the profiles to be saved
var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memprofile = flag.String("memprofile", "", "write memory profile to this file")

// HandleProfiling is the introduction function for all profiling related functionality
func HandleProfiling() error {
	fmt.Println("wtf")
	flag.Parse()
	if *cpuprofile != "" {
		err := cpuProfile()
		if err != nil {
			return errors.Wrap(err, "failed to start CPU profiling")
		}
		//fmt.Println("Finished CPU profiling")
	}
	if *memprofile != "" {
		err := ramProfile()
		if err != nil {
			return errors.Wrap(err, "failed to start memory profiling")
		}
		//fmt.Println("Finished memory profiling")
	}
	return nil
}

// Handles situation if cpu profiling is enabled
func cpuProfile() error {
	f, err := os.Create(*cpuprofile)
	if err != nil {
		return errors.Wrap(err, "failed to create specified cpu profile file")
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		return errors.Wrap(err, "failed to start cpu profiling")
	}
	return nil

}

var memFile *os.File

// Handles situation if memory profiling is enabled
func ramProfile() error {
	var err error
	memFile, err = os.Create(*memprofile)
	if err != nil {
		return errors.Wrap(err, "Failed to create specified memory profile file")
	}
	err = pprof.WriteHeapProfile(memFile)
	if err != nil {
		return errors.Wrap(err, "Failed to start memory profiling")
	}
	return nil

}
