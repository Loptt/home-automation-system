

class Time:

    def __init__(self, hour, minute):
        self.hour = hour
        self.minute = minute

    def __eq__(self, other):
        return self.hour == other.hour and self.minute == other.minute

    def __str__(self):
        hour = str(self.hour)
        minute = str(self.minute)

        if self.hour < 10:
            hour = "0" + hour
        if self.minute < 10:
            minute = "0" + minute

        return hour + ":" + minute


class Repetition:

    def __init__(self, times, date, current):
        self.times = times
        self.date = date
        self.current = current

    def __eq__(self, other):
        return self.times == other.times and self.date == other.date and self.current == other.current

    def __str__(self):
        return "  Times: " + str(self.times) + "\n  Date: " + str(self.date) + "\n  Current: " + str(self.current)


class Event:

    def __init__(self, pin, days, time, repetition, duration):
        self.pin = pin
        self.days = days
        self.time = time
        self.repetition = repetition
        self.duration = duration

    def __eq__(self, other):
        return self.pin == other.pin and self.days == other.days and self.time == other.time and self.repetition == other.repetition and self.duration == other.duration

    def __str__(self):
        return "Pin: " + str(self.pin) + "\nDays: " + str(self.days) + "\nTime: " + str(self.time) + "\nRepetition: " + str(self.repetition) + "\nDuration: " + str(self.duration)
