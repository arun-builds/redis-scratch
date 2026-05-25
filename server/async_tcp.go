package server

import (
	"log"
	"net"
	"syscall"
	"time"

	"github.com/arun-builds/redis-scratch/config"
	"github.com/arun-builds/redis-scratch/core"
)

var con_clients int = 0
var cronFrequency time.Duration = 1 * time.Second
var lastCronExecTime time.Time = time.Now()

func RunAsyncTCPServer() error {
	log.Println("starting an asynchronous TCP server on", config.Host, config.Port)

	max_clients := 20000

	var events []syscall.Kevent_t = make([]syscall.Kevent_t, max_clients)

	serverFD, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}

	defer syscall.Close(serverFD)

	// Set socket to non-blocking mode
	if err = syscall.SetNonblock(serverFD, true); err != nil {
		return err
	}

	// Bind IP and port
	ip4 := net.ParseIP(config.Host).To4()
	if err = syscall.Bind(serverFD, &syscall.SockaddrInet4{
		Port: config.Port,
		Addr: [4]byte{ip4[0], ip4[1], ip4[2], ip4[3]},
	}); err != nil {
		return err
	}

	// Start listening
	if err = syscall.Listen(serverFD, max_clients); err != nil {
		return err
	}

	// AsyncIO

	// Create KQUEUE instance

	kqueueFD, err := syscall.Kqueue()
	if err != nil {
		log.Fatal(err)
	}
	defer syscall.Close(kqueueFD)

	// Specify events we want hints about
	var socketServerEvent syscall.Kevent_t = syscall.Kevent_t{
		Ident:  uint64(serverFD),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD,
	}

	// Listen to read events on server socket
	_, err = syscall.Kevent(
		kqueueFD,
		[]syscall.Kevent_t{socketServerEvent},
		nil,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		nevents, err := syscall.Kevent(
			kqueueFD,
			nil,
			events,
			nil,
		)
		if err != nil {
			continue
		}

		for i := 0; i < nevents; i++ {

			if time.Now().After(lastCronExecTime.Add(cronFrequency)) {
				core.DeleteExpiredKeys()
				lastCronExecTime = time.Now()
			}

			fd := int(events[i].Ident)

			// Server socket ready
			if fd == serverFD {

				// Accept incoming client
				clientFD, _, err := syscall.Accept(serverFD)
				if err != nil {
					log.Println("err", err)
					continue
				}

				// Increase concurrent clients count
				con_clients++
				syscall.SetNonblock(clientFD, true)

				// Add new client socket to kqueue monitoring
				var socketClientEvent syscall.Kevent_t = syscall.Kevent_t{
					Ident:  uint64(clientFD),
					Filter: syscall.EVFILT_READ,
					Flags:  syscall.EV_ADD,
				}

				_, err = syscall.Kevent(
					kqueueFD,
					[]syscall.Kevent_t{socketClientEvent},
					nil,
					nil,
				)
				if err != nil {
					log.Fatal(err)
				}

			} else {

				comm := core.FDComm{
					Fd: fd,
				}

				cmds, err := readCommands(comm)
				if err != nil {
					syscall.Close(fd)
					con_clients -= 1
					continue
				}

				respond(cmds, comm)
			}
		}
	}
}
