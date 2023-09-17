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

type Processor interface {
	Process(*Request)
	SetNext(Processor)
}

type Kind int

const (
	Urgent Kind = 1 << iota
	Available
	Alert
)

type Request struct {
	Kind Kind
	Data string
}

type Printer struct {
	next Processor
}

func (p Printer) Process(req *Request) {
	if req.Kind&(Urgent|Alert) != 0 {
		log.Printf("Printer: %s", req.Data)
	}
	if p.next != nil {
		p.next.Process(req)
	}
}

func (p *Printer) SetNext(next Processor) {
	p.next = next
}

type Logger struct {
	next Processor
}

func (l *Logger) SetNext(next Processor) {
	l.next = next
}

func (l Logger) Process(req *Request) {
	if req.Kind&(Urgent|Available) != 0 {
		log.Printf("Logger: %s", req.Data)
	}
	if l.next != nil {
		l.next.Process(req)
	}
}

func DemoCoR() {
	l := &Logger{}
	p := &Printer{}
	l.SetNext(p)
	r1 := &Request{
		Urgent,
		"Urgent",
	}
	r2 := &Request{
		Available,
		"Available",
	}
	r3 := &Request{
		Alert,
		"Alert",
	}
	l.Process(r1)
	l.Process(r2)
	l.Process(r3)
}

// evictionAlgo — интерфейс strategy.
type evictionAlgo interface {
	evict(c *cache)
}

// реализация concreteStrategy

type fifo struct{}

func (l *fifo) evict(c *cache) {
	fmt.Println("Evicting by fifo strategy")
}

type lru struct{}

func (l *lru) evict(c *cache) {
	fmt.Println("Evicting by lru strategy")
}

type lfu struct{}

func (l *lfu) evict(c *cache) {
	fmt.Println("Evicting by lfu strategy")
}

// cache содержит контекст.
type cache struct {
	storage      map[string]string
	evictionAlgo evictionAlgo
	capacity     int
	maxCapacity  int
}

func initCache(e evictionAlgo) *cache {
	storage := make(map[string]string)
	return &cache{
		storage:      storage,
		evictionAlgo: e,
		capacity:     0,
		maxCapacity:  2,
	}
}

// setEvictionAlgo определяет алгоритм освобождения памяти.
func (c *cache) setEvictionAlgo(e evictionAlgo) {
	c.evictionAlgo = e
}

func (c *cache) add(key, value string) {
	if c.capacity == c.maxCapacity {
		c.evict()
	}
	c.capacity++
	c.storage[key] = value
}

func (c *cache) get(key string) {
	delete(c.storage, key)
}

func (c *cache) evict() {
	c.evictionAlgo.evict(c)
	c.capacity--
}

func DemoStrategy() {
	lfu := &lfu{}
	cache := initCache(lfu)
	cache.add("a", "1")
	cache.add("b", "2")
	cache.add("c", "3")
	lru := &lru{}
	cache.setEvictionAlgo(lru)
	cache.add("d", "4")
	fifo := &fifo{}
	cache.setEvictionAlgo(fifo)
	cache.add("e", "5")
}

func main() {
	border("Iterator")
	DemoIterator()
	border("Observer")
	DemoObserver()
	border("Command")
	DemoCommand()
	border("CoR")
	DemoCoR()
	border("Strategy")
	DemoStrategy()
}

func border(name string) {
	line := strings.Repeat("=", 80)
	out := fmt.Sprintf("%s\n\t\t\t%s\n%s\n", line, name, line)
	fmt.Println(out)
}
