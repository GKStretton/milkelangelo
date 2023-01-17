from datetime import datetime

def ts_format(ts: float) -> str:
	return datetime.utcfromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')
