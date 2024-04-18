package particles

import (
	"math"
	"math/rand"
	"slices"
	"strings"
	"time"
)

type Particle struct {
	lifetime int64
	speed    float64

	x float64
	y float64
}

func (p *Particle) Reset(params *ParticleParams) {
	p.lifetime = int64(math.Floor(float64(params.MaxLife) * rand.Float64()))
    speed := params.MaxSpeed * rand.Float64()
	p.speed = speed

	xPos := rand.NormFloat64()
	if xPos > 4 {
		xPos = 4
	} else if xPos < -4 {
		xPos = -4
	}

	p.x = xPos / 4 * float64(params.MaxX)
	p.y = 0

}

func (p *Particle) eol() bool {
	return p.lifetime <= 0
}

func (p *Particle) update(delta int64) {
	p.lifetime -= delta
	if p.lifetime <= 0 {
		return
	}

    increase := p.speed * (float64(delta) / 1000.0)
	p.y += increase
}

type ParticleAscii struct {
	Count  int
	Output []string
}

type ParticleParams struct {
	Ascii    []ParticleAscii
	MaxLife  int64
	MaxSpeed float64
	MaxX     int
    Height   int
}

type CoffeeSteam struct {
	ParticleParams

	count      int
	particles  []*Particle
	lastUpdate int64
}

func NewCoffeeSteam(params ParticleParams, count int) CoffeeSteam {
	particles := make([]*Particle, 0, count)

	return CoffeeSteam{
		ParticleParams: params,
		count:          count,
		particles:      particles,
	}
}

func (c *CoffeeSteam) Start() {
	c.lastUpdate = time.Now().UnixMilli()

    for i := 0; i < c.count; i++ {
        p := Particle{}
        c.particles = append(c.particles, &p)
		p.Reset(&c.ParticleParams)
	}
}

func (c *CoffeeSteam) Update() {
	now := time.Now().UnixMilli()
	delta := now - c.lastUpdate
	c.lastUpdate = now

	for _, p := range c.particles {
		if p.eol() {
			p.Reset(&c.ParticleParams)
		} else {
			p.update(delta)
		}
	}
}

func zero(length int) []int {
    out := make([]int, 0, length)
    for i := 0; i < length; i++ {
        out = append(out, 0)
    }
    return out
}

func (c *CoffeeSteam) toValue(count int) string {
    out := " "
    for _, c := range c.ParticleParams.Ascii {
        if c.Count > count {
            break
        }

        idx := int(math.Floor(rand.Float64() * float64(len(c.Output))))
        out = c.Output[idx]
    }
    return out
}

func (c *CoffeeSteam) Display() string {
    rows := make([][]int, 0)
    maxRowCount := c.ParticleParams.MaxX * 2 + 1

    for _, p := range c.particles {
        row := int(math.Floor(p.y))
        col := c.ParticleParams.MaxX + int(math.Floor(p.x))

        for len(rows) <= row {
            rows = append(rows, zero(maxRowCount))
        }

        rows[row][col] += 1
    }

    out := make([]string, 0)
    for _, row := range rows {
        r := ""
        for _, count := range row {
            r += c.toValue(count)
        }
        if len(strings.TrimSpace(r)) == 0 {
            continue
        }
        out = append(out, r)

        if len(out) >= c.ParticleParams.Height {
            break
        }
    }

    for len(out) < c.ParticleParams.Height {
        out = append(out, "")
    }
    slices.Reverse(out)
    return strings.Join(out, "\n")
}
