package main

import (
	"fmt"
	"time"

	"github.com/terminalhq/terminal/go/pkg/particles"
)

func main() {
    params := particles.ParticleParams{
        Ascii: []particles.ParticleAscii{
            {Count: 1, Output: []string{" "}},
            {Count: 2, Output: []string{"."}},
            {Count: 4, Output: []string{":"}},
            {Count: 6, Output: []string{"{", "}"}},
        },
        MaxX: 6,
        MaxLife: 2000,
        Height: 5,
        MaxSpeed: 4.75,
    }

    steam := particles.NewCoffeeSteam(params, 100)
    steam.Start()
    ticker := time.NewTicker(time.Millisecond * 100)
    for {
        <-ticker.C
        steam.Update()
        fmt.Print("\033[H\033[2J")
        fmt.Println(steam.Display())
    }
}


