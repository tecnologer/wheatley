import configparser
import logging
import sys
from twitch import Twitch
from telegram.ext import Updater, CommandHandler, MessageHandler, Filters
from telegram import InlineKeyboardButton, InlineKeyboardMarkup
from commands import Commands
from incomplete_command import IncompleteCommands
from language import Language

logging.basicConfig(filename='/tmp/twitch_bot_output.log', format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
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


def is_master_chat(update):
    return telegram_masterchat == update.effective_chat.id


def get_message_from_update(update):
    if update.channel_post is not None:
        msgText = update.channel_post.text
    else:
        msgText = update.message.text if update.edited_message is None else update.edited_message.text

    return msgText


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
    lang = Language(update)
    if update.effective_user.username is None:
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("is_allowed_fail"))
        return False
    user = update.effective_user.username.lower()
    if not user.startswith("@"):
        user = "@{0}".format(user)

    if not user in (user.lower() for user in telegram_whiteList):
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("is_allowed_fail"))
        return False

    return True


def get_param_value(update, command):
    '''returns a collection with the values after the command'''
    if not command.startswith("/"):
        command = "/{0}".format(command)

    lang = Language(update)
    text = get_message_from_update(update)
    if not text.startswith(command):
        return []

    botName = update.effective_user.bot.username.lower(
    ) if update.effective_user is not None else ""

    text = text.lower().replace(command.lower(), "").replace(
        botName, "").replace("@", "", -1)
    text = text.strip()
    if text == "":
        return []
    return text.split(" ")


def check_missing_data(update, context):
    lang = Language(update)
    if not t.is_client_id_set():
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("config_assist_client_id"))
        handle_twitch_client_id(update, context)
        return True

    if not t.is_client_secret_set():
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("config_assist_client_secret"))
        handle_twitch_client_secret(update, context)
        return True

    return False


def start(update, context):
    if not is_allowed(update, context):
        notify_to_master(update, context, commands.start.cmd)
        return

    lang = Language(update)
    context.bot.send_message(
        chat_id=update.effective_chat.id, text=lang.get("config_assist_greeting"))

    notify_to_master(update, context, commands.start.cmd, flag=0)

    if check_missing_data(update, context):
        return

    context.bot.send_message(
        chat_id=update.effective_chat.id, text=lang.get("config_assist_configured"))


def handle_twitch_client_id(update, context):
    if not is_allowed(update, context):
        notify_to_master(update, context, commands.set_client_id.cmd)
        return
    lang = Language(update)
    incopleteCmd.mark_completed(update)
    clientIds = get_param_value(update, commands.set_client_id.cmd)

    if len(clientIds) == 0:
        incopleteCmd.add(update, commands.set_client_id.cmd)
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("client_id_request"))
    else:
        clientId = clientIds[0]
        t.set_client_id(clientId)
        configure_notif_workers()

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("client_id_success"))
        notify_to_master(update, context, commands.set_client_id.cmd, flag=0)


def handle_twitch_client_secret(update, context):
    if not is_allowed(update, context):
        notify_to_master(update, context, commands.set_secret_id.cmd)
        return
    lang = Language(update)
    secrets = get_param_value(update, commands.set_secret_id.cmd)

    if len(secrets) == 0:
        incopleteCmd.add(update, commands.set_client_id.cmd)
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("client_secret_request"))
    else:
        secret = secrets[0]
        t.set_client_secret(secret)
        configure_notif_workers()

        if check_missing_data(update, context):
            return

        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("client_secret_success"))
        return


def isCommandFromChannel(update, context):
    postTxt = update.channel_post.text
    if postTxt == "":
        return False

    cmd = postTxt.split(" ")[0]
    if cmd == "" or len(cmd) == 1:
        return False

    cmd = cmd[1:]
    for command in commandsList:
        if command.cmd == cmd:
            command.handle(update, context)
            return True

    return False


def generic_handle(update, context):
    if update.channel_post is not None:
        if isCommandFromChannel(update, context):
            return

    msgText = get_message_from_update(update)

    if msgText.startswith("/") or msgText == "" or msgText is None:
        return

    lang = Language(update)
    need_check = False

    if incopleteCmd.is_incomplete(update, commands.set_client_id.cmd) and is_allowed(update, context):
        t.set_client_id(msgText)
        need_check = True
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("client_id_success"))
        notify_to_master(update, context, commands.set_client_id.cmd, flag=0)
    elif incopleteCmd.is_incomplete(update, commands.set_secret_id.cmd) and is_allowed(update, context):
        t.set_client_secret(msgText)

        need_check = True

        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("client_secret_success"))
        notify_to_master(update, context, commands.set_client_id.cmd, flag=0)
    elif incopleteCmd.is_incomplete(update, commands.add_user.cmd):
        res = t.add_user(msgText, update.effective_chat.id,
                         update.effective_chat.type != update.effective_chat.PRIVATE)

        if res is not None:
            context.bot.send_message(
                chat_id=update.effective_chat.id, text=res)
        else:
            need_check = True
            context.bot.send_message(chat_id=update.effective_chat.id,
                                     text=lang.get("add_user_success").format(msgText))
            notify_to_master(update, context, commands.add_user.cmd, msgText)

    elif incopleteCmd.is_incomplete(update, commands.remove_user.cmd):
        res = t.remove_user(msgText, update.effective_chat.id)

        notify_to_master(update, context, commands.remove_user.cmd, msgText)

        msg = lang.get("remove_user_success").format(
            res.username) if res is not None else lang.get("remove_user_not_found")

        context.bot.send_message(
            chat_id=update.effective_chat.id, text=msg)

    incopleteCmd.mark_completed(update)

    if need_check:
        check_missing_data(update, context)


def handle_cancel(update, context):
    incopleteCmd.mark_completed(update)
    lang = Language(update)
    context.bot.send_message(
        chat_id=update.effective_chat.id, text=lang.get("cancel_success"))


def handle_error(update, context):
    """Log Errors caused by Updates."""
    logger.warning('Update "%s" caused error "%s"', update, context.error)


def handle_help(update, context):
    global commandsList
    lang = Language(update)

    logger.info("/{0}".format(commands.help.cmd))

    helpMsg = u"**{0}**\n\n{1}\n\n".format(
        lang.get("bot_description"), lang.get("help_header"))
    for command in commandsList:
        if command.show_in_help or is_master_chat(update):
            helpMsg = helpMsg + \
                u"- /{0}: {1}\n".format(command.cmd,
                                        command.get_info(lang.code))

    context.bot.send_message(
        chat_id=update.effective_chat.id, text=helpMsg, parse_mode='MarkDown')


def handle_twitch_add_user(update, context):
    users = get_param_value(update, commands.add_user.cmd)
    chat_id = update.effective_chat.id

    lang = Language(update)
    logger.info(
        "/{0} {1} in chat {2}".format(commands.add_user.cmd, " ".join(users), chat_id))

    incopleteCmd.mark_completed(update)
    if len(t.unique_users_collection) >= 100:
        context.bot.send_message(
            chat_id=chat_id, text=lang.get("add_user_limit_reached").format(", ".join(telegram_whiteList)))
        return

    if len(users) == 0:
        incopleteCmd.add(update, commands.add_user.cmd)
        context.bot.send_message(
            chat_id=chat_id, text=lang.get("add_user_request_name"))
    else:
        added_users = []
        for user in users:
            res = t.add_user(user, chat_id,
                             update.effective_chat.type != update.effective_chat.PRIVATE)
            if res is not None:
                context.bot.send_message(
                    chat_id=chat_id, text=res)
            else:
                added_users.append(user)

        if len(added_users) > 0:
            context.bot.send_message(
                chat_id=chat_id, text=lang.get("add_user_success").format(", @".join(added_users)))

            notify_to_master(update, context, commands.add_user.cmd, users)


def handle_twitch_remove_user(update, context):
    incopleteCmd.mark_completed(update)
    users = get_param_value(update, commands.remove_user.cmd)
    chat_id = update.effective_chat.id
    logger.info(
        "/{0} {1} in chat: {2}".format(commands.remove_user.cmd, " ".join(users), chat_id))

    lang = Language(update)
    if len(users) == 0:
        incopleteCmd.add(update, commands.remove_user.cmd)
        context.bot.send_message(
            chat_id=chat_id, text=lang.get("remove_user_request_name"))
    else:
        removed_users = []
        invalid_users = []
        for user in users:
            res = t.remove_user(user, chat_id)
            if not res is None:
                removed_users.append(user)
            else:
                invalid_users.append(user)

        msg = lang.get("remove_user_success").format(
            ", @".join(removed_users)) if len(removed_users) > 0 else lang.get("remove_user_not_found").format(invalid_users)
        context.bot.send_message(
            chat_id=chat_id, text=msg)

        if len(removed_users) > 0:
            notify_to_master(update, context, commands.remove_user.cmd, users)


def handle_add_admin(update, context):
    if not is_allowed(update, context):
        return

    lang = Language(update)
    admins = get_param_value(update, commands.add_admin.cmd)
    chat_id = update.effective_chat.id
    logger.info(
        "/{0} {1} in chat: {2}".format(commands.add_admin.cmd, " ".join(admins), chat_id))

    if len(admins) == 0:
        context.bot.send_message(
            chat_id=chat_id, text=lang.get("add_admin_no_specified").format(commands.add_admin.cmd))

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
            chat_id=chat_id, text=lang.get("add_admin_success").format(", ".join(new_added)))
        notify_to_master(update, context, commands.add_admin.cmd, new_added)


def handle_remove_admin(update, context):
    if not is_allowed(update, context):
        return

    lang = Language(update)
    admins = get_param_value(update, commands.remove_admin.cmd)

    chat_id = update.effective_chat.id
    logger.info(
        "/{0} {1} in chat: {2}".format(commands.add_admin.cmd, " ".join(admins), chat_id))
    if len(admins) == 0:
        context.bot.send_message(
            chat_id=chat_id, text=lang.get("remove_admin_no_specified").format(commands.remove_admin.cmd))

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
            chat_id=chat_id, text=lang.get("remove_admin_success").format(", ".join(admin_removed)))
        notify_to_master(
            update, context, commands.remove_admin.cmd, admin_removed)


def handle_set_chat_master(update, context):
    if not is_allowed(update, context):
        return
    global telegram_masterchat

    lang = Language(update)
    if telegram_masterchat == update.effective_chat.id:
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("chat_master_no_change"))
        return

    msg = lang.get("chat_master_success")
    try:
        notify_to_master(update, context, commands.set_master_chat.cmd)
        telegram_masterchat = update.effective_chat.id
        updateData("telegram_masterchat", telegram_masterchat)
    except:
        msg = lang.get("chat_master_error")

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

    if cmd == commands.start.cmd and flag is None:
        msg = "{0} tried to start me".format(author)
    elif cmd == commands.start.cmd and flag == 0:
        msg = "{0} started me".format(author)
    elif cmd == commands.set_client_id.cmd:
        msg = "{0} tried to set the client id".format(author)
    elif cmd == commands.set_client_id.cmd and flag == 0:
        msg = "{0} updated the client id".format(author)
    elif cmd == commands.set_secret_id.cmd and flag == 0:
        msg = "{0} updated the client secret".format(author)
    elif cmd == commands.add_user.cmd:
        msg = '{0} added the user {1} to {2}({3}) with name "{4}"'.format(
            author, value, update.effective_chat.type, update.effective_chat.id, chat_title)
    elif cmd == commands.add_admin.cmd:
        msg = "{0} added {1} as new admin".format(
            author, value)
    elif cmd == commands.set_master_chat.cmd:
        msg = '{0} changed the master chat to {1}({2}) named "{3}"'.format(
            author, update.effective_chat.type, update.effective_chat.id, chat_title)
    elif cmd == commands.remove_user.cmd:
        msg = '{0} removed the user {1} from {2}({3}) with name "{4}"'.format(
            author, value, update.effective_chat.type, update.effective_chat.id, chat_title)
    elif cmd == commands.remove_admin.cmd:
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

    lang = Language(update)
    logger.info(
        "/{0}, users count: {1}".format(commands.get_users.cmd, len(users)))

    if len(users) == 0:
        context.bot.send_message(
            chat_id=chat_id, text=lang.get("get_users_no_users"))
        return

    msg = lang.get("get_users_header")
    invalid_users = []
    for user in users:
        if user.twitch_id is None:
            invalid_users.append(user)
            continue
        msg = "{0}• [{1}](twitch.tv/{1})\n".format(msg, user.username)

    if len(invalid_users) == len(users):
        msg = lang.get("get_users_header")

    if len(invalid_users) > 0:
        msg = lang.get("get_users_header_invalid_users").format(msg)

        for invalid in invalid_users:
            msg = "{0}• {1}\n".format(msg, invalid.username)

    context.bot.send_message(
        chat_id=chat_id, text=msg, parse_mode='MarkDown')


def handle_stream_status(update, context):
    users = get_param_value(update, commands.stream_status.cmd)
    lang = Language(update)
    if len(users) == 0:
        context.bot.send_message(
            chat_id=update.effective_chat.id, text=lang.get("stream_status_error").format(commands.stream_status.cmd))
        return

    t.get_stream_status(update, context, *users)


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
        commands.start.set_handle(start),
        commands.set_client_id.set_handle(handle_twitch_client_id),
        commands.set_secret_id.set_handle(handle_twitch_client_secret),
        commands.help.set_handle(handle_help),
        commands.cancel.set_handle(handle_cancel),
        commands.add_user.set_handle(handle_twitch_add_user),
        commands.remove_user.set_handle(handle_twitch_remove_user),
        commands.add_admin.set_handle(handle_add_admin),
        commands.remove_admin.set_handle(handle_remove_admin),
        commands.set_master_chat.set_handle(handle_set_chat_master),
        commands.get_users.set_handle(handle_get_users),
        commands.stream_status.set_handle(handle_stream_status),
    ]

    logger.info("configuring handlers")

    for command in commandsList:
        dp.add_handler(CommandHandler(command.cmd, command.handle))

    dp.add_error_handler(handle_error)
    dp.add_handler(MessageHandler(Filters.text, generic_handle))

    logger.info("configuring workers")
    configure_notif_workers()

    logger.info("starting polling")
    updater.start_polling()

    logger.info("Bot ready...")
    updater.idle()
