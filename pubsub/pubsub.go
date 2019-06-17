package pubsub

import (
	"errors"
	"fmt"
	"time"
)

type Message interface{}

type Sub struct {
	ReadChan chan Message
	input    chan Message
	quit     chan int
	done     chan int
}

func (s Sub) readRoutine() {
	for {
		select {
		case v := <-s.input:
			if v == nil {
				fmt.Println("Closing ReadChan as input closed")
				close(s.ReadChan)
				return
			}

			select {
			case s.ReadChan <- v:
			case <-time.After(10 * time.Second):
				// failed to write to the reading entity
				// Its a buffered channel so this should not happen
				// Close this subscriber
				fmt.Println("Closing ReadChan as ReadChan is not writeable")
				close(s.ReadChan)
				return
			}
		case <-s.quit:
			fmt.Println("Closing ReadChan as quit called")
			close(s.ReadChan)
			s.done <- 1
			return
		}
	}
}

func (s Sub) Quit() {
	s.quit <- 1
	<-s.done
}

type Pub struct {
	WriteChan chan Message
	output    chan Message
	quit      chan int
	done      chan int
	// Outgoing Quit Signal
	QuitChan chan int
	DoneChan chan int
}

func (p Pub) Quit() {
	p.quit <- 1
	<-p.done
}

func (p Pub) writeRoutine() {
	for {
		select {
		case v := <-p.WriteChan:
			if v == nil {
				return
			}

			select {
			case p.output <- v:
			default:
				fmt.Println("Closing writeRoutine as output closed")
				return
			}
		case <-p.quit:
			p.done <- 1
			fmt.Println("Closing writeRoutine as quit called")
			return
		}
	}
}

type PubSub struct {
	In   []Pub
	Out  []Sub
	pipe chan Message
	quit chan int
	done chan int
}

func NewPubSub() *PubSub {
	return &PubSub{In: []Pub{}, Out: []Sub{}, pipe: make(chan Message, 100),
		quit: make(chan int, 10),
		done: make(chan int, 10),
	}
}

func (ps *PubSub) Subscribe() Sub {
	sub := Sub{}
	sub.done = make(chan int, 10)
	sub.quit = make(chan int, 10)
	sub.ReadChan = make(chan Message, 10)
	sub.input = make(chan Message, 10)
	ps.Out = append(ps.Out, sub)
	go sub.readRoutine()
	return sub
}

func (ps *PubSub) Publish() Pub {
	pub := Pub{}
	pub.done = make(chan int, 10)
	pub.quit = make(chan int, 10)
	pub.QuitChan = make(chan int, 10)
	pub.DoneChan = make(chan int, 10)

	pub.WriteChan = make(chan Message, 10)
	pub.output = ps.pipe
	ps.In = append(ps.In, pub)
	go pub.writeRoutine()
	return pub
}

func (ps *PubSub) quitSubscribers() {
	for _, s := range ps.Out {
		if s.input != nil {
			close(s.input)
		}
	}
}

func (ps *PubSub) quitPublishers() {
	for _, p := range ps.In {
		close(p.QuitChan)
	}
}

func (ps *PubSub) writeToSubscribers(m Message) {
	for _, s := range ps.Out {
		select {
		case s.input <- m:
		default:
			close(s.input)
			s.input = nil
		}
	}

}

func (ps *PubSub) PipeRoutine() {
	for {
		select {
		case v := <-ps.pipe:
			if v == nil {
				ps.quitSubscribers()
				ps.quitPublishers()
				return
			}

			ps.writeToSubscribers(v)
		case <-ps.quit:
			fmt.Println("GotQuit")
			ps.quitSubscribers()
			fmt.Println("QuitSubscribers")
			ps.quitPublishers()
			fmt.Println("QuitPublishers")
			ps.done <- 1
			return
		}
	}

}

func (ps *PubSub) Quit() {
	ps.quit <- 1
	<-ps.done
}

func (ps *PubSub) Send(m Message) error {
	select {
	case ps.pipe <- m:
		return nil
	case <-time.After(10 * time.Second):
		return errors.New("Timeout")
	}
}
