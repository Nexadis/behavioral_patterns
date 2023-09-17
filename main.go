package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

func NewIter(max int) func() (int, bool) {
	n := 0
	return func() (int, bool) {
		if n >= max {
			return 0, false
		}
		n++
		return n - 1, true
	}
}

func DemoIterator() {
	iter := NewIter(5)
	for {
		n, ok := iter()
		if !ok {
			break
		}
		log.Println(n)
	}
}

type Subscriber struct {
	ID int
}

// Subscribe ожидает уведомления.
func (s Subscriber) Subscribe(c *sync.Cond) {
	for {
		c.L.Lock()
		c.Wait()
		fmt.Printf("Subscriber %v is notified\n", s.ID)
		c.L.Unlock()
	}
}

func DemoObserver() {
	cond := sync.NewCond(new(sync.Mutex))
	s1 := Subscriber{1}
	go s1.Subscribe(cond)
	s2 := Subscriber{2}
	go s2.Subscribe(cond)
	s3 := Subscriber{3}
	go s3.Subscribe(cond)
	time.Sleep(200 * time.Millisecond)
	// отправка уведомлений всем подписчикам
	cond.Broadcast()
	time.Sleep(200 * time.Millisecond)
	fmt.Println("Once more")
	cond.Broadcast()
	time.Sleep(200 * time.Millisecond)
}

type command interface {
	execute()
}

type receiver interface {
	action()
}

type invoker struct {
	commands map[string]command
}

func newInvoker() *invoker {
	i := new(invoker)
	i.commands = make(map[string]command)
	return i
}

func (i *invoker) do(c string) {
	log.Printf("Exec %s command\n", c)
	i.commands[c].execute()
}

type printer struct {
	receiver receiver
}

func (c *printer) execute() {
	c.receiver.action()
}

type rcvr struct {
	name string
}

func (r *rcvr) action() {
	fmt.Println(r.name)
}

func DemoCommand() {
	h := rcvr{"Hello"}
	hello := printer{&h}
	b := rcvr{"Bye Bye"}
	bye := printer{&b}
	i := newInvoker()
	i.commands["hello_command"] = &hello
	i.commands["bye_command"] = &bye
	i.do("hello_command")
	i.do("bye_command")
}

func main() {
	border("Iterator")
	DemoIterator()
	border("Observer")
	DemoObserver()
	border("Command")
	DemoCommand()
}

func border(name string) {
	line := strings.Repeat("=", 80)
	out := fmt.Sprintf("%s\n\t\t\t%s\n%s\n", line, name, line)
	fmt.Println(out)
}
