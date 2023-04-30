from videoediting.constants import *
from videoediting.properties.content_property_manager import *
from videoediting.properties.longform import LongFormPropertyManager
from videoediting.properties.shortform import ShortFormPropertyManager
from videoediting.properties.cleaning import CleaningPropertyManager

def create_property_manager(content_type: ContentType) -> BasePropertyManager:
	if content_type == ContentType.LONGFORM:
		return LongFormPropertyManager()
	elif content_type == ContentType.SHORTFORM:
		return ShortFormPropertyManager()
	elif content_type == ContentType.CLEANING:
		return CleaningPropertyManager()
	else:
		raise ValueError(f"Invalid content type: {content_type}")