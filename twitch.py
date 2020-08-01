import json
import dao
import threading
import time
import logging
import userTwitch
from requests import post
from requests import get
from userTwitch import UserTwitch

PREFIX_DB = "twitch"
unique_users_collection = []


logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

logger = logging.getLogger(__name__)


class Twitch:
    token = None
    client_id = ""
    client_secret = ""
    users = []

    def __init__(self):
        super().__init__()
        dao.__init__()
        self.token = dao.get(PREFIX_DB, "token")
        self.client_id = dao.get(PREFIX_DB, "client_id")
        self.client_secret = dao.get(PREFIX_DB, "client_secret")

        userTwitch.TOKEN = self.token
        userTwitch.CLIENT_ID = self.client_id
        self.load_users()

        if not self.is_token_set() or self.need_token_renew():
            self.get_token()
        else:
            self.token_renew_worker()

    def is_token_set(self):
        return not self.token is None and "access_token" in self.token

    def are_there_users(self):
        return len(self.users) > 0

    def is_client_id_set(self):
        return not (self.client_id == "" or self.client_id is None)

    def is_client_secret_set(self):
        return not (self.client_secret == "" or self.client_secret is None)

    def set_token(self, token):
        self.token = token
        dao.save(PREFIX_DB, "token", token)

    def set_client_id(self, client_id):
        self.client_id = client_id
        userTwitch.CLIENT_ID = self.client_id
        dao.save(PREFIX_DB, "client_id", client_id)
        self.get_token()

    def set_client_secret(self, client_secret):
        self.client_secret = client_secret
        dao.save(PREFIX_DB, "client_secret", client_secret)
        self.get_token()

    def load_users(self):
        global unique_users_collection
        usersDao = dao.get(PREFIX_DB, "users")
        self.users = []
        if usersDao is None:
            return

        for user in usersDao:
            if not user["twitch_id"] is None and not user["twitch_id"] in unique_users_collection:
                unique_users_collection.append(user["twitch_id"])

            self.users.append(UserTwitch(
                user["username"], user["chat_id"], user["is_group"], user["twitch_id"], user["last_stream"], user["is_streaming"]))

    def get_token(self):
        if not self.is_client_id_set() or not self.is_client_secret_set():
            return
        url = "https://id.twitch.tv/oauth2/token?client_id={0}&client_secret={1}&grant_type=client_credentials".format(
            self.client_id, self.client_secret)
        response = post(url)
        if response.status_code != 200:
            return response.text

        self.token = json.loads(response.text)
        self.token["renewel_at"] = time.time()
        self.token["expires_in"] = int(self.token['expires_in'])
        userTwitch.TOKEN = self.token
        self.token_renew_worker()
        dao.save(PREFIX_DB, "token", self.token)

    def add_user(self, user, chat_id, is_group=False):
        user = user.replace("@", "", 1).lower()
        if self.users is None:
            self.users = []

        for i, v in enumerate(self.users):
            if v.username.lower() == user.lower() and v.chat_id == chat_id:
                return "The user @{0} is already configured".format(user)

        newUser = UserTwitch(user, chat_id, is_group)

        if not newUser.twitch_id is None and not newUser.twitch_id in unique_users_collection:
            unique_users_collection.append(newUser.twitch_id)

        self.users.append(newUser)
        dao.save(PREFIX_DB, "users", self.users)

    def need_token_renew(self):
        return self.remaining_token_expiration() < 10

    def remaining_token_expiration(self):
        return self.token["expires_in"] - (time.time() - self.token["renewel_at"])

    def token_renew_worker(self):
        currentTimer = threading.Timer(
            self.remaining_token_expiration(), self.get_token).start()

    def send_notfications(self, context):
        global unique_users_collection
        headers = {
            "client-id": self.client_id,
            "Authorization": "Bearer {0}".format(self.token["access_token"])
        }

        query_params = ""
        for twitch_id in unique_users_collection:
            query_params = "{0}&user_id={1}".format(
                query_params, twitch_id)

        url = 'https://api.twitch.tv/helix/streams?{0}'.format(
            query_params)
        response = get(
            url, headers=headers)

        if response.status_code != 200:
            logger.warning(response.text)
            return  # set master chat to recive log errors

        data = json.loads(response.text)
        if not "data" in data:
            logger.warning(response.text)
            return

        streaming_users = []
        for twitch_data in data["data"]:
            userName = twitch_data['user_name']
            streaming_users.append(userName.lower())
            streamTitle = twitch_data['title']
            viewerCount = twitch_data['viewer_count']
            gameID = twitch_data['game_id']
            startedAt = twitch_data['started_at']

            response = get(
                'https://api.twitch.tv/helix/games?id={0}'.format(gameID), headers=headers)
            data = json.loads(response.text)
            gameName = data['data'][0]['name']

            users = self.get_users_by_username(userName)
            for user in users:
                if not user.requires_notif():
                    continue

                retval = context.bot.send_message(chat_id=user.chat_id, text='[{0}](https://twitch.tv/{0}) is live!! Streaming {1} with {2} viewers'.format(
                    userName, gameName, viewerCount), parse_mode='MarkDown')

                user.set_is_streaming(True)
                dao.save(PREFIX_DB, "users", self.users)

                if user.is_group:
                    try:
                        context.bot.pin_chat_message(
                            chat_id=user.chat_id, message_id=retval.message_id)
                    except:
                        continue

        for user in self.users:
            if user.username in streaming_users:
                continue

            if not user.is_streaming:
                continue

            user.set_is_streaming(False)
            dao.save(PREFIX_DB, "users", self.users)
            context.bot.send_message(
                chat_id=user.chat_id, text='{0} stream is not running :('.format(user.username))

            if user.is_group:
                try:
                    context.bot.unpin_chat_message(chat_id=user.chat_id)
                except:
                    continue

    def get_users_by_username(self, username):
        users = []
        for user in self.users:
            if user.username.lower() == username.lower():
                users.append(user)

        return users
