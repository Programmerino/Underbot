package profiling

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
func HandleProfiling() {
	flag.Parse()
	if *cpuprofile != "" {
		cpuProfile()
		fmt.Println("Finished CPU profiling")
	}
	if *memprofile != "" {
		ramProfile()
		fmt.Println("Finished memory profiling")
	}
}

// Handles situation if cpu profiling is enabled
func cpuProfile() {
	f, err := os.Create(*cpuprofile)
	if err != nil {
		panic(errors.Wrap(err, "Failed to create specified CPU profile file"))
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		panic(errors.Wrap(err, "Failed to start CPU profiling"))
	}
	defer pprof.StopCPUProfile()
}

// Handles situation if memory profiling is enabled
func ramProfile() {
	f, err := os.Create(*memprofile)
	if err != nil {
		panic(errors.Wrap(err, "Failed to create specified memory profile file"))
	}
	err = pprof.WriteHeapProfile(f)
	if err != nil {
		panic(errors.Wrap(err, "Failed to start memory profiling"))
	}
	defer f.Close()
}
