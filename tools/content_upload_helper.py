import argparse
import datetime
from machinepb import machine as pb
from content_plan.loader import *

if __name__ == "__main__":
	parser = argparse.ArgumentParser()
	parser.add_argument("-d", "--base-dir", action="store", help="base directory containing session_content and session_metadata", required=True)
	parser.add_argument("-n", "--session-number", action="store", help="session number e.g. 5", required=True)
	parser.add_argument("-i", "--ignore-uploaded", action="store_true", help="if true, skip posts that are marked as uploaded", default=True)
	args = parser.parse_args()

	content_statuses = load_content_statuses(args.base_dir, args.session_number)

	for content_type in pb.ContentType:
		if content_type.name not in content_statuses.content_statuses:
			continue
		
		content_status = content_statuses.content_statuses[content_type.name]
		for post in content_status.posts:
			if args.ignore_uploaded and post.uploaded:
				continue

			print("\n"*40)
			print('*'*20)

			print(content_type.name)
			print(post.platform.name)
			print(f"UPLOADED: {'YES' if post.uploaded else 'NO'}")
			schedule_str = datetime.fromtimestamp(post.scheduled_unix_timetamp).strftime('%Y-%m-%d %H:%M:%S')
			print(f"SCHEDULE: { 'IMMEDIATE' if post.scheduled_unix_timetamp == 0 else f'{schedule_str} Local Time'}")

			print(f"\n{'*'*10} TITLE {'*'*10}\n")
			print(post.title)

			print(f"\n\n{'*'*10} DESCRIPTION {'*'*10}\n")
			print(post.description)
			print("\n\n")

			resp = '...'
			while resp != 'x' and resp != '':
				resp = input("enter 'x' to mark as not uploaded, or enter nothing to mark uploaded: ")

			if resp == 'x':
				post.uploaded = False
				print("post marked as not uploaded")
			elif resp == '':
				post.uploaded = True
				print("post marked as uploaded")
	
	print("\n\nDone, writing content statuses...")
	write_content_statuses(content_statuses, args.base_dir, args.session_number)
	print("Exiting.")
	input()