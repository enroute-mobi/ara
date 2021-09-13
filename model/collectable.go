package model

import "time"

type Collectable struct {
	nextCollectAt time.Time
	collectedAt   time.Time
}

func (c *Collectable) NextCollectAt() time.Time {
	return c.nextCollectAt
}

func (c *Collectable) NextCollect(collectTime time.Time) {
	c.nextCollectAt = collectTime
}

func (c *Collectable) CollectedAt() time.Time {
	return c.collectedAt
}

func (c *Collectable) Updated(updateTime time.Time) {
	c.collectedAt = updateTime
}
