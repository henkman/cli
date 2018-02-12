package main

import (
	"bufio"
	"flag"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"time"

	ole "github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

/*
	#include <windows.h>
	struct Windows {
		UINT count;
		HWND window[4];
	};
	struct EnumData {
		DWORD processId;
		struct Windows windows;
	};
	BOOL CALLBACK EnumProc(HWND hWnd, LPARAM lParam) {
		if (!IsWindowVisible(hWnd)) return TRUE;
		struct EnumData *ed = (struct EnumData*)lParam;
		struct Windows *wins = (struct Windows*)&ed->windows;
		DWORD processId;
		GetWindowThreadProcessId(hWnd, &processId);
		if (ed->processId == processId) {
			wins->window[wins->count++] = hWnd;
		}
		return TRUE;
	}
	struct Windows FindWindowsOfProcess(DWORD processId) {
		struct EnumData ed = {
			.processId = processId,
			.windows = {
				0, {0}
			},
		};
		EnumWindows(EnumProc, (LPARAM)&ed);
		return ed.windows;
	}
*/
import "C"

var reValidFiles = regexp.MustCompile(
	"^.*\\.(?:wav|mp4|m4a|wmv|flv|webm|mkv|avi)$")

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
	var opts struct {
		Mpc_hc   string
		Dir      string
		Number   uint
		Vertical bool
	}
	flag.StringVar(&opts.Mpc_hc, "p",
		"C:\\Program Files (x86)\\MPC-HC\\mpc-hc.exe",
		"mpc-hc executable")
	flag.StringVar(&opts.Dir, "d", "", "dir")
	flag.UintVar(&opts.Number, "n", 2, "number")
	flag.BoolVar(&opts.Vertical, "v", false, "vertical instead of horizontal")
	flag.Parse()

	if opts.Dir == "" {
		flag.Usage()
		return
	}

	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	var w, h, rows, cols int
	{
		dw, dh := DesktopSize()
		w, h, rows, cols = TileNaive(dw, dh, int(opts.Number), opts.Vertical)
	}

	var wscript *ole.IDispatch
	{
		unknown, _ := oleutil.CreateObject("WScript.Shell")
		defer unknown.Release()
		wscript, _ = unknown.QueryInterface(ole.IID_IDispatch)
	}
	defer wscript.Release()

	rand.Seed(time.Now().UnixNano())

	files := make([]os.FileInfo, 0, opts.Number)
	{
		fd, err := os.Open(opts.Dir)
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
		if uint(len(files)) > opts.Number {
			files = files[:opts.Number]
		}
	}
	cmds := make([]*exec.Cmd, len(files))
	{
		for i, file := range files {
			cmds[i] = exec.Command(opts.Mpc_hc, filepath.Join(opts.Dir, file.Name()))
			cmds[i].Start()
		}
		tper := time.Duration(750)
		time.Sleep(time.Millisecond * tper * time.Duration(len(cmds)))
	}
	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			cmd := cmds[y*cols+x]
			oleutil.CallMethod(wscript, "AppActivate", cmd.Process.Pid)
			time.Sleep(time.Millisecond * time.Duration(100))
			oleutil.CallMethod(wscript, "SendKeys", "2", 0)
			time.Sleep(time.Millisecond * time.Duration(100))
			{
				wx := x * w
				wy := y * h
				wins := C.FindWindowsOfProcess(C.DWORD(cmd.Process.Pid))
				if wins.count == 0 {
					continue
				}
				hwnd := wins.window[0]
				C.MoveWindow(hwnd, C.int(wx), C.int(wy),
					C.int(w), C.int(h),
					C.WINBOOL(1))
			}
		}
	}
	{
		stdin := bufio.NewReader(os.Stdin)
		stdin.ReadString('\n')
	}
	for _, cmd := range cmds {
		cmd.Process.Kill()
	}
}
