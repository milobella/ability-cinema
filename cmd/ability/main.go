package main

import (
    "milobella/abilities/ability-sdk-go/pkg/ability"
    "milobella/oratio/pkg/anima"
    "milobella/oratio/pkg/cerebro"
)

// fun main()
func main() {
    server := ability.NewServer(10200)
    server.RegisterIntent("hello_world", helloWorldHandler)
    server.Serve()
}

func helloWorldHandler(nlu cerebro.NLU, nlg *anima.NLG) {
    nlg.Sentence = "Hello {{world}}"
    nlg.Params = map[string]string{"world": "world"}
}
