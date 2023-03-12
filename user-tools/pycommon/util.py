from datetime import datetime

def ts_format(ts: float) -> str:
	if ts is None:
		return "None"
	return datetime.utcfromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')

def ts_fmt(ts: float) -> str:
	if ts is None:
		return "None"
	return datetime.utcfromtimestamp(ts).strftime('%H:%M:%S')+".{}".format(int(ts%1*100))