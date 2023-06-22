import yaml

class CropConfig:
	def __init__(self, path):
		self.is_loaded = False
		yml = self._load_crop_config(path)
		self.raw_yml = yml
		if yml is None:
			return
		
		self.is_loaded = True
		self.x1 = yml['left_abs']
		self.x2 = yml['right_abs']
		self.y1 = yml['top_abs']
		self.y2 = yml['bottom_abs']

	def _load_crop_config(self, path):
		config = None
		try:
			with open(path, 'r') as f:
				config = yaml.load(f, Loader=yaml.FullLoader)
		except FileNotFoundError as err:
			print("error loading crop config at '{}': {}", path, err)

		return config