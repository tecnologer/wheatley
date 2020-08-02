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
isNotifWorkerRunning = False

telegram_whiteList = []
telegram_masterchat = 0

telegram_botToken = ""
updater = None
t = None

commands = None


def updateData(key, value):
    global DATA_FILE_PATH
    key = key.split("_")
    config = configparser.ConfigParser()
    config.read(DATA_FILE_PATH)
    cfgfile = open(DATA_FILE_PATH, 'w')
    config.set(key[0], key[1], str(value))
    config.write(cfgfile)
    cfgfile.close()


def configure_notif_workers():
    global updater, isNotifWorkerRunning
    if not t.is_client_id_set() or not t.is_client_secret_set() or isNotifWorkerRunning:
        return
    queue = updater.job_queue
    queue.run_repeating(t.send_notfications,
                        interval=NOTIFICATION_DELAY, first=0)
    isNotifWorkerRunning = True
    logger.info("Worker's notification is running")


def is_allowed(update, context):
    '''Checks if the telegram's user is allowed to execute the command'''
    global telegram_whiteList
    user = update.effective_user.username.lower()
    if not user.startswith("@"):
        user = "@{0}".format(user)

    if not user in (user.lower() for user in telegram_whiteList):
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Permissions required")
        return False

    return True


def reset_flags():
    global isConfigClientId, isConfigClientSecret, isAddingUser
    isAddingUser = False
    isConfigClientId = False
    isConfigClientSecret = False


def get_param_value(update, command):
    '''returns a collection with the values after the command'''
    text = update.message.text
    if not text.startswith(command):
        return []

    botName = update.effective_user.bot.username.lower()

    text = text.lower().replace(command.lower(), "").replace(
        botName, "").replace("@", "", -1)
    text = text.strip()
    if text == "":
        return []
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
    if not is_allowed(update, context):
        notify_to_master(update, context, "start")
        return
    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Hello, I'll assist you to configure everything")

    if check_missing_data(update, context):
        return

    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Congratulations!, everything is configured.")


def handle_twitch_client_id(update, context):
    if not is_allowed(update, context):
        notify_to_master(update, context, "setclientid")
        return

    global isConfigClientId
    reset_flags()

    clientIds = get_param_value(update, "/setclientid")

    if len(clientIds) == 0:
        isConfigClientId = True
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Please type the twitch Client ID:")
    else:
        clientId = clientIds[0]
        t.set_client_id(clientId)
        configure_notif_workers()

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client Id configured")
        return


def handle_twitch_client_secret(update, context):
    if not is_allowed(update, context):
        notify_to_master(update, context, "setsecretid")
        return
    global isConfigClientSecret
    reset_flags()

    secrets = get_param_value(update, "/setsecretid", update.message.text)

    if len(secrets) == 0:
        isConfigClientSecret = True
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Please type the twitch Client secret:")
    else:
        secret = secrets[0]
        t.set_client_secret(secret)
        configure_notif_workers()

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client secret configured")
        return


def generic_handle(update, context):
    global isConfigClientId, isConfigClientSecret, isAddingUser

    msgText = update.message.text
    if msgText.startswith("/") or msgText == "" or msgText is None:
        return

    if isConfigClientId and is_allowed(update, context):
        isConfigClientId = False
        t.set_client_id(msgText)

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client Id configured")
        return

    if isConfigClientSecret and is_allowed(update, context):
        isConfigClientSecret = False
        t.set_client_secret(msgText)

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client secret configured")
        return

    if isAddingUser and is_allowed(update, context):
        isAddingUser = False
        res = t.add_user(msgText, update.effective_chat.id,
                         update.effective_chat.type != 'private')

        notify_to_master(update, context, "adduser", msgText)

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
    reset_flags()
    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Done!, everything is canceled:")


def handle_error(update, context):
    """Log Errors caused by Updates."""
    logger.warning('Update "%s" caused error "%s"', update, context.error)


def handle_help(update, context):
    global commands
    helpMsg = u"Available commands: \n\n"
    for command in commands:
        if command["inHelp"]:
            helpMsg = helpMsg + \
                u"- /{0}: {1}\n".format(command["command"], command["info"])

    context.bot.send_message(
        chat_id=update.effective_chat.id, text=helpMsg)


def handle_twitch_add_user(update, context):
    global isConfigClientId, isConfigClientSecret, isAddingUser, telegram_masterchat
    isConfigClientSecret = False
    isConfigClientId = False
    isAddingUser = False
    users = get_param_value(update, "/adduser")
    if len(users) == 0:
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
            chat_id=update.effective_chat.id, text="Now I'll notify you in this chat when @{0} {1} streaming".format(", @".join(users), isOrAre))

        notify_to_master(update, context, "adduser", users)


def handle_add_admin(update, context):
    if not is_allowed(update, context):
        return

    admins = get_param_value(update, "/addadmin")
    if len(admins) == 0:
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="The name of admin(s) is required. use: /addadmin <telegram_username>")

    new_added = []
    for admin in admins:
        if not admin.startswith("@"):
            admin = "@{0}".format(admin)

        if admin in (user.lower() for user in telegram_whiteList):
            continue
        telegram_whiteList.append(admin)
        new_added.append(admin)

    if len(new_added) > 0:
        updateData("telegram_whiteList", telegram_whiteList)
        notify_to_master(update, context, "addadmin", new_added)


def handle_remove_admin(update, context):
    if not is_allowed(update, context):
        return

    admins = get_param_value(update, "/removeadmin")
    if len(admins) == 0:
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="The name of admin(s) is required. use: /removeadmin <telegram_username>")

    admin_removed = []
    for admin in admins:
        if not admin.startswith("@"):
            admin = "@{0}".format(admin)

        if not admin in (user.lower() for user in telegram_whiteList):
            continue
        telegram_whiteList.remove(admin)
        admin_removed.append(admin)

    if len(admin_removed) > 0:
        updateData("telegram_whiteList", telegram_whiteList)
        notify_to_master(update, context, "addadmin", admin_removed)


def handle_set_chat_master(update, context):
    if not is_allowed(update, context):
        return
    global telegram_masterchat

    if telegram_masterchat == update.effective_chat.id:
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Noting to chance")
        return

    msg = "Now I'll notify you here if any weird happens"
    try:
        updateData("telegram_masterchat", telegram_masterchat)
        notify_to_master(update, context, "setmasterchat")
        telegram_masterchat = update.effective_chat.id
    except:
        msg = "an exception occurred"

    context.bot.send_message(
        chat_id=update.effective_chat.id, text=msg)


def notify_to_master(update, context, cmd, value=None):
    global telegram_masterchat
    msg = "paso algo, pero no se que"
    author = update.effective_user.username
    if not value is None:
        value = str(value)

    if cmd == "start":
        msg = "{0} tried to start me".format(author)
    elif cmd == "setclientid":
        msg = "{0} tried to set the client id".format(author)
    elif cmd == "setsecretid":
        msg = "{0} tried to set the client secret".format(author)
    elif cmd == "adduser":
        msg = "{0} added the user {1} to {2}({3}) with name {4}".format(
            author, value, update.effective_chat.type, update.effective_chat.id, update.effective_chat.name)
    elif cmd == "addadmin":
        msg = "{0} added {1} as new admin".format(
            author, value)
    elif cmd == "setmasterchat":
        msg = "{0} changed the master chat to {1}({2}) named {3}".format(
            author, update.effective_chat.type, update.effective_chat.id, update.effective_chat.name)

    context.bot.send_message(
        chat_id=telegram_masterchat, text=msg)


def __init__():
    global t, commands, telegram_whiteList, telegram_masterchat, updater
    t = Twitch()
    config = configparser.ConfigParser()

    logger.info("reading bot configuration")
    config.read(DATA_FILE_PATH)
    telegram_botToken = config.get("telegram", "botToken")
    telegram_masterchat = int(config.get("telegram", "masterchat"))
    telegram_whiteList = eval(config.get("telegram", "whiteList"))

    logger.info("creating bot updater")
    updater = Updater(telegram_botToken, use_context=True)
    dp = updater.dispatcher

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
        },
        {
            "command": "addadmin",
            "handle": handle_add_admin,
            "info": "Adds a new telegram's user to whitelist (permissions). Use users separated by space to add multiple.",
            "inHelp": False
        },
        {
            "command": "removeadmin",
            "handle": handle_remove_admin,
            "info": "Removes a telegram's user from whitelist (permissions). Use users separated by space to add multiple.",
            "inHelp": False
        },
        {
            "command": "setmasterchat",
            "handle": handle_set_chat_master,
            "info": "Marks this chat as master to recive notifications (debugging purpose).",
            "inHelp": False
        }
    ]

    logger.info("configuring handlers")

    for i, v in enumerate(commands):
        dp.add_handler(CommandHandler(v["command"], v["handle"]))

    dp.add_error_handler(handle_error)
    dp.add_handler(MessageHandler(Filters.text, generic_handle))

    logger.info("configuring workers")
    configure_notif_workers()

    logger.info("starting polling")
    updater.start_polling()

    logger.info("Bot ready...")
    updater.idle()
