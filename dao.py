import json
import pickledb

DB = None


def save(prefix_db, key, value):
    if prefix_db != "":
        key = "{0}.{1}".format(prefix_db, key)

    if str(type(value)) == "<class 'list'>":
        dumpValue = []
        for i, v in enumerate(value):
            if str(type(v)).startswith("<class"):
                dumpValue.append(v.to_dict())

        DB.set(key, dumpValue)
    else:
        DB.set(key, value)
    DB.dump()


def get(prefix_db, key):
    if prefix_db != "":
        key = "{0}.{1}".format(prefix_db, key)
    value = DB.get(key)
    if not value:
        return None

    return value


def __init__():
    global DB
    DB = pickledb.load('./twitch_bot.json', False)
