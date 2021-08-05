package main

import (
	"context"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"

	nested "github.com/antonfisher/nested-logrus-formatter"

	"idiocy/action"

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
	// var err error
	var matches []string

	// dir := c.String("dir")
	// if !strings.EqualFold(dir, "") {
	// 	if !common.IsDir(dir) {
	// 		return fmt.Errorf("%s doesn't exist", dir)
	// 	}
	// 	if matches, err = filepath.Glob(path.Join(dir, "*.go")); err != nil {
	// 		return err
	// 	}
	// }

	logrus.Print(matches)
	files := c.StringSlice("file")
	if len(files) != 0 {
		matches = append(matches, files...)
	}

	for _, v := range matches {
		action.GenerateDoc(v)
	}

	return nil
}
