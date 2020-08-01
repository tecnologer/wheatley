import configparser
import logging
from twitch import Twitch
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters
from telegram import InlineKeyboardButton, InlineKeyboardMarkup

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

logger = logging.getLogger(__name__)

# constans
DATA_FILE_PATH = './data.ini'
NOTIFICATION_DELAY = 60


isConfigClientSecret = False
isConfigClientId = False
isAddingUser = False

telegram_whiteList = []

telegram_botToken = ""
updater = None
t = None

commands = None


def is_allowed(update):
    '''Checks if the telegram's user is allowed to execute the command'''
    global telegram_whiteList
    user = update.effective_user.username.lower()
    if user in (name.lower() for user in telegram_whiteList):
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Permissions required")
        return False

    return True


def get_param_value(command, text):
    '''returns a collection with the values after the command'''
    if not text.startswith(command):
        return None
    text = text.replace(command, "")
    text = text.strip()
    if text == "":
        return None
    return text.split(" ")


def check_missing_data(update, context):
    if not t.is_client_id_set():
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Before continue, We need set the twitch client id.")
        handle_twitch_client_id(update, context)
        return True

    if not t.is_client_secret_set():
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Also, We need set the twitch client secret.")
        handle_twitch_client_secret(update, context)
        return True

    if not t.are_there_users():
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="It's required at least one User's twitch.")
        handle_twitch_add_user(update, context)
        return True

    return False


def start(update, context):
    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Hello, I'll assist you to configure everything")

    if check_missing_data(update, context):
        return

    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Congratulations!, everything is configured.")


def handle_twitch_client_id(update, context):
    global isConfigClientId, isConfigClientSecret, isAddingUser
    isConfigClientSecret = False
    isAddingUser = False
    isConfigClientId = True
    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Please type the twitch Client ID:")


def handle_twitch_client_secret(update, context):
    if not is_allowed(update):
        return
    global isConfigClientId, isConfigClientSecret, isAddingUser
    isConfigClientSecret = True
    isAddingUser = False
    isConfigClientId = False

    users = get_param_value("/setclientid", update.message.text)

    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Please type the twitch Client secret:")


def generic_handle(update, context):
    global isConfigClientId, isConfigClientSecret, isAddingUser

    msgText = update.message.text
    if msgText.startswith("/") or msgText == "" or msgText is None:
        return

    if isConfigClientId:
        isConfigClientId = False
        t.set_client_id(msgText)

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client Id configured")
        return

    if isConfigClientSecret:
        isConfigClientSecret = False
        t.set_client_secret(msgText)

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client secret configured")
        return

    if isAddingUser:
        isAddingUser = False
        res = t.add_user(msgText, update.effective_chat.id,
                         update.effective_chat.type != 'private')

        if res is not None:
            context.bot.send_message(
                chat_id=update.effective_chat.id, text=res)
            return

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Now I'll notify you in this chat when @{0} is streaming".format(msgText))
        return


def handle_cancel(update, context):
    isConfigClientSecret = False
    isConfigClientId = False
    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Done!, everything is canceled:")


def handle_error(update, context):
    """Log Errors caused by Updates."""
    logger.warning('Update "%s" caused error "%s"', update, context.error)


def handle_help(update, context):
    global commands
    helpMsg = u"Available commands: \n\n"
    for i, v in enumerate(commands):
        if v["inHelp"]:
            helpMsg = helpMsg + \
                u"- /{0}: {1}\n".format(v["command"], v["info"])

    context.bot.send_message(
        chat_id=update.effective_chat.id, text=helpMsg)


def handle_twitch_add_user(update, context):
    global isConfigClientId, isConfigClientSecret, isAddingUser
    isConfigClientSecret = False
    isConfigClientId = False
    isAddingUser = False
    users = get_param_value("/adduser", update.message.text)
    if len(user) == 0:
        isAddingUser = True
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Type the username of User's twitch:")
    else:
        for user in users:
            res = t.add_user(user, update.effective_chat.id,
                             update.effective_chat.type != update.effective_chat.PRIVATE)
            if res is not None:
                context.bot.send_message(
                    chat_id=update.effective_chat.id, text=res)
                return

        isOrAre = "is" if len(users) == 1 else "are"
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Now I'll notify you in this chat when @{0} {1} streaming".format(", ".join(users), isOrAre))


def __init__():
    global t, commands, telegram_whiteList
    t = Twitch()
    config = configparser.ConfigParser()

    logger.info("reading bot configuration")
    config.read(DATA_FILE_PATH)
    telegram_botToken = config.get("telegram", "botToken")
    telegram_whiteList = config.get("telegram", "whiteList").split(",")

    logger.info("creating bot updater")
    updater = Updater(telegram_botToken, use_context=True)
    dp = updater.dispatcher
    queue = updater.job_queue

    commands = [
        {
            "command": "start",
            "handle": start,
            "info": "Starts the bot",
            "inHelp": False
        },
        {
            "command": "setclientid",
            "handle": handle_twitch_client_id,
            "info": "Stores the twitch client id",
            "inHelp": True
        },
        {
            "command": "setsecretid",
            "handle": handle_twitch_client_secret,
            "info": "Stores the twitch client secret",
            "inHelp": True
        },
        {
            "command": "help",
            "handle": handle_help,
            "info": "Shows the basic info of commands",
            "inHelp": False
        },
        {
            "command": "cancel",
            "handle": handle_cancel,
            "info": "Cancels the active command",
            "inHelp": True
        },
        {
            "command": "adduser",
            "handle": handle_twitch_add_user,
            "info": "Adds a new user(s) to the list to monitor its stream status. Use users separated by space to add multiple.",
            "inHelp": True
        }
    ]

    logger.info("configuring handlers")

    for i, v in enumerate(commands):
        dp.add_handler(CommandHandler(v["command"], v["handle"]))

    dp.add_error_handler(handle_error)
    dp.add_handler(MessageHandler(Filters.text, generic_handle))

    logger.info("configuring workers")
    queue.run_repeating(t.send_notfications,
                        interval=NOTIFICATION_DELAY, first=0)

    logger.info("starting polling")
    updater.start_polling()

    logger.info("Bot ready...")
    updater.idle()
