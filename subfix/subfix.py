import io
import re
import sys
import argparse

parser = argparse.ArgumentParser()
parser.add_argument('-f', dest='file', type=str, help='file',required=True)
parser.add_argument('-s', dest='seconds', type=int, default=0, help='seconds')
parser.add_argument('-ms', dest='milliseconds', type=int, default=0, help='milliseconds')
args = parser.parse_args()

reTime = re.compile(r"(\d{2}):(\d{2}):(\d{2}),(\d{3}) --> (\d{2}):(\d{2}):(\d{2}),(\d{3})")

with io.open(args.file, "rb") as fd:
	sub = fd.read().decode('utf-8')

def timeMod(hour, min, sec, ms):
	ms += args.milliseconds
	sec += args.seconds
	if ms >= 1000:
		ms %= 1000
		sec += 1
	if sec >= 60:
		sec %= 60
		min += 1
	if min >= 60:
		min %= 60
		hour += 1
	return (hour, min, sec, ms)

def timeRepl(m):
	hour1, min1, sec1, ms1 = [int(x, 10) for x in m.groups()[:4]]
	hour1, min1, sec1, ms1 = timeMod(hour1, min1, sec1, ms1)
	hour2, min2, sec2, ms2 = [int(x, 10) for x in m.groups()[4:]]
	hour2, min2, sec2, ms2 = timeMod(hour2, min2, sec2, ms2)
	return "{:02}:{:02}:{:02},{:03} --> {:02}:{:02}:{:02},{:03}".format(
		hour1, min1, sec1, ms1,
		hour2, min2, sec2, ms2
	)

sys.stdout.write(reTime.sub(timeRepl, sub))
