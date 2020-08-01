FROM python:3
COPY . /
RUN pip install telegram pickledb
RUN pip install python-telegram-bot --upgrade
CMD ["python", "/Bot.py"]