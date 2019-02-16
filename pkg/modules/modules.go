package modules

import (
	"github.com/autom8ter/fsctl"
	"github.com/autom8ter/fsctl/clone"
	"github.com/autom8ter/goexec/pkg/util"
	"io"
)

func init() {
	var err error
	Fs, err = fsctl.NewFs(clone.NewAuthCloner())
	if err != nil {
		panic(err)
	}
}

var (
	Fs  *fsctl.Fs
	app string
)

func createInit() io.Writer {
	app = Fs.Prompt("app", util.GreenStringf("please provide a name for your goexec program ---> "))
	err := Fs.MkdirAll(app, 0755)
	util.Panic(err, "failed to create goexec directory: %s\n", app)
	f, err := Fs.Create(app + "/main.go")
	util.PrintErr(err, "failed to create file: %s\n", app+"/main.go")
	return f
}
