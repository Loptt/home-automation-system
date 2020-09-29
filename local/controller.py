import requests
import time
import os
import sys
import json
from getpass import getpass
import event as e
import configuration as c

#SERVER_URL = "https://home-automation-289621.uc.r.appspot.com"
SERVER_URL = "http://127.0.0.1:4747"

current_events = []


def initial_setup():
    username = input("Enter your username: ")
    password = getpass("Enter your password: ")

    pload = json.dumps({"username": username, "password": password})

    r = requests.post(SERVER_URL + "/login", data=pload,
                      headers={'Content-type': 'application/json'})
    r_dict = r.json()

    if not r_dict["valid"]:
        print("Invalid username/password")
        print("Run program again to try again")
        sys.exit()

    print("Successful login...")
    print("Saving configuration...")
    f = open("config.txt", "w")
    f.write(r_dict["user"])
    f.close()

    return r_dict["user"]


def get_user():
    f = open("config.txt", "r")
    user = f.readline()

    r = requests.get(SERVER_URL + "/users/" + user)

    if r.status_code == 200:
        print("Successful login...")
        return user
    else:
        print("Invalid user... Reinitializing configuration")
        return initial_setup()


def get_configuration(user):

    r = requests.get(SERVER_URL + "/configurations/by-user/" + user)

    if r.status_code != 200:
        print("Error retrieving configuration, check internet connection.")
        sys.exit()

    r_dict = r.json()

    return c.Configuration(
        r_dict["systemStatus"], r_dict["rainPercentage"], r_dict["defaultDuration"], r_dict["update"])


def update_schedules(user):
    r = requests.get(
        SERVER_URL + "/devices/by-user-with-events/" + user)

    devices = r.json()

    for device in devices:
        for event in device["events"]:
            


def main():
    user = ""
    if not os.path.isfile("./config.txt"):
        print("No configuration found... Initializing configuration")
        user = initial_setup()
    else:
        print("Validating user...")
        user = get_user()

    print("Initializing routine...")

    while True:
        configuration = get_configuration(user)
        if configuration.update:
            print("Updating schedule...")
            update_schedules(user)
        time.sleep(1)


if __name__ == "__main__":
    main()
