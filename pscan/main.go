/*
	TO IMPLEMENT:
	- shuffle, requires that all ips are know at the start
	- only send syn packet(http://www.darkcoding.net/uncategorized/raw-sockets-in-go-ip-layer/)
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type Target interface {
	HasNext() bool
	Next() string
	Count() int
}

type Ips struct {
	IPs     []net.IP
	Ports   []uint16
	CurPort int
	CurIP   int
}

func (i *Ips) HasNext() bool {
	return i.CurPort < len(i.Ports)
}
func (i *Ips) Next() string {
	defer func() {
		i.CurIP++
		if i.CurIP >= len(i.IPs) {
			i.CurIP = 0
			i.CurPort++
		}
	}()
	return fmt.Sprintf("%s:%d", i.IPs[i.CurIP].String(), i.Ports[i.CurPort])
}
func (i *Ips) Count() int {
	return len(i.IPs) * len(i.Ports)
}

type IpMask struct {
	Start   net.IP
	IP      net.IP
	IPNet   *net.IPNet
	Ports   []uint16
	IPCount int
	CurPort int
}

func (i *IpMask) HasNext() bool {
	return i.CurPort < len(i.Ports)
}
func (i *IpMask) Next() string {
	defer func() {
		incIp(i.IP)
		if !i.IPNet.Contains(i.IP) {
			i.IP = i.Start
			i.CurPort++
		}
	}()
	return fmt.Sprintf("%s:%d", i.IP.String(), i.Ports[i.CurPort])
}
func (i *IpMask) Count() int {
	return i.IPCount * len(i.Ports)
}

type IpRange struct {
	Start    net.IP
	End      net.IP
	IP       net.IP
	Ports    []uint16
	CurPort  int
	IPCount  int
	Overflow bool
}

func (i *IpRange) HasNext() bool {
	return i.CurPort < len(i.Ports)
}
func (i *IpRange) Next() string {
	defer func() {
		incIp(i.IP)
		if i.Overflow {
			i.Overflow = false
			i.IP = i.Start
			i.CurPort++
		}
		if i.IP.Equal(i.End) {
			i.Overflow = true
		}
	}()
	return fmt.Sprintf("%s:%d", i.IP.String(), i.Ports[i.CurPort])
}
func (i *IpRange) Count() int {
	return i.IPCount * len(i.Ports)
}

type Scan struct {
	target  Target
	workers uint
	timeout time.Duration
}

func NewScan(target Target, workers uint, timeout time.Duration) *Scan {
	s := new(Scan)
	s.target = target
	s.workers = workers
	s.timeout = timeout
	return s
}

func (s *Scan) check(hosts <-chan string, results chan<- string) {
	for host := range hosts {
		c, err := net.DialTimeout("tcp", host, s.timeout)
		if err == nil {
			c.Close()
			results <- host
		} else {
			results <- ""
		}
	}
}

func (s *Scan) Run() {
	results := make(chan string, 100)
	hosts := make(chan string, 100)

	fmt.Println("scanning", s.target.Count(), "targets")

	var i uint
	for i = 0; i < s.workers; i++ {
		go s.check(hosts, results)
	}

	go func() {
		for s.target.HasNext() {
			hosts <- s.target.Next()
		}
		close(hosts)
	}()

	for i := 0; i < s.target.Count(); i++ {
		host := <-results
		if host != "" {
			fmt.Println(host)
		}
	}
}

func incIp(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func IPsInNet(sip net.IP, ipnet *net.IPNet) int {
	var r int
	for ip := sip.Mask(ipnet.Mask); ipnet.Contains(ip); incIp(ip) {
		r++
	}
	return r
}

func IPsInRange(s, e net.IP) int {
	i := s
	var r int
	for !i.Equal(e) {
		incIp(i)
		r++
	}
	return r + 1
}

func parseTarget(sports, shosts string) Target {
	var ports []uint16
	if strings.Contains(sports, "-") {
		hs := strings.Split(sports, "-")
		if len(hs) != 2 {
			log.Fatal("")
		}
		b, err := strconv.Atoi(hs[0])
		if err != nil {
			log.Fatal(err)
		}
		e, err := strconv.Atoi(hs[1])
		if err != nil {
			log.Fatal(err)
		}
		if b > e {
			log.Fatal("range invalid")
		}
		ports = make([]uint16, e-b+1)
		for i := b; i <= e; i++ {
			ports[i-b] = uint16(i)
		}
	} else {
		ps := strings.Split(sports, ",")
		ports = make([]uint16, len(ps))
		for i, s := range ps {
			p, err := strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
			ports[i] = uint16(p)
		}
	}

	if strings.Contains(shosts, "/") {
		ip, ipnet, err := net.ParseCIDR(shosts)
		if err != nil {
			log.Fatal(err)
		}
		c := ip
		i := ip
		return &IpMask{ip.Mask(ipnet.Mask),
			c.Mask(ipnet.Mask),
			ipnet,
			ports,
			IPsInNet(i, ipnet),
			0}
	}

	if strings.Contains(shosts, "-") {
		hs := strings.Split(shosts, "-")
		if len(hs) != 2 {
			log.Fatal("")
		}
		i := net.ParseIP(hs[0])
		c := net.ParseIP(hs[0])
		s := net.ParseIP(hs[0])
		e := net.ParseIP(hs[1])
		return &IpRange{s, e, c, ports, 0, IPsInRange(i, e), false}
	}

	hs := strings.Split(shosts, ",")
	ips := make([]net.IP, len(hs))
	for i, s := range hs {
		ips[i] = net.ParseIP(s)
	}
	return &Ips{ips, ports, 0, 0}
}

func main() {
	var opts struct {
		Ports   string
		Hosts   string
		Workers uint
		Timeout uint
	}
	flag.StringVar(&opts.Ports, "p", "", "Port(s) comma separated or range with port-port")
	flag.StringVar(&opts.Hosts, "h", "",
		"Host(s), either comma separated, a masked net or range with ip-ip")
	flag.UintVar(&opts.Workers, "w", 200, "Number of workers, default is 200")
	flag.UintVar(&opts.Timeout, "t", 2, "Timeout in seconds, default is 2")
	flag.Parse()

	if opts.Hosts == "" || opts.Ports == "" {
		flag.Usage()
		return
	}

	target := parseTarget(opts.Ports, opts.Hosts)
	s := NewScan(target, opts.Workers, time.Second*time.Duration(opts.Timeout))
	s.Run()
}
