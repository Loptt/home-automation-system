import requests
import time
import os
import sys
import json
from getpass import getpass

SERVER_URL = "https://home-automation-289621.uc.r.appspot.com"
#SERVER_URL = "http://127.0.0.1:4747"


def initial_setup():
    print("No configuration found... Initializing configuration")
    username = input("Enter your username: ")
    password = getpass("Enter your password: ")

    pload = json.dumps({"username": username, "password": password})

    r = requests.post(SERVER_URL + "/login", data=pload,
                      headers={'Content-type': 'application/json'})
    r_dict = r.json()

    print(r_dict)

    if not r_dict["valid"]:
        print("Invalid username/password")
        print("Run program again to try again")
        sys.exit()

    f = open("config.txt", "w")
    f.write(r_dict["user"])
    f.close()

    return r_dict["user"]


def get_user():
    f = open("config.txt", "r")
    return f.readline()


def main():
    user = ""
    if not os.path.isfile("./config.txt"):
        user = initial_setup()
    else:
        user = get_user()

    while True:
        devices = requests.get(SERVER_URL + "/devices/by-user/" + user)
        print(devices.json())
        time.sleep(1)


if __name__ == "__main__":
    main()
