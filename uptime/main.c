/*
	uptime for windows
	build with:
	gcc -Wall -std=c99 -pedantic -o uptime main.c -O3 -s -nostdlib -fno-asynchronous-unwind-tables -fno-ident -ffunction-sections -static -Wl,-e,_main -lkernel32 -ladvapi32
*/
#define WIN32_LEAN_AND_MEAN
#include <windows.h>

#define console_print_const(c, msg) \
	console_write(c, msg, sizeof(msg))

typedef struct {
	HANDLE stdout;
} Console;

static void console_init(Console *c)
{
	AllocConsole();
	c->stdout = GetStdHandle(STD_OUTPUT_HANDLE);
}

static void console_deinit(Console *c)
{
	FreeConsole();
}

static unsigned long console_write(Console *c, void *msg, size_t len)
{
	unsigned long n;
	WriteConsoleA(c->stdout, msg, len, &n, NULL);
	return n;
}

static unsigned long console_put(Console *c, char b)
{
	return console_write(c, (char *)&b, 1);
}

static unsigned long console_printu(Console *c, unsigned v)
{
	char *tp, *rtp;
	char t[11];
	size_t n;
	char x;
	if(v < 10) {
		return console_put(c, v+'0');
	}
	tp = rtp = t;
	while(v) {
		*tp++ = (v%10)+'0';
		v = v/10;
	}
	n = tp-t;
	tp--;
	while(tp > rtp) {
		x = *tp;
		*tp = *rtp;
		*rtp = x;
		tp--;
		rtp++;
	}
	return console_write(c, t, n);
}

static void UnixTimeToFileTime(unsigned long t, FILETIME *ft)
{
	LONGLONG ll;
	ll = Int32x32To64(t, 10000000) + 116444736000000000;
	ft->dwLowDateTime = (DWORD)ll;
	ft->dwHighDateTime = ll >> 32;
}

static inline void printu_prefixpadzero(Console *c, unsigned v)
{
	if(v < 10) {
		console_put(c, '0');
	}
	console_printu(c, v);
}

static void print_systemtime(Console *c, SYSTEMTIME *st, TIME_ZONE_INFORMATION *tz)
{
	SYSTEMTIME sttz;
	SystemTimeToTzSpecificLocalTime(tz, st, &sttz);
	printu_prefixpadzero(c, sttz.wDay);
	console_put(c, '.');
	printu_prefixpadzero(c, sttz.wMonth);
	console_put(c, '.');
	console_printu(c, sttz.wYear);
	console_put(c, ' ');
	printu_prefixpadzero(c, sttz.wHour);
	console_put(c, ':');
	printu_prefixpadzero(c, sttz.wMinute);
	console_put(c, ':');
	printu_prefixpadzero(c, sttz.wSecond);
}

static void print_timestamp(Console *c, unsigned ts, TIME_ZONE_INFORMATION *tz)
{
	FILETIME ft;
	SYSTEMTIME st;
	UnixTimeToFileTime(ts, &ft);
	FileTimeToSystemTime(&ft, &st);
	print_systemtime(c, &st, tz);
}

#define EVENT_ID_SHUTDOWN 6006
#define EVENT_ID_STARTUP 6009

#define RECORDS_SIZE (8*1024)

static void uptime_print(Console *c)
{
	char *records;
	EVENTLOGRECORD *record;
	SYSTEMTIME time;
	TIME_ZONE_INFORMATION tz;
	HANDLE eventlog;
	unsigned long read, req;
	unsigned o, id, ts, lts;
	unsigned start, findstart;
	
	eventlog = OpenEventLog("", "system");
	if(!eventlog) {
		console_print_const(c, "could not open event log\n");
		return;
	}
	GetTimeZoneInformation(&tz);
	records = HeapAlloc(GetProcessHeap(), 0, RECORDS_SIZE);
	findstart = 1;
	while(ReadEventLog(eventlog,
		EVENTLOG_SEQUENTIAL_READ|EVENTLOG_FORWARDS_READ,
		0, records, RECORDS_SIZE, &read, &req)
	) {
		for(o=0; o<read; o+=record->Length) {
			record = (EVENTLOGRECORD *)&records[o];
			id = record->EventID & 0xFFFF;
			ts = record->TimeGenerated;
			if(findstart) {
				if(id == EVENT_ID_STARTUP) {
					start = ts;
					findstart = 0;
				}
			} else {
				if(id == EVENT_ID_STARTUP) {
					print_timestamp(c, start, &tz);
					console_put(c, ';');
					print_timestamp(c, lts, &tz);
					console_put(c, '\n');
					findstart = 1;
				} else if(id == EVENT_ID_SHUTDOWN) {
					print_timestamp(c, start, &tz);
					console_put(c, ';');
					print_timestamp(c, ts, &tz);
					console_put(c, '\n');
					findstart = 1;
				}
			}
			lts = ts;
		}
	}
	if(!findstart) {
		GetSystemTime(&time);
		print_timestamp(c, start, &tz);
		console_put(c, ';');
		print_systemtime(c, &time, &tz);
		console_put(c, '\n');
	}
	HeapFree(GetProcessHeap(), 0, records);
	CloseEventLog(eventlog);
}

void _main(void) __asm__("_main");
void _main(void)
{
	Console c;
	
	console_init(&c);
	uptime_print(&c);
	console_deinit(&c);
	ExitProcess(0);
}