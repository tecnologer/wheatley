import configparser
import logging
from twitch import Twitch
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters
from telegram import InlineKeyboardButton, InlineKeyboardMarkup
from commands import Commands
from incomplete_command import IncompleteCommands

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

logger = logging.getLogger(__name__)

# constans
DATA_FILE_PATH = './data.ini'
NOTIFICATION_DELAY = 60

commands = Commands()
incopleteCmd = IncompleteCommands()

isNotifWorkerRunning = False

telegram_whiteList = []
telegram_masterchat = 0

telegram_botToken = ""
updater = None
t = None

commandsList = None


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
    if update.effective_user.username is None:
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Permissions required")
        return False
    user = update.effective_user.username.lower()
    if not user.startswith("@"):
        user = "@{0}".format(user)

    if not user in (user.lower() for user in telegram_whiteList):
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Permissions required")
        return False

    return True


def get_param_value(update, command):
    '''returns a collection with the values after the command'''
    if not command.startswith("/"):
        command = "/{0}".format(command)

    text = update.message.text if update.edited_message is None else update.edited_message.text
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
        notify_to_master(update, context, commands.start)
        return
    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Hello, I'll assist you to configure everything")

    notify_to_master(update, context, commands.start, flag=0)

    if check_missing_data(update, context):
        return

    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Congratulations!, everything is configured.")


def handle_twitch_client_id(update, context):
    if not is_allowed(update, context):
        notify_to_master(update, context, commands.set_client_id)
        return

    incopleteCmd.mark_completed(update)
    clientIds = get_param_value(update, commands.set_client_id)

    if len(clientIds) == 0:
        incopleteCmd.add(update, commands.set_client_id)
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
        notify_to_master(update, context, commands.set_client_id, flag=0)


def handle_twitch_client_secret(update, context):
    if not is_allowed(update, context):
        notify_to_master(update, context, commands.set_secret_id)
        return

    secrets = get_param_value(update, commands.set_secret_id)

    if len(secrets) == 0:
        incopleteCmd.add(update, commands.set_client_id)
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Please type the twitch Client secret:")
    else:
        secret = secrets[0]
        t.set_client_secret(secret)
        configure_notif_workers()

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client secret is configured")
        return


def generic_handle(update, context):
    msgText = update.message.text if update.edited_message is None else update.edited_message.text

    if msgText.startswith("/") or msgText == "" or msgText is None:
        return

    need_check = False

    if incopleteCmd.is_incomplete(update, commands.set_client_id) and is_allowed(update, context):
        t.set_client_id(msgText)
        need_check = True
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client Id is configured")
        notify_to_master(update, context, commands.set_client_id, flag=0)
    elif incopleteCmd.is_incomplete(update, commands.set_secret_id) and is_allowed(update, context):
        t.set_client_secret(msgText)

        need_check = True

        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Client secret is configured")
        notify_to_master(update, context, commands.set_client_id, flag=0)
    elif incopleteCmd.is_incomplete(update, commands.add_user):
        res = t.add_user(msgText, update.effective_chat.id,
                         update.effective_chat.type != update.effective_chat.PRIVATE)

        if res is not None:
            context.bot.send_message(
                chat_id=update.effective_chat.id, text=res)
        else:
            need_check = True
            context.bot.send_message(chat_id=update.effective_chat.id,
                                     text="Now I'll notify you in this chat when @{0} is streaming".format(msgText))
            notify_to_master(update, context, commands.add_user, msgText)

    elif incopleteCmd.is_incomplete(update, commands.remove_user):
        res = t.remove_user(msgText, update.effective_chat.id)

        notify_to_master(update, context, commands.remove_user, msgText)

        msg = "The notifications for @{0} are turned off".format(
            res.username) if res is not None else "Users is not configured"

        context.bot.send_message(
            chat_id=update.effective_chat.id, text=msg)

    incopleteCmd.mark_completed(update)

    if need_check:
        check_missing_data(update, context)


def handle_cancel(update, context):
    incopleteCmd.clear()
    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Done!, everything is canceled:")


def handle_error(update, context):
    """Log Errors caused by Updates."""
    logger.warning('Update "%s" caused error "%s"', update, context.error)


def handle_help(update, context):
    global commandsList
    helpMsg = u"Available commands: \n\n"
    for command in commandsList:
        if command["inHelp"]:
            helpMsg = helpMsg + \
                u"- /{0}: {1}\n".format(command["command"], command["info"])

    context.bot.send_message(
        chat_id=update.effective_chat.id, text=helpMsg)


def handle_twitch_add_user(update, context):
    users = get_param_value(update, commands.add_user)

    incopleteCmd.mark_completed(update)
    if len(t.unique_users_collection) >= 100:
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="The limit of users has been reached, please contact any of this administrators: {0}".format(", ".join(telegram_whiteList)))
        return

    if len(users) == 0:
        incopleteCmd.add(update, commands.add_user)
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Type the username of User's twitch:")
    else:
        added_users = []
        for user in users:
            res = t.add_user(user, update.effective_chat.id,
                             update.effective_chat.type != update.effective_chat.PRIVATE)
            if res is not None:
                context.bot.send_message(
                    chat_id=update.effective_chat.id, text=res)
            else:
                added_users.append(user)

        if len(added_users) > 0:
            isOrAre = "is" if len(users) == 1 else "are"
            context.bot.send_message(
                chat_id=update.effective_chat.id, text="Now I'll notify you in this chat when @{0} {1} streaming".format(", @".join(added_users), isOrAre))

            notify_to_master(update, context, commands.add_user, users)


def handle_twitch_remove_user(update, context):
    incopleteCmd.mark_completed(update)
    users = get_param_value(update, commands.remove_user)
    if len(users) == 0:
        incopleteCmd.add(update, commands.remove_user)
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="Type the username of User's twitch:")
    else:
        removed_users = []
        invalid_users = []
        for user in users:
            res = t.remove_user(user, update.effective_chat.id)
            if not res is None:
                removed_users.append(user)
            else:
                invalid_users.append(user)

        msg = "The notifications for @{0} are turned off".format(
            ", @".join(removed_users)) if len(removed_users) > 0 else "This users {0} are not configured".format(invalid_users)
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=msg)

        if len(removed_users) > 0:
            notify_to_master(update, context, commands.remove_user, users)


def handle_add_admin(update, context):
    if not is_allowed(update, context):
        return

    admins = get_param_value(update, commands.add_admin)
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
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="{0} is admin 😉".format(", ".join(new_added)))
        notify_to_master(update, context, commands.add_admin, new_added)


def handle_remove_admin(update, context):
    if not is_allowed(update, context):
        return

    admins = get_param_value(update, commands.remove_admin)
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
        context.bot.send_message(
            chat_id=update.effective_chat.id, text="{0} are no longer admin 😢".format(", ".join(admin_removed)))
        notify_to_master(update, context, commands.remove_admin, admin_removed)


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
        notify_to_master(update, context, commands.set_master_chat)
        telegram_masterchat = update.effective_chat.id
        updateData("telegram_masterchat", telegram_masterchat)
    except:
        msg = "an exception occurred"

    context.bot.send_message(
        chat_id=update.effective_chat.id, text=msg)


def notify_to_master(update, context, cmd, value=None, flag=None):
    global telegram_masterchat
    if telegram_masterchat == 0 or telegram_masterchat is None or telegram_masterchat == update.effective_chat.id:
        return
    msg = None
    author = update.effective_user.username if not update.effective_user.username is None else update.effective_user.first_name
    if author is None:
        author = update.effective_user.id

    if not value is None:
        value = str(value)

    chat_title = author if update.effective_chat.title is None else update.effective_chat.title

    if cmd == commands.start and flag is None:
        msg = "{0} tried to start me".format(author)
    elif cmd == commands.start and flag == 0:
        msg = "{0} started me".format(author)
    elif cmd == commands.set_client_id:
        msg = "{0} tried to set the client id".format(author)
    elif cmd == commands.set_client_id and flag == 0:
        msg = "{0} updated the client id".format(author)
    elif cmd == commands.set_secret_id and flag == 0:
        msg = "{0} updated the client secret".format(author)
    elif cmd == commands.add_user:
        msg = '{0} added the user {1} to {2}({3}) with name "{4}"'.format(
            author, value, update.effective_chat.type, update.effective_chat.id, chat_title)
    elif cmd == commands.add_admin:
        msg = "{0} added {1} as new admin".format(
            author, value)
    elif cmd == commands.set_master_chat:
        msg = '{0} changed the master chat to {1}({2}) named "{3}"'.format(
            author, update.effective_chat.type, update.effective_chat.id, chat_title)
    elif cmd == commands.remove_user:
        msg = '{0} removed the user {1} from {2}({3}) with name "{4}"'.format(
            author, value, update.effective_chat.type, update.effective_chat.id, chat_title)
    elif cmd == commands.remove_admin:
        msg = "{0} removed {1} as admin".format(
            author, value)
    else:
        msg = "paso algo, pero no se que.\n\n **User:** {0}\n **Chat:** {1}\n, **Command:** {2}\n, **Value:** {3}\n, **Flag:** {4}\n".format(
            author, chat_title, cmd, value, flag)
    context.bot.send_message(
        chat_id=telegram_masterchat, text=msg, parse_mode='MarkDown')


def handle_get_users(update, context):
    chat_id = update.effective_chat.id
    users = t.get_users_by_chat(chat_id)

    if len(users) == 0:
        context.bot.send_message(
            chat_id=chat_id, text="There are not users twitch configured in this chat")
        return

    msg = "*The following users twitch are registered in this chat:*\n\n"
    invalid_users = []
    for user in users:
        if user.twitch_id is None:
            invalid_users.append(user)
            continue
        msg = "{0}• twitch.tv/{1}\n".format(msg, user.username)

    if len(invalid_users) == len(users):
        msg = "*The following users twitch are registered in this chat:*\n\n"

    if len(invalid_users) > 0:
        msg = "{0}\n*The following users twitch are invalid:*\n\n".format(msg)

        for invalid in invalid_users:
            msg = "{0}• {1}\n".format(msg, invalid.username)

    context.bot.send_message(
        chat_id=chat_id, text=msg, parse_mode='MarkDown')


def run():
    global t, commandsList, telegram_whiteList, telegram_masterchat, updater
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

    commandsList = [
        {
            "command": commands.start,
            "handle": start,
            "info": "Starts the bot",
            "inHelp": False
        },
        {
            "command": commands.set_client_id,
            "handle": handle_twitch_client_id,
            "info": "Stores the twitch client id",
            "inHelp": False
        },
        {
            "command": commands.set_secret_id,
            "handle": handle_twitch_client_secret,
            "info": "Stores the twitch client secret",
            "inHelp": False
        },
        {
            "command": commands.help,
            "handle": handle_help,
            "info": "Shows the basic info of commands",
            "inHelp": False
        },
        {
            "command": commands.cancel,
            "handle": handle_cancel,
            "info": "Cancels the active command",
            "inHelp": True
        },
        {
            "command": commands.add_user,
            "handle": handle_twitch_add_user,
            "info": "Adds a new user(s) to the list to monitor its stream status. Use users separated by space to add multiple.",
            "inHelp": True
        },
        {
            "command": commands.remove_user,
            "handle": handle_twitch_remove_user,
            "info": "Removes a user(s) from the list to monitor its stream status. Use users separated by space to add multiple.",
            "inHelp": True
        },
        {
            "command": commands.add_admin,
            "handle": handle_add_admin,
            "info": "Adds a new telegram's user to whitelist (permissions). Use users separated by space to add multiple.",
            "inHelp": False
        },
        {
            "command": commands.remove_admin,
            "handle": handle_remove_admin,
            "info": "Removes a telegram's user from whitelist (permissions). Use users separated by space to add multiple.",
            "inHelp": False
        },
        {
            "command": commands.set_master_chat,
            "handle": handle_set_chat_master,
            "info": "Marks this chat as master to recive notifications (debugging purpose).",
            "inHelp": False
        },
        {
            "command": commands.get_users,
            "handle": handle_get_users,
            "info": "Returns the list of users twitch in the chat",
            "inHelp": True
        }
    ]

    logger.info("configuring handlers")

    for i, v in enumerate(commandsList):
        dp.add_handler(CommandHandler(v["command"], v["handle"]))

    dp.add_error_handler(handle_error)
    dp.add_handler(MessageHandler(Filters.text, generic_handle))

    logger.info("configuring workers")
    configure_notif_workers()

    logger.info("starting polling")
    updater.start_polling()

    logger.info("Bot ready...")
    updater.idle()
