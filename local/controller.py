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
import RPi.GPIO as GPIO

#SERVER_URL = "https://home-automation-289621.uc.r.appspot.com"
#SERVER_URL = "http://127.0.0.1:4747"
SERVER_URL = "http://192.168.11.117:4747"

current_events = []


def calculate_max_duration(time):
    hours = 23 - time.hour
    minutes = 60 - time.minute

    return hours * 60 + minutes


def turn_on(pin):
    print("Turn on " + str(pin))
    GPIO.output(pin, GPIO.HIGH)


def turn_off(pin):
    print("Turn off " + str(pin))
    GPIO.output(pin, GPIO.LOW)


def schedule_off(time, day, duration, pin):
    new_day = day
    end_time = e.Time(0, 0)

    if duration > calculate_max_duration(time):
        # Next day calculation
        new_day = day + 1
        off_duration = duration - calculate_max_duration(time)
        end_time.hour = off_duration // 60
        end_time.minute = off_duration % 60
    else:
        # Same day calculation
        end_time.hour = time.hour + \
            (duration // 60) + (time.minute + (duration % 60)) // 60
        end_time.minute = (time.minute + duration % 60) % 60

    if new_day > 7:
        new_day = 1

    if new_day == 1:
        schedule.every().monday.at(str(end_time)).do(turn_off, pin)
    elif new_day == 2:
        schedule.every().tuesday.at(str(end_time)).do(turn_off, pin)
    elif new_day == 3:
        schedule.every().wednesday.at(str(end_time)).do(turn_off, pin)
    elif new_day == 4:
        schedule.every().thursday.at(str(end_time)).do(turn_off, pin)
    elif new_day == 5:
        schedule.every().friday.at(str(end_time)).do(turn_off, pin)
    elif new_day == 6:
        schedule.every().saturday.at(str(end_time)).do(turn_off, pin)
    elif new_day == 7:
        schedule.every().sunday.at(str(end_time)).do(turn_off, pin)


def schedule_job(event):
    GPIO.setup(event.pin, GPIO.OUT)

    if len(event.days) == 0 or len(event.days) == 7:
        schedule.every().day.at(str(event.time)).do(turn_on, event.pin)
    else:
        if 1 in event.days:
            schedule.every().monday.at(str(event.time)).do(turn_on, event.pin)
            schedule_off(event.time, 1, event.duration, event.pin)
        if 2 in event.days:
            schedule.every().tuesday.at(str(event.time)).do(turn_on, event.pin)
            schedule_off(event.time, 2, event.duration, event.pin)
        if 3 in event.days:
            schedule.every().wednesday.at(str(event.time)).do(turn_on, event.pin)
            schedule_off(event.time, 3, event.duration, event.pin)
        if 4 in event.days:
            schedule.every().thursday.at(str(event.time)).do(turn_on, event.pin)
            schedule_off(event.time, 4, event.duration, event.pin)
        if 5 in event.days:
            schedule.every().friday.at(str(event.time)).do(turn_on, event.pin)
            schedule_off(event.time, 5, event.duration, event.pin)
        if 6 in event.days:
            schedule.every().saturday.at(str(event.time)).do(turn_on, event.pin)
            schedule_off(event.time, 6, event.duration, event.pin)
        if 7 in event.days:
            schedule.every().sunday.at(str(event.time)).do(turn_on, event.pin)
            schedule_off(event.time, 7, event.duration, event.pin)


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
    GPIO.cleanup()

    for device in devices:
        for event in device["events"]:
            print(event)
            schedule_job(e.Event(
                device["pin"], event["days"], e.Time(event["time"]["hour"], event["time"]["minute"]), e.Repetition(event["repetition"]["times"], event["repetition"]["date"], event["repetition"]["current"]), event["duration"]))


def setup():
    GPIO.setmode(GPIO.BCM)
    GPIO.setwarnings(False)


def main():
    setup()
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
    thread.start()

    print("Schedule running.")

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
