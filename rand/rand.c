#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include <wincrypt.h>

#define ARR_LEN(x) \
	(sizeof(x)/sizeof(x[0]))

typedef enum {
	ArgProg,
	ArgBool,
	ArgValue
} ArgType;

typedef struct {
	ArgType type;
	char *key;
	union {
		struct {
			char *value;
			int len;
		};
		int isset;
	};
} Arg;

typedef struct {
	HANDLE stdout;
} Console;

static void args_make(Arg *arg, ArgType type, char *key)
{
	arg->type = type;
	arg->key = key;
	if(type == ArgBool) {
		arg->isset = 0;
	} else if(type == ArgValue || type == ArgProg) {
		arg->value = NULL;
		arg->len = 0;
	}
}

static void args_get(Arg **args, int count, const char *cmdline)
{
	static char cmd[8*1024];
	int i, o, n, ignorespace;
	Arg *key;

	lstrcpy(cmd, cmdline);

	/* parse program name */
	o = 0;
	while(cmd[o] && cmd[o] != ' ') {
		o++;
	}
	cmd[o] = 0;
	for(i=0; i<count; i++) {
		if(args[i]->type == ArgProg) {
			args[i]->value = &cmd[0];
			args[i]->len = o;
			break;
		}
	}
	o++;
	if(!cmd[o]) {
		return;
	}
	while(cmd[o] == ' ' || cmd[o] == '\t') {
		o++;
	}

	/* begin key value pairs */
	key = NULL;
	n = o;
	ignorespace = 0;
	while(cmd[o]) {
		if(cmd[o] == '"') {
			ignorespace = !ignorespace;
			o++;
		} else if(cmd[o] == ' ' && !ignorespace) {
			cmd[o] = 0;
			if(key) {
				key->value = &cmd[n];
				key->len = o-n;
				if(key->value[0] == '"') {
					key->value++;
					key->len--;
				}
				if(key->value[key->len-1] == '"') {
					key->value[--key->len] = 0;
				}
				key = NULL;
			} else {
				for(i=0; i<count; i++) {
					if(lstrcmp(args[i]->key, &cmd[n]) == 0) {
						if(args[i]->type == ArgBool) {
							args[i]->isset = 1;
						} else {
							key = args[i];
						}
						break;
					}
				}
			}
			o++;
			if(!cmd[o]) {
				break;
			}
			while(cmd[o] == ' ' || cmd[o] == '\t') {
				o++;
			}
			n = o;
		} else {
			o++;
		}
	}

	/* the last value, if any */
	if(key) {
		key->value = &cmd[n];
		key->len = o-n;
		if(key->value[0] == '"') {
			key->value++;
			key->len--;
		}
		if(key->value[key->len-1] == '"') {
			key->value[--key->len] = 0;
		}
	} else {
		for(i=0; i<count; i++) {
			if(lstrcmp(args[i]->key, &cmd[n]) == 0) {
				if(args[i]->type == ArgBool) {
					args[i]->isset = 1;
				}
				break;
			}
		}
	}
}

static void console_init(Console *c)
{
	AllocConsole();
	c->stdout = GetStdHandle(STD_OUTPUT_HANDLE);
}

static void console_deinit(Console *c)
{
	FreeConsole();
}

#define console_print_const(c, msg) \
	console_write(c, msg, sizeof(msg))

static unsigned long console_write(Console *c, const void *msg, size_t len)
{
	unsigned long n;
	WriteConsoleA(c->stdout, msg, len, &n, NULL);
	return n;
}

static unsigned long console_put(Console *c, const char b)
{
	char msg[1] = {b};
	return console_write(c, msg, sizeof(msg));
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

static unsigned stou(const char *s)
{
	char *p = (char *) s;
	unsigned r = 0;
	while(*p) {
		if(*p < '0' || *p > '9') {
			return 0;
		}
		r *= 10;
		r += *p - '0';
		p++;
	}
	return r;
}

void _main(void) asm("_main");
void _main(void)
{
	HCRYPTPROV hcrypt;
	unsigned r;
	unsigned ret = 0;
	Console c;
	Arg coinmode, start, end;
	Arg *args[] = {
		&coinmode, &start, &end,
	};

	console_init(&c);
	args_make(&start, ArgValue, "-s");
	args_make(&end, ArgValue, "-e");
	args_make(&coinmode, ArgBool, "-c");
	args_get(args, ARR_LEN(args), GetCommandLine());

	if(!coinmode.isset && !(start.value && end.value)) {
		console_print_const(&c, "Usage of rand:\n\
  -s:  lower range end\n\
  -e:  upper range end\n\
  -c:  coin mode, only prints yes or no, and returns 1 or 0\n");
		goto end;
	}

	if(!CryptAcquireContext(&hcrypt, "quite random", NULL,
			PROV_RSA_FULL, 0)
	   && !CryptAcquireContext(&hcrypt, "quite random", NULL,
			PROV_RSA_FULL, CRYPT_NEWKEYSET)) {
		console_print_const(&c, "crypt api failed");
		goto end;
	}

	if(coinmode.isset) {
		CryptGenRandom(hcrypt, sizeof(unsigned), (BYTE *)&r);
		r = r & (1<<15);
		if(r) {
			console_print_const(&c, "yes\n");
		} else {
			console_print_const(&c, "no\n");
		}
		ret = r;
	} else {
		unsigned s = stou(start.value);
		unsigned e = stou(end.value);
		unsigned dist = e-s;
		if(dist == 0) {
			r = 0;
		} else {
			unsigned max;
			if(dist <= 0x80000000u) {
				unsigned left = (0x80000000u % dist) * 2;
				if(left >= dist) {
					left -= dist;
				}
				max = 0xffffffffu - left;
			} else {
				max = dist - 1;
			}
			do {
				CryptGenRandom(hcrypt, sizeof(unsigned), (BYTE *)&r);
			} while(r > max);
			r %= dist;
		}
		r = s+r;
		console_printu(&c, r);
		console_put(&c, '\n');
		ret = r;
	}
end:
	CryptReleaseContext(hcrypt, 0);
	console_deinit(&c);
	ExitProcess(ret);
}
