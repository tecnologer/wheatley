FROM python:3
COPY . /
# RUN ls -la
RUN python -m pip install --upgrade pip
RUN pip install telegram pickledb requests
RUN pip install python-telegram-bot --upgrade
CMD ["python", "Bot.py"]