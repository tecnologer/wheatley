import hashlib
import logging

logging.basicConfig(filename='ouput.log', format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)

logger = logging.getLogger(__name__)


class IncompleteCommands:
    __commands = {}

    def __init__(self):
        super().__init__()

    def add(self, update, command):
        author = self.get_author(update)
        self.__cancel_other_commands(author)

        self.__commands[author] = command
        return None

    def mark_completed(self, update):
        author = self.get_author(update)
        self.__cancel_other_commands(author)

    def get_author(self, update):
        userid = update.effective_user.id
        chatid = update.effective_chat.id
        author = "{0}-{1}".format(userid, chatid)
        return hashlib.md5(str.encode(author)).hexdigest()

    def user_has_incomplete_command(self, update):
        author = self.get_author(update)
        return author in self.__commands

    def is_incomplete(self, update, command):
        author = self.get_author(update)
        return author in self.__commands and command == self.__commands[author]

    def clear(self):
        logger.info(
            "Clear imcomplete commands. {0}".format(self.__commands))
        self.__commands = {}

    def __cancel_other_commands(self, author):
        if not author in self.__commands:
            logger.info(
                "Clear empty imcomplete")
            return
        logger.info(
            "Clear imcomplete for user. {0}".format(self.__commands[author]))
        del self.__commands[author]
