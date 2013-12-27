package main

import (
	flags "github.com/jessevdk/go-flags"
	"log"
	"os"
	"runtime"
)

const version = "v0.0.1"

var opts struct {
	Version    bool   `short:"v" long:"version" description:"Shows gowsdl version"`
	Package    string `short:"p" long:"package" description:"Package under which code will be generated" default:"myservice"`
	OutputFile string `short:"o" long:"output" description:"File where the generated code will be saved" default:"myservice.go"`
}

var logger *log.Logger

func init() {
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.SetPrefix("🍀  ")
}

func main() {
	args, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		log.Println(version)
		os.Exit(0)
	}

	if len(args) == 0 {
		log.Fatalln("WSDL file is required to start the party")
	}

	if opts.OutputFile == args[0] {
		log.Fatalln("Output file cannot be the same WSDL file")
	}

	gowsdl, err := NewGoWsdl(args[0], opts.Package)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	gocode, err := gowsdl.Start()
	if err != nil {
		log.Fatalln(err)
	}

	pkg := "./" + opts.Package
	err = os.Mkdir(pkg, 0744)

	if perr, ok := err.(*os.PathError); ok && os.IsExist(perr.Err) {
		log.Println("Package directory already exist, skipping creation")
	} else {
		if err != nil {
			log.Fatalln(err)
		}
	}

	fd, err := os.Create(pkg + "/" + opts.OutputFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer fd.Close()

	fd.Write(gocode["types"])
	fd.Write(gocode["messages"])
	fd.Write(gocode["operations"])
	fd.Write(gocode["proxy"])

	log.Println("Done 💩")
}