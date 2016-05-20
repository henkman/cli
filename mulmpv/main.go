package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"
)

/*
	#include <windows.h>
	#include <tlhelp32.h>
	static void kill_process(DWORD pid)
	{
		PROCESSENTRY32 pe = {0};
		pe.dwSize = sizeof(PROCESSENTRY32);
		HANDLE hSnap = CreateToolhelp32Snapshot(TH32CS_SNAPPROCESS, 0);
		if (!Process32First(hSnap, &pe)) {
			return;
		}
		BOOL bContinue = TRUE;
		while (bContinue) {
			if (pe.th32ParentProcessID == pid) {
				HANDLE hChildProc = OpenProcess(PROCESS_ALL_ACCESS,
					FALSE, pe.th32ProcessID);
				if (hChildProc) {
					TerminateProcess(hChildProc, 1);
					CloseHandle(hChildProc);
				}
			}
			bContinue = Process32Next(hSnap, &pe);
		}
		HANDLE hProc = OpenProcess(PROCESS_ALL_ACCESS, FALSE, pid);
		if (hProc) {
			TerminateProcess(hProc, 1);
			CloseHandle(hProc);
		}
	}
*/
import "C"

var (
	reValidFiles = regexp.MustCompile(
		"^.*\\.(?:wav|mp4|m4a|wmv|flv|webm|mkv|avi)$")

	_mpv      string
	_dir      string
	_start    uint
	_number   uint
	_vertical bool
)

func init() {
	flag.StringVar(&_mpv, "p", "mpv", "mpv executable")
	flag.StringVar(&_dir, "d", "", "dir")
	flag.UintVar(&_number, "n", 2, "number")
	flag.UintVar(&_start, "s", 0, "start")
	flag.BoolVar(&_vertical, "v", false, "vertical instead of horizontal")
	flag.Parse()
}

func DesktopSize() (int, int) {
	return int(C.GetSystemMetrics(C.SM_CXSCREEN)),
		int(C.GetSystemMetrics(C.SM_CYSCREEN))
}

func CalcPrimeFactors(n int) []int {
	x := n
	rv := []int{}
	ch := make(chan int)
	go func(max int, ch chan<- int) {
		ch <- 2
		for i := 3; i <= max; i += 2 {
			ch <- i
		}
		ch <- -1
	}(x, ch)
	for prime := <-ch; (prime != -1) && (x > 1); prime = <-ch {
	pl:
		for x%prime == 0 {
			x = x / prime
			for _, f := range rv {
				if prime == f {
					continue pl
				}
			}
			rv = append(rv, prime)
		}
		ch1 := make(chan int)
		go func(in <-chan int, out chan<- int, prime int) {
			for i := <-in; i != -1; i = <-in {
				if i%prime != 0 {
					out <- i
				}
			}
			out <- -1
		}(ch, ch1, prime)
		ch = ch1
	}
	return rv
}

func TileNaive(tw, th, n int, vertical bool) (int, int, int, int) {
	best := struct {
		w, h int
		r, c int
		diff int
	}{
		0, 0, 0, 0, 1e9,
	}
	ds := CalcPrimeFactors(n)
	if vertical {
		for _, r := range ds {
			c := (n / r)
			w := tw / c
			h := th / r
			diff := int(math.Abs(float64(w - h)))
			if diff < best.diff {
				best.diff = diff
				best.w = w
				best.h = h
				best.r = r
				best.c = c
			}
		}
	} else {
		for _, c := range ds {
			r := (n / c)
			w := tw / c
			h := th / r
			diff := int(math.Abs(float64(w - h)))
			if diff < best.diff {
				best.diff = diff
				best.w = w
				best.h = h
				best.r = r
				best.c = c
			}
		}
	}

	return best.w, best.h, best.r, best.c
}

func main() {
	if _dir == "" {
		flag.Usage()
		return
	}

	var w, h, rows, cols int
	{
		dw, dh := DesktopSize()
		w, h, rows, cols = TileNaive(dw, dh, int(_number), _vertical)
	}

	rand.Seed(time.Now().UnixNano())

	files := make([]os.FileInfo, 0, _number)
	{
		fd, err := os.Open(_dir)
		if err != nil {
			panic(err)
		}
		defer fd.Close()
		fis, err := fd.Readdir(-1)
		if err != nil {
			panic(err)
		}
		for _, fi := range fis {
			if !fi.IsDir() && reValidFiles.MatchString(fi.Name()) {
				files = append(files, fi)
			}
		}
		p := rand.Perm(len(files))
		for i, _ := range files {
			o := p[i]
			files[i], files[o] = files[o], files[i]
		}
		if uint(len(files)) > _number {
			files = files[:_number]
		}
	}
	cmds := make([]*exec.Cmd, len(files))
	{
		for y := 0; y < rows; y++ {
			for x := 0; x < cols; x++ {
				o := y*cols + x
				f := files[o]
				wx := x * w
				wy := y * h
				cmds[o] = exec.Command(_mpv,
					filepath.Join(_dir, f.Name()),
					fmt.Sprintf("--geometry=%dx%d+%d+%d", w, h, wx, wy),
					"--no-border",
					"--idle=yes",
				)
				if _start > 0 {
					cmds[o].Args = append(cmds[o].Args,
						fmt.Sprintf("--start=%d%%", _start))
				}
				cmds[o].Start()
			}
		}
	}
	{
		stdin := bufio.NewReader(os.Stdin)
		stdin.ReadString('\n')
	}
	for _, cmd := range cmds {
		C.kill_process(C.DWORD(cmd.Process.Pid))
	}
}
