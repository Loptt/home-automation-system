import requests
import time


def main():
    while True:
        r = requests.get("https://home-automation-289621.uc.r.appspot.com")
        print(r.json())
        time.sleep(1)


if __name__ == "__main__":
    main()
