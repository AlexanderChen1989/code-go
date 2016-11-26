package main

import (
	"fmt"
	"time"
)

/*
	1,2,3,4 -> 5,6,7,8 -> 6,8
*/

type Observer interface {
	OnNext(int)
	OnComplete()
	OnError(error)
}

type Observable interface {
	Subscribe(Observer)
}

func Map(ob Observable, mapper func(int) int) Observable {
	m := &MapOp{
		mapper: mapper,
	}
	ob.Subscribe(m)
	return m
}

func Filter(ob Observable, filter func(int) bool) Observable {
	f := &FilterOp{
		filter: filter,
	}
	ob.Subscribe(f)
	return f
}

type MapOp struct {
	obs    Observer
	mapper func(int) int
}

func (m *MapOp) OnNext(next int) {
	m.obs.OnNext(m.mapper(next))
}

func (m *MapOp) OnComplete() {
	m.obs.OnComplete()
}

func (m *MapOp) OnError(err error) {
	m.obs.OnError(err)
}

func (m *MapOp) Subscribe(obs Observer) {
	m.obs = obs
}

type FilterOp struct {
	obs    Observer
	filter func(int) bool
}

func (m *FilterOp) OnNext(next int) {
	if m.filter(next) {
		m.obs.OnNext(next)
	}
}

func (m *FilterOp) OnComplete() {
	m.obs.OnComplete()
}

func (m *FilterOp) OnError(err error) {
	m.obs.OnError(err)
}

func (m *FilterOp) Subscribe(obs Observer) {
	m.obs = obs
}

type Counter struct {
	obs Observer
}

func (c *Counter) Subscribe(obs Observer) {
	c.obs = obs
}

func (c *Counter) Start() {
	for i := 0; i < 10; i++ {
		c.obs.OnNext(i)
		time.Sleep(time.Second)
	}
	c.obs.OnComplete()
}

type IntObserver struct{}

func (obs *IntObserver) OnNext(v int) {
	fmt.Println(v)
}
func (obs *IntObserver) OnComplete() {
	fmt.Println("complete")
}
func (obs *IntObserver) OnError(err error) {
	fmt.Println("error", err)
}

type Combinator struct {
	ob Observable
}

func From(ob Observable) *Combinator {
	return &Combinator{ob: ob}
}

func (c *Combinator) Map(mapper func(int) int) *Combinator {
	c.ob = Map(c.ob, mapper)
	return c
}

func (c *Combinator) Filter(filter func(int) bool) *Combinator {
	c.ob = Filter(c.ob, filter)
	return c
}

func (c *Combinator) Subscribe(obs Observer) {
	c.ob.Subscribe(obs)
}

func main() {
	c := &Counter{}

	From(c).Map(func(v int) int {
		return v + 10
	}).Filter(func(v int) bool {
		return v%2 == 0
	}).Subscribe(&IntObserver{})

	c.Start()
}
