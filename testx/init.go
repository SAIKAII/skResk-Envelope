package testx

import (
	"github.com/SAIKAII/skResk-Infra/base"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
)

func init() {
	file := kvs.GetCurrentFilePath("../../config/config.ini", 1)
	conf := ini.NewIniFileCompositeConfigSource(file)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DbxDatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	//infra.Register(&base.IrisServerStarter{})

	app := infra.New(conf)
	app.Start()
}
