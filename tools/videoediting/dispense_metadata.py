import yaml
import logging
from typing import Optional
import machinepb.machine as pb
import os


class DispenseMetadataWrapper:
    def __init__(self, base_dir: str, session_number: int):
        self.base_dir = base_dir
        self.session_number = session_number

        self.metadata = self.load_dispense_metadata(base_dir, session_number)

    def load_dispense_metadata(self, base_dir: str, session_number: int) -> Optional[pb.DispenseMetadataMap]:
        path = os.path.join(base_dir, "session_content", str(session_number), "dispense-metadata.yml")
        try:
            with open(path, 'r') as f:
                yml = yaml.load(f, Loader=yaml.FullLoader)
                meta: pb.DispenseMetadataMap = pb.DispenseMetadataMap().from_dict(value=yml)
                return meta
        except FileNotFoundError:
            logging.error("no dispense metadata found")
            return None

    def get_dispense_metadata(self, startup_counter: int, dispense_request_number: int) -> Optional[pb.DispenseMetadata]:
        if self.metadata is None:
            return None

        key = f"{startup_counter}_{dispense_request_number}"
        return self.metadata.dispense_metadata.get(key)

    def get_dispense_metadata_from_sr(self, sr: pb.StateReport) -> Optional[pb.DispenseMetadata]:
        return self.get_dispense_metadata(sr.startup_counter, sr.pipette_state.dispense_request_number)
