package main

import (
	"context"
	"fmt"
	"idiocy/logger"
	"idiocy/schema"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	nested "github.com/antonfisher/nested-logrus-formatter"

	"github.com/urfave/cli/v2"
)

// Foo 结构体
type Foo struct {
	i int
}

// Bar 接口
type Bar interface {
	Do(ctx context.Context) error
}

func init() {
	logrus.SetFormatter(&nested.Formatter{
		TrimMessages:    true,
		TimestampFormat: "15:04:05",
		NoFieldsSpace:   true,
		HideKeys:        false,
		ShowFullLevel:   true,
		CallerFirst:     false,
		FieldsOrder:     []string{"component", "category"},
		CustomCallerFormatter: func(f *runtime.Frame) string {
			s := strings.Split(f.Function, ".")
			funcName := s[len(s)-1]
			return fmt.Sprintf(" [%s:%d][%s()]", path.Base(f.File), f.Line, funcName)
		},
	})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.TraceLevel)
}

func main() {
	app := &cli.App{
		Name: "idiocy",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "conf",
				Aliases: []string{`c`},
				Value:   "./config.toml",
			},
		},
		Commands: []*cli.Command{
			&genDoc,
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}

}

var genDoc = cli.Command{
	Name:  "gen",
	Usage: "generate doc",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "dir",
			Aliases: []string{`d`},
		},
		&cli.StringSliceFlag{
			Name:    "file",
			Aliases: []string{`f`},
		},
	},
	Action: run,
}

func run(c *cli.Context) error {
	projSchema := schema.NewSchema(c.String("dir"))
	projSchema.LoadSourceFiles()

	for _, v := range projSchema.SourceFile {
		v.ParseFile(v.FullPath)
		logger.S.Infof("%#v", v)

		v.BuildStacks()

		// f.BuildStacks()
		// f.FindDecals()
		// f.FindGinInstance()

	}

	return nil
}
