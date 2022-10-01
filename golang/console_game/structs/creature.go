package structs

type Creature struct {
	holeLength int
	health     int
	respect    int
	weight     int
}

func (c *Creature) dig(intensively bool) {

	if intensively {
		c.holeLength += 5
		c.health -= 30
	} else {
		c.holeLength += 2
		c.health -= 10
	}
}

func (c *Creature) eat(green bool) {
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

func (c *Creature) fightWith(enemy Creature) {

}
