package creature

import (
	"math/rand"
)

type Creature struct {
	HoleLength int
	Health     int
	Respect    int
	Weight     int
}

func (c *Creature) Dig(intensively bool) {

	if intensively {
		c.HoleLength += 5
		c.Health -= 30
	} else {
		c.HoleLength += 2
		c.Health -= 10
	}
}

func (c *Creature) Eat(green bool) {
	if green {
		if c.Respect < 30 {
			c.Health -= 30
		} else {
			c.Health += 30
			c.Weight += 30
		}
	} else {
		c.Health += 10
		c.Weight += 15
	}
}

func (c *Creature) FightWith(enemy Creature) bool {
	result := false
	x := float32(enemy.Weight) / float32(c.Weight)
	if rand.Intn(c.Weight+enemy.Weight+1) <= c.Weight {
		result = true
		c.Respect += int(50 * x)
	}
	c.Health -= int(50 * x)

	return result
}

func (c *Creature) Sleep() {
	c.HoleLength -= 2
	c.Health += 20
	c.Respect -= 2
	c.Weight -= 5
}
