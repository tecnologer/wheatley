import dao

PREFIX_DB = "games"


def get_game(id):
    return dao.get(PREFIX_DB, id)


def add_game(id, name):
    dao.save(PREFIX_DB, id, name)
