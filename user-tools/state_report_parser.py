# This parses the state reports from a session and sets the dispense_request_number
# Used for sessions recorded before the dispense_request_number was added to the state report

import yaml
import machinepb.machine as pb

INPUT = "/mnt/md0/light-stores/session_content/59/state-reports.yml"
OUTPUT = "/mnt/md0/light-stores/session_content/59/state-reports_modified.yml"

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

isDispensing = False
dispenseCounter = 0
for i in range(len(reports)):
	print(i)
	report = pb.StateReport().from_dict(reports[i])

	if not isDispensing and report.status == pb.Status.DISPENSING:
		isDispensing = True
		dispenseCounter += 1
	
	if isDispensing and report.status != pb.Status.DISPENSING:
		isDispensing = False

	report.pipette_state.dispense_request_number = dispenseCounter

	new_reports.append(report.to_dict(include_default_values=True))

with open(OUTPUT, 'w') as f:
	yaml.dump(new_reports, f)