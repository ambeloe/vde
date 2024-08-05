package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
)

const (
	ArgInFile = iota
	ArgOutDir
)

func main() {
	os.Exit(rMain())
}

const helpText = `Var dump Extract All
usage: vdea {uefivardump json file} {output directory}`

func rMain() int {
	var err error

	var dump []UefiVar
	var fileName string
	var tFile []byte

	flag.Parse()

	if flag.NArg() != 2 {
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

	for _, d := range dump {
		fileName = fmt.Sprintf("%s_%s_nv-%t_bs-%t_rs-%t_%x.bin", d.Name, d.VendorGuid, d.Attributes.NonVolatile, d.Attributes.BootserviceAccess, d.Attributes.RuntimeAccess, d.DataLen)
		err = os.WriteFile(path.Join(flag.Arg(ArgOutDir), fileName), d.Data, 0644)
		if err != nil {
			fmt.Printf("error writing file %s: %v\n", fileName, err)
			return 1
		}
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
