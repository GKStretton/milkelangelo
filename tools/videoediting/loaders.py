import logging
import typing
import os
from dataclasses import dataclass

import numpy as np
from PIL import Image
import machinepb.machine as pb
import pycommon.image as image

import yaml


def get_session_metadata(base_dir: str, session_number: int):
    filename = "{}_session.yml".format(session_number)
    path = os.path.join(base_dir, "session_metadata", filename)
    yml = None
    with open(path, 'r') as f:
        yml = yaml.load(f, Loader=yaml.FullLoader)
    logging.info("Loaded session metadata: {}\n".format(yml))
    return yml


def get_session_content_path(base_dir: str, session_number: int):
    return os.path.join(base_dir, "session_content", str(session_number))


def get_state_reports(base_dir: str, session_number: int) -> typing.Tuple[float, pb.StateReport]:
    content_path = get_session_content_path(base_dir, session_number)
    raw_state_reports = None
    with open(os.path.join(content_path, "state-reports.yml"), 'r') as f:
        raw_state_reports = yaml.load(f, yaml.FullLoader)

    state_report_list = pb.StateReportList().from_dict(raw_state_reports).state_reports
    logging.info("Loaded {} state report entries\n".format(len(state_report_list)))

    state_reports: typing.Tuple[float, pb.StateReport] = []
    for _, v in enumerate(state_report_list):
        report = v
        report_ts = float(report.timestamp_unix_micros) / 1.0e6
        state_reports.append((report_ts, report))

    return state_reports


def get_content_plan(base_dir: str, session_number: int) -> pb.ContentTypeStatuses:
    content_path = get_session_content_path(base_dir, session_number)
    with open(os.path.join(content_path, "content_plan.yml"), 'r') as f:
        raw_plan = yaml.load(f, yaml.FullLoader)

    return pb.ContentTypeStatuses().from_dict(raw_plan)


def get_still_dir(base_dir: str, session_number: int) -> str:
    return os.path.join(get_session_content_path(base_dir, session_number), "still")


# frame selected from the top-cam recording in the interface
def get_selected_still_path(base_dir: str, session_number: int) -> str:
    return os.path.join(get_still_dir(base_dir, session_number), "selected.jpg")


# selected still with the vignette overlay applied, as an RGB array
def get_selected_still_image(base_dir: str, session_number: int) -> np.ndarray:
    img = np.array(Image.open(get_selected_still_path(base_dir, session_number)).convert("RGB"))
    image.add_overlay(img)
    return img


# unix timestamp (s) of the frame selected in the interface
def get_selected_still_timestamp(base_dir: str, session_number: int) -> float:
    with open(os.path.join(get_still_dir(base_dir, session_number), "selected.yml"), 'r') as f:
        selected = yaml.load(f, yaml.FullLoader)

    return float(selected['timestamp_unix'])


@dataclass
class MiscData:
    """Miscellaneous data to be passed to the pipeline"""
    selected_timestamp_s: float = 0
    start_timestamp_s: int = 0


def get_misc_data(base_dir: str, session_number: int, state_reports: typing.Tuple[float, pb.StateReport]) -> MiscData:
    return MiscData(
        selected_timestamp_s=get_selected_still_timestamp(base_dir, session_number),
        start_timestamp_s=state_reports[0][0]
    )


def get_profile_snapshot(base_dir: str, session_number: int) -> pb.SystemVialConfigurationSnapshot:
    content_path = get_session_content_path(base_dir, session_number)
    with open(os.path.join(content_path, "vial-profiles.yml"), 'r') as f:
        raw_profiles = yaml.load(f, yaml.FullLoader)

    return pb.SystemVialConfigurationSnapshot().from_dict(raw_profiles)
