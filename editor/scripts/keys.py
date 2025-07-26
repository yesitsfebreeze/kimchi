import sys
import tty
import termios
from contextlib import contextmanager

@contextmanager
def raw_mode():
	fd = sys.stdin.fileno()
	old_settings = termios.tcgetattr(fd)
	try:
		tty.setraw(fd)
		sys.stdout.write("\033[?1000h")
		sys.stdout.flush()
		yield
	finally:
		termios.tcsetattr(fd, termios.TCSADRAIN, old_settings)

def get_key():
	# not unicode necessarily
	return sys.stdin.buffer.read(1)

def identify_key(char):
	try:
		char = char.decode("utf-8")
	except UnicodeDecodeError:
		pass
	if ord(char) < 32:
		return f"ctrl-{chr(ord(char) + 64)}"
	elif ord(char) == 127:
		return "backspace"
	elif ord(char) <= 127:
		return char
	elif ord(char) >= 127 and ord(char) <= 255:
		return f"alt-{chr(ord(char) - 225 + 97)}"
	else:
		return ord(char)

def main():
	print("Press keys (Ctrl-C to exit):")
	with raw_mode():
		while True:
			char = get_key()
			key_name = identify_key(char)
			print(f"{key_name}", end="")
			sys.stdout.flush()
			if char == b"\x03":
				return


if __name__ == "__main__":
	main()
