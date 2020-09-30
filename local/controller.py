import requests
import time
import os
import sys
import json
import threading
from getpass import getpass
import schedule
import event as e
import configuration as c

#SERVER_URL = "https://home-automation-289621.uc.r.appspot.com"
SERVER_URL = "http://127.0.0.1:4747"

current_events = []


def turn_on(pin):
    print("Turn on " + str(pin))


def turn_off(pin):
    print("Turn off " + str(pin))


def schedule_job(event):
    if len(event.days) == 0 or len(event.days) == 7:
        schedule.every().day.at(str(event.time)).do(turn_on, event.pin)
    else:
        if 1 in event.days:
            schedule.every().monday.at(str(event.time)).do(turn_on, event.pin)
        if 2 in event.days:
            schedule.every().tuesday.at(str(event.time)).do(turn_on, event.pin)
        if 3 in event.days:
            schedule.every().wednesday.at(str(event.time)).do(turn_on, event.pin)
        if 4 in event.days:
            schedule.every().thursday.at(str(event.time)).do(turn_on, event.pin)
        if 5 in event.days:
            schedule.every().friday.at(str(event.time)).do(turn_on, event.pin)
        if 6 in event.days:
            schedule.every().saturday.at(str(event.time)).do(turn_on, event.pin)
        if 7 in event.days:
            schedule.every().sunday.at(str(event.time)).do(turn_on, event.pin)


def run_scheduling():
    while True:
        schedule.run_pending()
        time.sleep(1)


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
        r_dict["systemStatus"], r_dict["rainPercentage"], r_dict["defaultDuration"], r_dict["update"], r_dict["_id"])


def set_update_off(configuration):
    r = requests.put(
        SERVER_URL + "/configurations/set-off-update/" + configuration.id)

    if r.status_code >= 400:
        print(
            "Error updating configuration status... Possible reconfiguration on next cycle")
    else:
        print("Update set off.")


def update_schedules(user):
    r = requests.get(
        SERVER_URL + "/devices/by-user-with-events/" + user)

    devices = r.json()

    schedule.clear()

    for device in devices:
        for event in device["events"]:
            print(event)
            schedule_job(e.Event(
                device["pin"], event["days"], e.Time(event["time"]["hour"], event["time"]["minute"]), e.Repetition(event["repetition"]["times"], event["repetition"]["date"], event["repetition"]["current"]), event["duration"]))


def main():
    user = ""
    if not os.path.isfile("./config.txt"):
        print("No configuration found... Initializing configuration")
        user = initial_setup()
    else:
        print("Validating user...")
        user = get_user()

    print("Initializing routine...")

    # Initialize separate thread to run scheduling jobs
    thread = threading.Thread(None, run_scheduling, "Schedule")
    thread.run()

    while True:
        configuration = get_configuration(user)
        if configuration.update:
            print("Updating schedule...")
            update_schedules(user)
            set_update_off(configuration)
        time.sleep(1)

    thread.join()


if __name__ == "__main__":
    main()
