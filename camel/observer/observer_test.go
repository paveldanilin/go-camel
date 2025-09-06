package observer

import (
	"fmt"
	"testing"
)

type changeState struct {
	name     string
	oldValue any
	newValue any
}

type person struct {
	*Subject

	name string
	age  int
}

func (p *person) SetAge(newAge int) {
	state := &changeState{
		name:     "age",
		oldValue: p.age,
		newValue: newAge,
	}

	p.age = newAge
	p.Notify(state)
}

func (p *person) SerName(newName string) {
	state := &changeState{
		name:     "name",
		oldValue: p.name,
		newValue: newName,
	}

	p.name = newName
	p.Notify(state)
}

func TestObserverSubscribe(t *testing.T) {
	p := &person{
		Subject: NewSubject(),
		name:    "John",
		age:     10,
	}

	var updates []string

	p.Subscribe(func(state State) {
		if cs, isChangeState := state.(*changeState); isChangeState {
			fmt.Printf("[%s]: %v => %v\n", cs.name, cs.oldValue, cs.newValue)
			updates = append(updates, cs.name)
		}
	})

	p.SetAge(11)
	p.SerName("Adam")

	if len(updates) != 2 {
		t.Errorf("expected 2 updates, but got %v", len(updates))
	}

	if updates[0] != "age" {
		t.Errorf("expected update of 'age', but got '%s'", updates[0])
	}

	if updates[1] != "name" {
		t.Errorf("expected update of 'name', but got '%s'", updates[1])
	}
}

func TestObserverUnsubscribe(t *testing.T) {

}