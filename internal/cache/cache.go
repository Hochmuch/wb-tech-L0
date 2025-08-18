package cache

import (
	"github.com/google/uuid"
	"wb-tech-L0/internal/model"
)

type Cache struct {
	orders map[uuid.UUID]*model.Order
	queue  chan uuid.UUID
}

func New() *Cache {
	c := &Cache{}
	c.orders = make(map[uuid.UUID]*model.Order)
	c.queue = make(chan uuid.UUID, 100)
	return c
}

func (c *Cache) HasSpace() bool {
	return 100-len(c.orders) > 0
}

func (c *Cache) AddOrder(order *model.Order) {
	if c.HasSpace() {
		c.orders[order.OrderUID] = order
		c.queue <- order.OrderUID
	} else {
		toDelete := <-c.queue
		delete(c.orders, toDelete)

		c.orders[order.OrderUID] = order
		c.queue <- order.OrderUID
	}
}

func (c *Cache) GetOrder(orderUID uuid.UUID) (*model.Order, bool) {
	order, ok := c.orders[orderUID]
	return order, ok
}

func (c *Cache) GetCache(orders []*model.Order) {
	c.orders = make(map[uuid.UUID]*model.Order)
	c.queue = make(chan uuid.UUID, 100)

	for _, order := range orders {
		c.queue <- order.OrderUID
		c.orders[order.OrderUID] = order
	}
}
