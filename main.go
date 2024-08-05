package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

const (
	ArgMode = iota
	ArgInFile
	ArgVarName
	ArgWorkFile
)

func main() {
	os.Exit(rMain())
}

//go:embed help.txt
var helpText string

func rMain() int {
	var err error

	var dump []UefiVar
	var offset int = -1
	var tFile []byte

	flag.Parse()

	if flag.NArg() != 4 {
		fmt.Println("invalid number of arguments")
		fmt.Println(helpText)
		return 1
	}

	tFile, err = os.ReadFile(flag.Arg(ArgInFile))
	if err != nil {
		fmt.Printf("error opening file %s: %v\n", flag.Arg(ArgInFile), err)
		return 1
	}

	err = json.Unmarshal(tFile, &dump)
	if err != nil {
		fmt.Printf("error parsing dump file %s: %v\n", flag.Arg(ArgInFile), err)
		return 1
	}

	//find variable
	for i := range dump {
		if dump[i].Name == flag.Arg(ArgVarName) {
			offset = i
			break
		}
	}
	if offset == -1 {
		fmt.Printf("could not find variable %s\n", flag.Arg(ArgVarName))
		return 1
	}

	switch flag.Arg(ArgMode) {
	case "r":
		err = os.WriteFile(flag.Arg(ArgWorkFile), dump[offset].Data, 0o660)
		if err != nil {
			fmt.Printf("error writing to output file %s: %v\n", flag.Arg(ArgWorkFile), err)
			return 1
		}
	case "w":
		tFile, err = os.ReadFile(flag.Arg(ArgWorkFile))
		if err != nil {
			fmt.Printf("error reading from input file %s: %v\n", flag.Arg(ArgWorkFile), err)
			return 1
		}

		dump[offset].DataLen = len(tFile)
		dump[offset].Data = tFile

		tFile, err = json.MarshalIndent(dump, "", "\t")
		if err != nil {
			fmt.Printf("error marshalling dump: %v\n", err)
			return 1
		}

		err = os.WriteFile(flag.Arg(ArgInFile), tFile, 0o660)
		if err != nil {
			fmt.Printf("error writing re-marshalled dump to output file %s: %v\n", flag.Arg(ArgWorkFile), err)
			return 1
		}
	default:
		fmt.Printf("invalid mode %d\n", flag.Arg(ArgMode))
		println(helpText)
		return 1
	}

	return 0
}

type UefiVar struct {
	Name       string `json:"name"`
	VendorGuid string `json:"vendor_guid"`

	Attributes struct {
		NonVolatile                       bool `json:"non_volatile"`
		BootserviceAccess                 bool `json:"bootservice_access"`
		RuntimeAccess                     bool `json:"runtime_access"`
		HardwareErrorRecord               bool `json:"hardware_error_record"`
		AuthenticatedWriteAccess          bool `json:"authenticated_write_access"`
		TimeBasedAuthenticatedWriteAccess bool `json:"time_based_authenticated_write_access"`
		EnhancedAuthenticatedAccess       bool `json:"enhanced_authenticated_access"`
	} `json:"attributes"`

	DataLen int    `json:"data_len"`
	Data    []byte `json:"data"`
}
