package util

type CooldownTimer struct {
	Cooldown   int
	ticker     int
	InCooldown bool
}

func NewCooldownTimer(cooldown int) *CooldownTimer {
	return &CooldownTimer{Cooldown: cooldown}
}

func (c *CooldownTimer) IncrementTicker() {
	if c.InCooldown {
		c.ticker++
	}
}

func (c *CooldownTimer) CheckCooldown() {
	if c.ticker >= c.Cooldown {
		c.ticker = 0
		c.InCooldown = false
	}
}

func (c *CooldownTimer) StartCooldown() {
	c.InCooldown = true
}

func (c *CooldownTimer) GetDisplay() int {
	var cd int = 0
	if c.InCooldown {
		cd = max(c.Cooldown-c.ticker, 0)
	}
	return cd
}
