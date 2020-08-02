import json
import time
from requests import get
from datetime import datetime

TOKEN = None
CLIENT_ID = None
'''
Minimum of seconds to send again the notification, each 21600 secs = 6 hours
'''
NOTIF_DELAY = 6*60*60


class UserTwitch:
    username = None
    twitch_id = None
    chat_id = None
    is_streaming = False
    is_group = None
    last_stream = 0

    def default(self, o):
        return o.__dict__

    def __init__(self, username, chat_id, is_group=False, twitch_id=None, last_stream=None, is_streaming=False):
        super().__init__()
        self.username = username
        self.chat_id = chat_id
        self.is_group = is_group
        self.is_streaming = is_streaming
        self.last_stream = last_stream
        if twitch_id is None:
            self.get_user_twitch_id()
        else:
            self.twitch_id = twitch_id

    def to_dict(self):
        return {
            "username": self.username,
            "twitch_id": self.twitch_id,
            "chat_id": self.chat_id,
            "is_group": self.is_group,
            "is_streaming": self.is_streaming,
            "last_stream": self.last_stream
        }

    def get_user_twitch_id(self):
        global TOKEN, CLIENT_ID
        if TOKEN is None:
            return
        headers = {
            "client-id": CLIENT_ID,
            "Authorization": "Bearer {0}".format(TOKEN["access_token"])
        }

        url = "https://api.twitch.tv/helix/users?login={0}".format(
            self.username)
        response = get(url, headers=headers)
        if response.status_code != 200:
            return response.text
        dataRes = json.loads(response.text)
        if not "data" in dataRes or len(dataRes["data"]) < 1 or not "id" in dataRes["data"][0]:
            self.twitch_id = None
            return

        self.twitch_id = dataRes["data"][0]["id"]

    def set_is_streaming(self, is_streaming):
        if is_streaming:
            self.last_stream = datetime.utcnow().timestamp()

        self.is_streaming = is_streaming

    def requires_notif(self):
        if self.last_stream is None or self.last_stream == 0:
            return True

        last_stream = datetime.fromtimestamp(self.last_stream)
        delta = datetime.utcnow() - last_stream
        return delta.seconds > NOTIF_DELAY
