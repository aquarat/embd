/*
 * Copyright (c) Karan Misra 2014
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and
 * associated documentation files (the "Software"), to deal in the Software without restriction,
 * including without limitation the rights to use, copy, modify, merge, publish, distribute,
 * sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or
 * substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT
 * NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
 * NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

// Generic Interrupt Pins.

package generic

import (
	"errors"
	"fmt"
	"sync"
	"syscall"

	"github.com/aquarat/embd"
)

const (
	MaxGPIOInterrupt = 64
)

var ErrorPinAlreadyRegistered = errors.New("pin interrupt already registered")

type interrupt struct {
	pin            embd.DigitalPin
	initialTrigger bool
	handler        func(embd.DigitalPin)
}

func (i *interrupt) Signal() {
	if !i.initialTrigger {
		i.initialTrigger = true
		return
	}
	i.handler(i.pin)
}

type epollListener struct {
	mu                sync.Mutex // Guards the following.
	fd                int
	interruptablePins map[int]*interrupt
}

var epollListenerInstance *epollListener

func getEpollListenerInstance() *epollListener {
	if epollListenerInstance == nil {
		epollListenerInstance = initEpollListener()
	}
	return epollListenerInstance
}

func initEpollListener() *epollListener {
	fd, err := syscall.EpollCreate1(0)
	if err != nil {
		panic(fmt.Sprintf("Unable to create epoll: %v", err))
	}
	listener := &epollListener{fd: fd, interruptablePins: make(map[int]*interrupt)}

	go func() {
		var epollEvents [MaxGPIOInterrupt]syscall.EpollEvent

		for {
			n, err := syscall.EpollWait(listener.fd, epollEvents[:], -1)
			if err != nil {
				panic(fmt.Sprintf("EpollWait error: %v", err))
			}
			listener.mu.Lock()
			for i := 0; i < n; i++ {
				if irq, ok := listener.interruptablePins[int(epollEvents[i].Fd)]; ok {
					irq.Signal()
				}
			}
			listener.mu.Unlock()
		}
	}()
	return listener
}

func registerInterrupt(pin *digitalPin, handler func(embd.DigitalPin)) error {
	l := getEpollListenerInstance()

	pinFd := int(pin.val.Fd())

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.interruptablePins[pinFd]; ok {
		return ErrorPinAlreadyRegistered
	}

	var event syscall.EpollEvent
	event.Events = syscall.EPOLLIN | (syscall.EPOLLET & 0xffffffff) | syscall.EPOLLPRI

	if err := syscall.SetNonblock(pinFd, true); err != nil {
		return err
	}

	event.Fd = int32(pinFd)

	if err := syscall.EpollCtl(l.fd, syscall.EPOLL_CTL_ADD, pinFd, &event); err != nil {
		return err
	}

	l.interruptablePins[pinFd] = &interrupt{pin: pin, handler: handler}

	return nil
}

func unregisterInterrupt(pin *digitalPin) error {
	l := getEpollListenerInstance()

	pinFd := int(pin.val.Fd())

	l.mu.Lock()
	defer l.mu.Unlock()

	if _, ok := l.interruptablePins[pinFd]; !ok {
		return nil
	}

	if err := syscall.EpollCtl(l.fd, syscall.EPOLL_CTL_DEL, pinFd, nil); err != nil {
		return err
	}

	if err := syscall.SetNonblock(pinFd, false); err != nil {
		return err
	}

	delete(l.interruptablePins, pinFd)
	return nil
}
