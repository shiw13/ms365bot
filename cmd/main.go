package main

import (
	"log"

	"ms365bot/pkg/bot"

	"github.com/shiw13/go-one/pkg/startup"
)

type program struct{}

func main() {
	prg := &program{}

	if err := startup.Run(prg); err != nil {
		log.Fatalf("%s", err)
	}
}

func (p *program) Initialize() error {
	return bot.InitBot()
}

func (p *program) OnStart() error {
	bot.StartBot()
	return nil
}

func (p *program) OnStop() error {
	bot.StopBot()
	return nil
}
