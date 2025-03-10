package cmd

import (
	"fmt"

	"github.com/yagoyudi/cheat/internal/config"
)

func Conf(_ map[string]interface{}, conf config.Config) {
	fmt.Println(conf.Path)
}
