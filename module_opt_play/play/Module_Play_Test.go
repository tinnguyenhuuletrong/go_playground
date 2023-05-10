package play

import (
	"log"

	"ttin.com/play2022/module_opt_play"
)

func Play() {
	ins := module_opt_play.NewModuleX(
		module_opt_play.WithAddress("0.0.0.0"),
		module_opt_play.WithSSL(true),
		module_opt_play.WithPort(433))
	log.Println(ins)
}
