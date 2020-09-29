

class Time:

    def __init__(self, hour, minute):
        self.hour = hour
        self.minute = minute

    def __eq__(self, other):
        return self.hour == other.hour and self.minute == other.minute


class Repetition:

    def __init__(self, times, date, current):
        self.times = times
        self.date = date
        self.current = current

    def __eq__(self, other):
        return self.times == other.times and self.date == other.date and self.current == other.current


class Event:

    def __init__(self, pin, days, time, repetition, duration):
        self.pin = pin
        self.days = days
        self.time = time
        self.repetition = repetition
        self.duration = duration

    def __eq__(self, other):
        return self.pin == other.pin and self.days == other.days and self.time == other.time and self.repetition == other.repetition and self.duration == other.duration
