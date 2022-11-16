package creature

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDig(t *testing.T) {

	t.Run("intensively_hole_length", func(t *testing.T) {
		c := New(100, 100, 100, 100)
		c.Dig(true)
		got := c.GetHoleLength()
		exp := 105
		assert.Equal(t, exp, got)
	})

	t.Run("intensively_health", func(t *testing.T) {
		c := New(100, 100, 100, 100)
		c.Dig(true)
		got := c.GetHealth()
		exp := 70
		assert.Equal(t, exp, got)

	})

	t.Run("basic_hole_length", func(t *testing.T) {
		c := New(100, 100, 100, 100)
		c.Dig(false)
		got := c.GetHoleLength()
		exp := 102
		assert.Equal(t, exp, got)
	})

	t.Run("basic_health", func(t *testing.T) {
		c := New(100, 100, 100, 100)
		c.Dig(false)
		got := c.GetHealth()
		exp := 90
		assert.Equal(t, exp, got)
	})
}

func TestEat(t *testing.T) {

	t.Run("health_green_respect_less_30", func(t *testing.T) {

		c := New(100, 100, 20, 100)
		c.Eat(true)
		got := c.GetHealth()
		exp := 70
		assert.Equal(t, exp, got)
	})

	t.Run("health_green_respect_more_30", func(t *testing.T) {
		c := New(100, 100, 100, 100)
		c.Eat(true)
		got := c.GetHealth()
		exp := 130
		assert.Equal(t, exp, got)
	})

	t.Run("weight_green_respect_more_30", func(t *testing.T) {
		c := New(100, 100, 100, 100)
		c.Eat(true)
		got := c.GetWeight()
		exp := 130
		assert.Equal(t, exp, got)
	})

	t.Run("health_withered", func(t *testing.T) {
		c := New(100, 100, 100, 100)
		c.Eat(false)
		got := c.GetHealth()
		exp := 110
		assert.Equal(t, exp, got)
	})

	t.Run("weight_withered", func(t *testing.T) {
		c := New(100, 100, 100, 100)
		c.Eat(false)
		got := c.GetWeight()
		exp := 115
		assert.Equal(t, exp, got)
	})
}

func TestFightWith(t *testing.T) {

	t.Run("respect", func(t *testing.T) {
		c := New(100, 100, 100, 50)
		enemy := New(100, 100, 100, 50)

		res := c.FightWith(enemy)
		got := c.GetRespect()

		exp := 100
		if res {
			exp = 150
		}
		assert.Equal(t, exp, got)
	})

	t.Run("health", func(t *testing.T) {
		c := New(100, 100, 100, 50)
		enemy := New(100, 100, 100, 100)

		c.FightWith(enemy)
		got := c.GetHealth()

		exp := 0
		assert.Equal(t, exp, got)
	})

}
