package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LoliGothick/cbt/internal/core/solutions"
	"github.com/LoliGothick/cbt/internal/wandbox"
	"github.com/LoliGothick/cbt/internal/wandbox/expand"
	"github.com/LoliGothick/freyja/cutil"
	"github.com/LoliGothick/freyja/maybe"
	"github.com/mattn/go-colorable"
	"github.com/urfave/cli"
)

type CLI struct {
	app *cli.App
}

func NewCLI() *CLI {
	_cli := new(CLI)

	app := cli.NewApp()
	app.Name = "cbt"
	app.Usage = "Build Tool Wandbox"
	app.Version = "0.3.0"
	app.Commands = []cli.Command{
		{
			Name:    "wandbox",
			Aliases: []string{"wb"},
			Usage:   "Sending local codes to wandbox",
			Subcommands: []cli.Command{
				{
					Name:   "cpp",
					Usage:  "C++",
					Action: WandboxCpp,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "compiler, x",
							Usage: "specify C++ compiler [fmt: gcc-x.x.x, clang-x.x.x or zapcc]",
							Value: "gcc-head",
						},
						cli.StringFlag{
							Name:  "std",
							Usage: "specify C++ version [fmt: c++** or gnu++**]",
							Value: "c++14",
						},
						cli.StringFlag{
							Name:  "stdin,in",
							Usage: "specify standard input [text or file both accept]",
							Value: "",
						},
						cli.BoolFlag{
							Name:  "warning, w",
							Usage: `enable warning options [warning means "-Wall -Wextra"]`,
						},
						cli.StringFlag{
							Name:  "pedantic, p",
							Usage: "specify pedantic option [no(default), yes or eoors]",
							Value: "no",
						},
						cli.BoolFlag{
							Name:  "verbose, v",
							Usage: "enable verbose option",
						},
						cli.BoolFlag{
							Name:  "optimize, o",
							Usage: `enable optimize options [optimize means "-O2 -march=native"]`,
						},
						cli.BoolFlag{
							Name:  "sprout",
							Usage: "enable sprout library",
						},
						cli.StringFlag{
							Name:  "boost",
							Usage: "specify boost library version [fmt: x.xx.x]",
							Value: "nothing",
						},
						cli.BoolFlag{
							Name:  "msgpack, m",
							Usage: "enable massage pack",
						},
						cli.StringFlag{
							Name:  "compile-option, c",
							Usage: "specify compile options",
							Value: "",
						},
						cli.StringFlag{
							Name:  "runtime-option, r",
							Usage: "specify runtime options",
							Value: "",
						},
						cli.BoolFlag{
							Name:  "save, s",
							Usage: "publishing permanent link",
						},
						cli.BoolFlag{
							Name:  "bash",
							Usage: "Secret",
						},
					},
				},
				{
					Name:   "c",
					Usage:  "C",
					Action: WandboxC,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "compiler, x",
							Usage: "specify C++ compiler [fmt: gcc-x.x.x or clang-x.x.x]",
							Value: "gcc-head",
						},
						cli.StringFlag{
							Name:  "std",
							Usage: "specify C++ version [c89, gnu89, c99, gnu99, c11(default), or gnu11]",
							Value: "c11",
						},
						cli.StringFlag{
							Name:  "stdin,in",
							Usage: "specify standard input [text or file both accept]",
							Value: "",
						},
						cli.BoolFlag{
							Name:  "warning, w",
							Usage: `enable warning options [warning means "-Wall -Wextra"]`,
						},
						cli.StringFlag{
							Name:  "pedantic, p",
							Usage: "specify pedantic option [no(default), yes or eoors]",
							Value: "no",
						},
						cli.BoolFlag{
							Name:  "verbose, v",
							Usage: "enable verbose option",
						},
						cli.BoolFlag{
							Name:  "optimize, o",
							Usage: "enable optimize option",
						},
						cli.StringFlag{
							Name:  "compile-option, c",
							Usage: "specify compile options",
							Value: "",
						},
						cli.StringFlag{
							Name:  "runtime-option, r",
							Usage: "specify runtime options",
							Value: "",
						},
						cli.BoolFlag{
							Name:  "save, s",
							Usage: "publishing permanent link",
						},
						cli.BoolFlag{
							Name:  "bash",
							Usage: "Secret",
						},
					},
				},
				{
					Name:   "go",
					Usage:  "Go",
					Action: WandboxGo,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "compiler, x",
							Usage: "specify Go version [fmt: go-x.x]",
							Value: "go-head",
						},
						cli.StringFlag{
							Name:  "runtime-option, r",
							Usage: "specify runtime options",
							Value: "",
						},
						cli.StringFlag{
							Name:  "stdin,in",
							Usage: "specify standard input [text or file both accept]",
							Value: "",
						},
						cli.BoolFlag{
							Name:  "save, s",
							Usage: "publishing permanent link",
						},
					},
				},
				{
					Name:   "ruby",
					Usage:  "Ruby",
					Action: WandboxRuby,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "compiler, x",
							Usage: "specify Go version [fmt: ruby-x.x or mruby-x.x]",
							Value: "ruby-head",
						},
						cli.StringFlag{
							Name:  "runtime-option, r",
							Usage: "specify runtime options",
							Value: "",
						},
						cli.StringFlag{
							Name:  "stdin,in",
							Usage: "specify standard input [text or file both accept]",
							Value: "",
						},
						cli.BoolFlag{
							Name:  "save, s",
							Usage: "publishing permanent link",
						},
					},
				},
			},
		},
		{
			Name:    "solution",
			Aliases: []string{"sln"},
			Usage:   "solution management",
			Subcommands: []cli.Command{
				{
					Name:   "init",
					Usage:  "Initialize solution",
					Action: solutions.SolutionInitial,
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "lang",
							Value: "cpp",
						},
					},
				},
				{
					Name:   "update",
					Usage:  "update solution",
					Action: solutions.SolutionUpdate,
				},
			},
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Println("(´･_･`)? Command not found")
		return nil
	}
	_cli.app = app
	return _cli
}

func (_cli *CLI) Run() {
	_cli.app.Writer = colorable.NewColorableStdout()
	_cli.app.Run(os.Args)
}

func (_cli *CLI) TestRun(args []string) ([]byte, error) {
	outStream, errStream := new(bytes.Buffer), new(bytes.Buffer)
	_cli.app.Writer = outStream
	_cli.app.ErrWriter = errStream
	err := _cli.app.Run(args)
	return outStream.Bytes(), err
}

func WandboxC(c *cli.Context) {
	// preprocessing

	// prepare JSON struct
	config := wandbox.Request{}
	// prepare stdin
	var stdin string
	switch in := cutil.OrElse(c.String("in") == "", "", maybe.Expected(ioutil.ReadFile(c.String("in"))).UnwrapOr(c.String("in"))); in.(type) {
	case []byte:
		stdin = string(in.([]byte))
	case string:
		stdin = in.(string)
	case error:
		panic(in.(error))
	}

	// Let's Making JSON!
	if !c.Bool("bash") {
		options := c.String("std")
		if c.Bool("w") {
			options += ",warning"
		}
		switch c.String("p") {
		case "no":
			options += ",cpp-no-pedantic"
		case "yes":
			options += ",cpp-pedantic"
		case "errors":
			options += ",cpp-pedantic-errors"
		}
		if c.Bool("v") {
			options += ",cpp-verbose"
		}
		if c.Bool("o") {
			options += ",optimize"
		}
		if len(c.Args()) < 2 {
			code, codes := expand.ExpandInclude(string(c.Args().First()), `#include.*".*"|".*"/\*cbt-require\*/`)
			// JSON configure
			config = wandbox.Request{
				Compiler:          c.String("x") + "-c",
				Code:              code,
				Codes:             wandbox.TransformToCodes(codes),
				Options:           options,
				Stdin:             stdin,
				CompilerOptionRaw: c.String("c"),
				RuntimeOptionRaw:  c.String("r"),
				Save:              c.Bool("s"),
			}
		} else {
			targets := []string{}
			targets = c.Args()
			code, src, codes := expand.ExpandIncludeMulti(targets, `#include.*".*"|".*"/\*cbt-require\*/`)
			config = wandbox.Request{
				Compiler:          c.String("x") + "-c",
				Code:              code,
				Codes:             wandbox.TransformToCodes(codes),
				Options:           options,
				Stdin:             stdin,
				CompilerOptionRaw: strings.Join(src, "\n") + "\n" + c.String("c"),
				RuntimeOptionRaw:  c.String("r"),
				Save:              c.Bool("s"),
			}
		}
	} else {
		{ // else target is multiple src-file
			// set target
			var target expand.PathSlice
			target = ([]string)(c.Args())
			// code analyze
			codes := expand.ExpandAll(target, `#include.*".*"|".*"/\*cbt-require\*/`)
			// generate template (shell)
			shell_tmpl := `
echo 'compiler:' {{.Compiler}}
echo 'target:' {{.Target}}
{{if .Clang}}
/opt/wandbox/{{.Compiler}}/bin/clang {{.Target}} {{.Option}} && ./a.out{{else}}/opt/wandbox/{{.Compiler}}/bin/gcc {{.Target}} -std={{.VER}} {{.Option}} && ./a.out{{end}}{{if .StdinFlag}} <<- EOS
{{.Stdin}}
EOS{{end}}
`
			options := ""
			if c.Bool("w") {
				options += ` -Wall -Wextra`
			}
			switch c.String("p") {
			case "no":
			case "yes":
				options += " -pedantic"
			case "errors":
				options += " -pedantic-errors"
			}
			if c.Bool("v") {
				options += " -v"
			}
			if c.Bool("o") {
				options += " -O2 -march=native"
			}

			tmpl := template.Must(template.New("bash").Parse(shell_tmpl))
			bash := &wandbox.Bash{
				Compiler:  c.String("x"),
				Target:    strings.Join(target.ToBase(), " "),
				VER:       c.String("std"),
				Option:    options,
				StdinFlag: stdin != "",
				Stdin:     stdin,
				Clang:     c.String("x")[0:3] != "gcc",
			}
			var shell = ""
			buf := bytes.NewBufferString(shell)
			err := tmpl.Execute(buf, bash)
			if err != nil {
				panic(err)
			}
			// JSON configure
			config = wandbox.Request{
				Compiler: "bash",
				Code:     buf.String(),
				Codes:    wandbox.TransformToCodes(codes),
				Save:     c.Bool("s"),
			}
		}
	}
	postRequest(config, c.Bool("s"), c.App.Writer, c.App.ErrWriter)
}

func WandboxCpp(c *cli.Context) {
	// preprocessing

	// prepare JSON struct
	config := wandbox.Request{}
	// prepare stdin
	var stdin string
	switch in := cutil.OrElse(c.String("in") == "", "", maybe.Expected(ioutil.ReadFile(c.String("in"))).UnwrapOr(c.String("in"))); in.(type) {
	case []byte:
		stdin = string(in.([]byte))
	case string:
		stdin = in.(string)
	case error:
		panic(in.(error))
	}

	// Let's Making JSON!
	if !c.Bool("bash") {
		options := c.String("std")
		if c.Bool("w") {
			options += ",warning"
		}
		switch c.String("p") {
		case "no":
			options += ",cpp-no-pedantic"
		case "yes":
			options += ",cpp-pedantic"
		case "errors":
			options += ",cpp-pedantic-errors"
		}
		if c.Bool("v") {
			options += ",cpp-verbose"
		}
		if c.Bool("o") {
			options += ",optimize"
		}
		if c.Bool("sprout") {
			options += ",sprout"
		}
		if c.Bool("msgpack") {
			options += ",msgpack"
		}
		options += ",boost-" + c.String("boost") + "-" + c.String("x")
		if len(c.Args()) < 2 {
			code, codes := expand.ExpandInclude(string(c.Args().First()), `#include.*".*"|".*"/\*cbt-require\*/`)
			// JSON configure
			config = wandbox.Request{
				Compiler:          c.String("x"),
				Code:              code,
				Codes:             wandbox.TransformToCodes(codes),
				Options:           options,
				Stdin:             string(stdin),
				CompilerOptionRaw: c.String("c"),
				RuntimeOptionRaw:  c.String("r"),
				Save:              c.Bool("s"),
			}
		} else {
			targets := []string{}
			targets = c.Args()
			code, src, codes := expand.ExpandIncludeMulti(targets, `#include.*".*"|".*"/\*cbt-require\*/`)
			config = wandbox.Request{
				Compiler:          c.String("x"),
				Code:              code,
				Codes:             wandbox.TransformToCodes(codes),
				Options:           options,
				Stdin:             string(stdin),
				CompilerOptionRaw: strings.Join(src, "\n") + "\n" + c.String("c"),
				RuntimeOptionRaw:  c.String("r"),
				Save:              c.Bool("s"),
			}
		}
	} else {
		{ // else target is multiple src-file
			// set target
			var target expand.PathSlice
			target = ([]string)(c.Args())
			// code analyze
			codes := expand.ExpandAll(target, `#include.*".*"|".*"/\*cbt-require\*/`)
			// generate template (shell)
			shell_tmpl := `
echo 'compiler:' {{.Compiler}}
echo 'target:' {{.Target}}
{{if .Clang}}
/opt/wandbox/{{.Compiler}}/bin/clang++ {{.Target}} {{.Option}} && ./a.out{{else}}/opt/wandbox/{{.Compiler}}/bin/g++ {{.Target}} -std={{.CXX}}++{{.VER}} {{.Option}} && ./a.out{{end}}{{if .StdinFlag}} <<- EOS
{{.Stdin}}
EOS{{end}}
`
			cxx := strings.Split(c.String("std"), "++")
			options := ""
			if c.Bool("w") {
				options += ` -Wall -Wextra`
			}
			switch c.String("p") {
			case "no":
			case "yes":
				options += " -pedantic"
			case "errors":
				options += " -pedantic-errors"
			}
			if c.Bool("v") {
				options += " -v"
			}
			if c.Bool("o") {
				options += " -O2 -march=native"
			}
			if c.Bool("sprout") {
				options += " -I/opt/wandbox/sprout"
			}
			if c.Bool("msgpack") {
				options += " -I/opt/wandbox/msgpack/include"
			}
			if c.String("boost") != "nothing" {
				options += " -I/opt/wandbox/boost-" + c.String("boost") + "/" + c.String("x") + "/include"
			}

			tmpl := template.Must(template.New("bash").Parse(shell_tmpl))
			bash := &wandbox.Bash{
				Compiler:  c.String("x"),
				Target:    strings.Join(target.ToBase(), " "),
				CXX:       cxx[0],
				VER:       cxx[1],
				Option:    options,
				StdinFlag: string(stdin) != "",
				Stdin:     string(stdin),
				Clang:     c.String("x")[0:3] != "gcc",
			}
			var shell = ""
			buf := bytes.NewBufferString(shell)
			err := tmpl.Execute(buf, bash)
			if err != nil {
				panic(err)
			}
			// JSON configure
			config = wandbox.Request{
				Compiler: "bash",
				Code:     buf.String(),
				Codes:    wandbox.TransformToCodes(codes),
				Save:     c.Bool("s"),
			}
		}
	}
	postRequest(config, c.Bool("s"), c.App.Writer, c.App.ErrWriter)
}

func WandboxGo(c *cli.Context) {
	// preprocessing

	// prepare JSON struct
	config := wandbox.Request{}
	// prepare stdin
	var stdin string
	switch in := cutil.OrElse(c.String("in") == "", "", maybe.Expected(ioutil.ReadFile(c.String("in"))).UnwrapOr(c.String("in"))); in.(type) {
	case []byte:
		stdin = string(in.([]byte))
	case string:
		stdin = in.(string)
	case error:
		panic(in.(error))
	}

	code, codes := expand.ExpandGo(string(c.Args().First()))
	// JSON configure
	config = wandbox.Request{
		Compiler:         c.String("x"),
		Code:             code,
		Codes:            wandbox.TransformToCodes(codes),
		Stdin:            stdin,
		RuntimeOptionRaw: c.String("r"),
		Save:             c.Bool("s"),
	}

	postRequest(config, c.Bool("s"), c.App.Writer, c.App.ErrWriter)
}

func WandboxRuby(c *cli.Context) {
	// preprocessing

	// prepare JSON struct
	config := wandbox.Request{}
	// prepare stdin
	var stdin string
	switch in := cutil.OrElse(c.String("in") == "", "", maybe.Expected(ioutil.ReadFile(c.String("in"))).UnwrapOr(c.String("in"))); in.(type) {
	case []byte:
		stdin = string(in.([]byte))
	case string:
		stdin = in.(string)
	case error:
		panic(in.(error))
	}

	// Let's Making JSON!
		if len(c.Args()) < 2 {
			code, codes := expand.ExpandRubyRequire(string(c.Args().First()), `require.*'.*'`)
			// JSON configure
			config = wandbox.Request{
				Compiler:          c.String("x"),
				Code:              code,
				Codes:             wandbox.TransformToCodes(codes),
				Stdin:             stdin,
				RuntimeOptionRaw:  c.String("r"),
				Save:              c.Bool("s"),
			}
		}

	postRequest(config, c.Bool("s"), c.App.Writer, c.App.ErrWriter)
}



func postRequest(config wandbox.Request, save bool, stdout, stderr io.Writer) *wandbox.Result {
	// Marshal JSON
	cppJSONBytes, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	out := new(bytes.Buffer)
	json.Indent(out, cppJSONBytes, "", "    ") // pretty

	// Client : Wait Time 30s
	client := &http.Client{Timeout: time.Duration(30) * time.Second}
	// Request : POST JSON
	req, err := http.NewRequest("POST", "https://wandbox.org/api/compile.json?", strings.NewReader(out.String()))
	if err != nil {
		panic(err)
	}
	// Header : Content-type <- application/json
	req.Header.Add("Content-type", "application/json")

	// Send POST
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	result := new(wandbox.Result)
	if err := json.Unmarshal(([]byte)(body), result); err != nil {
		panic(fmt.Errorf(`%s:\n%s`, err.Error(), body))
	}

	switch {
	case result.ProgramMessage != "":
		stdout.Write([]byte(result.ProgramMessage))
	case result.CompilerError != "":
		stdout.Write([]byte("\033[31m"))
		stdout.Write([]byte("Compilation Error!:"))
		stdout.Write([]byte(result.CompilerError))
	case result.ProgramError != "":
		stdout.Write([]byte("Runtime Error!:"))
		stdout.Write([]byte(result.ProgramError))
	}
	stdout.Write([]byte("\033[0m"))
	if save {
		stdout.Write([]byte("\n\nPermlink: " + result.Permlink))
		stdout.Write([]byte("\nURL: " + result.URL))
	}
	return result
}
