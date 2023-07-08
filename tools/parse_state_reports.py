# This parses the state reports from a session and sets the dispense_request_number
# Used for sessions recorded before the dispense_request_number was added to the state report

import yaml
from betterproto import Casing
import machinepb.machine as pb

INPUT = "/home/greg/Downloads/session_content/59/state-reports.yml"
OUTPUT = "/home/greg/Downloads/session_content/59/state-reports_modified.yml"


def get_state_reports():
    state_reports = None
    with open(INPUT, 'r') as f:
        state_reports = yaml.load(f, yaml.FullLoader)
    print("Loaded {} state report entries\n".format(len(state_reports)))
    return state_reports


print("in:", INPUT)
print("out:", OUTPUT)
print()

reports = get_state_reports()
new_reports = []

for i in range(len(reports)):
    print(i)
    report = pb.StateReport().from_dict(reports[i])

    if report.timestamp_unix_micros / 1000000.0 > 1681149877.765481586:
        report.latest_dslr_file_number = 90 + 5
    else:
        report.latest_dslr_file_number = 1

    new_reports.append(report.to_dict(include_default_values=True, casing=Casing.SNAKE))

with open(OUTPUT, 'w') as f:
    yaml.dump(new_reports, f)
