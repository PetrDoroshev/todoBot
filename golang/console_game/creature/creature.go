package creature

import (
	"math/rand"
)

type Creature struct {
	holeLength int
	health     int
	respect    int
	weight     int
}

func New(_holeLength int, _health int, _respect int, _weight int) *Creature {
	c := new(Creature)

	c.holeLength = _holeLength
	c.health = _health
	c.respect = _respect
	c.weight = _weight

	return c
}

func (c *Creature) IsStateOk() bool {
	return c.health > 0 && c.holeLength > 0 && c.respect > 0 && c.weight > 0
}

func (c *Creature) GetHoleLength() int {
	return c.holeLength
}

func (c *Creature) GetHealth() int {
	return c.health
}

func (c *Creature) GetRespect() int {
	return c.respect
}

func (c *Creature) GetWeight() int {
	return c.weight
}

func (c *Creature) SetHoleLength(_holeLength int) {
	c.holeLength = _holeLength
}

func (c *Creature) SetHealth(_health int) {
	c.health = _health
}

func (c *Creature) SetRespect(_respect int) {
	c.respect = _respect
}

func (c *Creature) SetWeight(_weight int) {
	c.weight = _weight
}

func (c *Creature) Dig(intensively bool) {

	if intensively {
		c.holeLength += 5
		c.health -= 30
	} else {
		c.holeLength += 2
		c.health -= 10
	}
}

func (c *Creature) Eat(green bool) {
	if green {
		if c.respect < 30 {
			c.health -= 30
		} else {
			c.health += 30
			c.weight += 30
		}
	} else {
		c.health += 10
		c.weight += 15
	}
}

func (c *Creature) FightWith(enemy *Creature) bool {
	result := false
	x := float32(enemy.weight) / float32(c.weight)
	if rand.Intn(c.weight+enemy.weight+1) <= c.weight {
		result = true
		c.respect += int(50 * x)
	}
	c.health -= int(50 * x)

	return result
}

func (c *Creature) Sleep() {
	c.holeLength -= 2
	c.health += 20
	c.respect -= 2
	c.weight -= 5
}
