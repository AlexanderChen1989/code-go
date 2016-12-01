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

type Train struct {
	ob Observable

	source *Source
}

type Source struct {
	obs Observer
	emitFn func(Observer)
}

func (c *Source) Subscribe(obs Observer) {
	c.obs = obs
}

func (c *Source) Start() {
	c.emitFn(c.obs)
}

func From(fn func(Observer)) *Train {
	s := &Source{emitFn: fn}
	return &Train{source: s, ob: s}
}

func (t *Train) Map(mapper func(int) int) *Train {
	m := &MapOp{mapper: mapper}
	t.ob.Subscribe(m)
	t.ob = m
	return t
}

func (t *Train) Filter(filter func(int) bool) *Train {
	f := &FilterOp{filter: filter}
	t.ob.Subscribe(f) 
	t.ob = f
	return t
}

func (t *Train) Subscribe(obs Observer) {
	t.ob.Subscribe(obs)
	t.source.Start()
}


func main() {
	obs := From(func(ob Observer) {
		for i := 0; i < 20; i++ {
			ob.OnNext(i)
			time.Sleep(50 * time.Millisecond)
		}
		ob.OnComplete()
	}).Map(
		addOne,
	).Filter(
		odd,
	)

	obs.Subscribe(&IntObserver{})
	obs.Subscribe(&IntObserver{})
}

func addOne(v int) int { return v + 1}
func odd(v int) bool { return v%2 != 0}

