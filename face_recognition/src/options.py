import sys
from enroll import enroll_user
from verify import verify_user

def main():
    if len(sys.argv) < 2:
        print("Usage: python options.py [enroll|verify]")
        return

    cmd = sys.argv[1]
    if cmd == "enroll":
        enroll_user()
    elif cmd == "verify":
        verify_user()
    else:
        print("Unknown command. Use 'enroll' or 'verify'.")

if __name__ == "__main__":
    main()
