import json
import datetime
import telegrambot
import configparser
import logging
from requests import post
from requests import get
from pytz import timezone


#from __future__ import unicode_literals
logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

logger = logging.getLogger(__name__)

isAlive = False
globalViewerCount = 0
dataFile = './data.ini'

# create a application on dev.twitch.com
twitch_clientId = ""
twitch_clientSecret = ""
twitch_userId = ""
twitch_userName = ""
twitch_token = ""
telegram_botToken = ""
telegram_whiteList = ""
telegram_chatId = 1


def error(update, context):
    """Log Errors caused by Updates."""
    logger.warning('Update "%s" caused error "%s"', update, context.error)


def updateData(key, value):
    global dataFile
    key = key.split("_")
    config = configparser.ConfigParser()
    config.read(dataFile)
    cfgfile = open(dataFile, 'w')
    config.set(key[0], key[1], str(value))
    config.write(cfgfile)
    cfgfile.close()


def readData():
    global dataFile
    global twitch_clientId
    global twitch_clientSecret
    global twitch_userId
    global twitch_userName
    global telegram_botToken
    global telegram_whiteList
    global telegram_chatId
    global twitch_token
    config = configparser.ConfigParser()
    config.read(dataFile)

    twitch_clientId = config.get("twitch", "clientId")
    twitch_clientSecret = config.get("twitch", "clientSecret")
    twitch_userId = config.get("twitch", "userId")
    twitch_userName = config.get("twitch", "userName")
    twitch_token = config.get("twitch", "token")
    telegram_botToken = config.get("telegram", "botToken")
    telegram_whiteList = config.get("telegram", "whiteList").split(",")
    telegram_chatId = config.get("telegram", "chatId")


def checkToken():
    global twitch_clientId
    global twitch_clientSecret
    global twitch_token

    url = 'https://id.twitch.tv/oauth2/validate'

    myobj = {'Authorization': 'OAuth {0}'.format(twitch_token)}
    x = get(url, headers=myobj)
    data = json.loads(x.text)

    expire = int(data['expires_in'])/86400

    if expire < 2:
        url = 'https://id.twitch.tv/oauth2/revoke'
        myobj = {'client_id': '{0}'.format(
            twitch_clientId), 'token': '{0}'.format(twitch_token)}

        x = post(url, data=myobj)

# use the part below to generate a new token, if you dont have a token to start with :)
        url = 'https://id.twitch.tv/oauth2/token'
        myobj = {'client_id': '{0}'.format(twitch_clientId), 'client_secret': '{0}'.format(
            twitch_clientSecret), 'grant_type': 'client_credentials'}
        x = post(url, data=myobj)

        data = json.loads(x.text)
#        print data['access_token']

        twitch_token = str(data['access_token'])
        updateData("twitch_token", twitch_token)

        return 'Token Renewed'
    else:
        return 'Token ok'


def streamstatus(update, context):
    global telegram_chatId
    if telegram_chatId == 1 or not isAllowedTelegramUser(update.effective_user.name):
        return
    global twitch_token
    global twitch_clientId

    myobj = {"client-id": "{0}".format(twitch_clientId),
             "Authorization": "Bearer {0}".format(twitch_token)}

    response = get(
        'https://api.twitch.tv/helix/streams?user_id={0}'.format(twitch_userId), headers=myobj)

#    response =  get('https://api.twitch.tv/helix/streams?user_id=51956085', headers={"client-id":"8yy2cqs86znp95v3bvht8ebgxvseh5"})
    data = json.loads(response.text)

    if len(data['data']) > 0:
        data2 = data['data'][0]
        userName = data2['user_name']
        streamTitle = data2['title']
        viewerCount = data2['viewer_count']
        gameID = data2['game_id']
        startedAt = data2['started_at']
        response = get(
            'https://api.twitch.tv/helix/games?id={0}'.format(gameID), headers=myobj)
        data = json.loads(response.text)
        gameName = data['data'][0]['name']

        context.bot.send_message(chat_id=telegram_chatId, text='{0} is live!! Streaming {1} with {2} viewers'.format(
            userName, gameName, viewerCount))
    else:
        context.bot.send_message(
            chat_id=telegram_chatId, text=u'{0} stream not online :('.format(twitch_userName))


def callback_minute(context):
    global telegram_chatId
    if telegram_chatId == 1:
        return
    global isAlive
    global globalViewerCount
    global twitch_clientId
    global twitch_token

    myobj = {"client-id": "{0}".format(twitch_clientId),
             "Authorization": "Bearer {0}".format(twitch_token)}
    response = get(
        'https://api.twitch.tv/helix/streams?user_id={0}'.format(twitch_userId), headers=myobj)

    data = json.loads(response.text)

    if len(data['data']) > 0:
        data2 = data['data'][0]
        userName = data2['user_name']
        streamTitle = data2['title']
        viewerCount = data2['viewer_count']
        gameID = data2['game_id']
        startedAt = data2['started_at']
        newValueIsAlive = data2["type"] == 'live'

        if newValueIsAlive and not isAlive:
            globalViewerCount = viewerCount
            response = get(
                'https://api.twitch.tv/helix/games?id={0}'.format(gameID), headers=myobj)
            data = json.loads(response.text)
            gameName = data['data'][0]['name']
            retval = context.bot.send_message(chat_id=telegram_chatId, text='[{0}](https://twitch.tv/{0}) is live!! Streaming {1} with {2} viewers'.format(
                userName, gameName, viewerCount), parse_mode='MarkDown')
            context.bot.pin_chat_message(
                chat_id=telegram_chatId, message_id=retval.message_id)

        isAlive = newValueIsAlive

#	else:
#	    if viewerCount != globalViewerCount:
#		msg = 'Viewercount RubenSaurus {0} naar {1}'.format('gedaald' if viewerCount< globalViewerCount else 'gestegen', viewerCount)
#		context.bot.send_message(chat_id=gotfChatId, text=msg)
    else:
        if isAlive == True:
            isAlive = False
            context.bot.send_message(
                chat_id=telegram_chatId, text='{0} stream is not running :('.format(twitch_userName))
            context.bot.unpin_chat_message(chat_id=telegram_chatId)
#    [like this](http://someurl)

#    context.bot.send_message(chat_id=gotfChatId, text='One message every minute')


def callback_daily(context):
    print(checkToken())  # replace Twitch token if neccesary


def countMessages(update, context):
    global telegram_chatId
    if not isAllowedTelegramUser(update.effective_user.name):
        return
    telegram_chatId = update.effective_chat.id
    chat_id = update.effective_chat.id
#    userid = update.message.from_user
    userid = update.effective_user.id
    username = update.effective_user.first_name
    bot = context.bot
    l_message = update.message.text


def showhelp(update, context):
    if not isAllowedTelegramUser(update.effective_user.name):
        return
    chat_id = update.effective_chat.id
    #    userid = update.message.from_user
    userid = update.effective_user.id
    username = update.effective_user.first_name
    bot = context.bot
    l_message = update.message.text


def start(update, context):
    global telegram_chatId
    if not isAllowedTelegramUser(update.effective_user.name):
        return
    telegram_chatId = update.effective_chat.id
    updateData("telegram_chatid", telegram_chatId)
    context.bot.send_message(
        chat_id=update.effective_chat.id, text="Ok, I'm ready, now I'll notify you in this chat")


def isAllowedTelegramUser(userName):
    return userName in telegram_whiteList


# def main():
#     """Start the bot."""
    # Create the Updater and pass it your bot's token.
    # Make sure to set use_context=True to use the new context based callbacks
    # Post version 12 this will no longer be necessary
    # updateData("test_key1", "hola")
    # readData()

    # updater = Updater(telegram_botToken,
    #                   use_context=True)  # put in Bot token
    # j = updater.job_queue
    # dp = updater.dispatcher

    # # on different commands - answer in Telegram
    # dp.add_handler(CommandHandler("start", start))
    # dp.add_handler(CommandHandler("help", showhelp))
    # dp.add_handler(CommandHandler("streamstatus", streamstatus))

    # dp.add_handler(MessageHandler(Filters.text, countMessages))

    # dp.add_error_handler(error)

    # # Start the Bot
    # updater.start_polling()
    # job_minute = j.run_repeating(callback_minute, interval=60, first=0)
    # job_daily = j.run_daily(callback_daily, datetime.time(
    #     0, 0, 0, tzinfo=timezone('Europe/Amsterdam')))

    # Run the bot until you press Ctrl-C or the process receives SIGINT,
    # SIGTERM or SIGABRT. This should be used most of the time, since
    # start_polling() is non-blocking and will stop the bot gracefully.

    # updater.idle()


if __name__ == '__main__':
    telegrambot.__init__()
