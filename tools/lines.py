import os

lines = 0

# bot dir doesnt have any of subdirs so i can just use listdir
files = os.listdir("bot")

for file in files:
    file = "bot/" + file

# and dont forget main.go
files.append("main.go")


for file in files:
    with open(file) as f:
        # this is probably one of the worst ways to do that
        # but uhh i dont really care
        lines += len(f.readlines())

print(f"Total lines: {lines}")
