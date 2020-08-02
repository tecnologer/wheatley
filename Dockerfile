FROM python:3
COPY . /
RUN pip3 install telegram pickledb requests
RUN pip3 install python-telegram-bot --upgrade
CMD ["python", "/Bot.py"]