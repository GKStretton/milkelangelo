from datetime import datetime
import typing

def ts_format(ts: float) -> str:
	if ts is None:
		return "None"
	return datetime.utcfromtimestamp(ts).strftime('%Y-%m-%d %H:%M:%S')

def ts_fmt(ts: float) -> str:
	if ts is None:
		return "None"
	return datetime.utcfromtimestamp(ts).strftime('%H:%M:%S')+".{}".format(int(ts%1*100))

def floats_are_equal(threshold: float, nums: typing.List[float]):
	if len(nums) == 0:
		return True
	
	value = nums[0]
	for i in range(len(nums)):
		if abs(nums[i] - value) > threshold:
			return False
	
	return True