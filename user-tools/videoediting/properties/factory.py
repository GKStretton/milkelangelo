import machinepb.machine as pb

from videoediting.properties.content_property_manager import BasePropertyManager
from videoediting.properties.longform import LongFormPropertyManager
from videoediting.properties.shortform import ShortFormPropertyManager
from videoediting.properties.cleaning import CleaningPropertyManager


def create_property_manager(content_type: pb.ContentType) -> BasePropertyManager:
    if content_type == pb.ContentType.CONTENT_TYPE_LONGFORM:
        return LongFormPropertyManager()
    elif content_type == pb.ContentType.CONTENT_TYPE_SHORTFORM:
        return ShortFormPropertyManager()
    elif content_type == pb.ContentType.CONTENT_TYPE_CLEANING:
        return CleaningPropertyManager()
    else:
        raise ValueError(f"Invalid content type: {content_type}")
