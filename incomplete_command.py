import hashlib


class IncompleteCommands:
    __commands = {}

    def __init__(self):
        super().__init__()

    def add(self, update, command):
        author = self.get_author(update)
        if author in self.commands:
            return "Please complete {0} first or user /cancel".format(self.__commands[author])

        self.commands[author] = command
        return None

    def mark_completed(self, update, command):
        author = self.get_author(update)
        if author in self.commands:
            return "Nothing pending"
        del self.commands[author]
        return None

    def get_author(self, update):
        userid = update.effective_user.id
        chatid = update.effective_chat.id
        author = b"{0}-{1}".format(userid, chatid)
        return hashlib.md5(str.encode(author)).hexdigest()

    def user_has_incomplete_command(self, update):
        author = self.get_author(update)
        return author in self.__commands

    def clear(self):
        self.__commands = {}
