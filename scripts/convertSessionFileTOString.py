# Convert Session Files To StringSession

import asyncio
import json
import os
import sys
from os import path
import socks
import random
import traceback
from telethon.errors import SessionPasswordNeededError
from telethon.sessions import StringSession
from telethon.sync import TelegramClient

########################################################################

api_id = '1558897'
api_hash = '22cef8c309dc02618e7ed494f55cd211'
phonelist = []
index = 0

proxies = []


# get all files from sessions folder and split it by .session
def get_sessions():
    sessions = []
    for file in os.listdir('sessions'):
        if file.endswith('.session'):
            phoneNumber = file.split('.')[0]
            sessions.append(phoneNumber)
    return sessions


def get_proxies():
    proxies = []
    with open("proxy.csv", "r") as f:
        for line in f.readlines():
            line = line.strip()
            proxies.append(line)
    return proxies


class mani():
    def __init__(self):
        try:
            if not path.isdir("sessions"):
                raise Exception("sessions Not Exits Create One")
        except Exception as es:
            print(es)
            input("Press Any Key To Exit")
            sys.exit(1)
        global phonelist
        phonelist = get_sessions()
        global proxies
        proxies = get_proxies()

        if phonelist == []:
            print("No Sessions Found")
            input("Press Any Key To Exit")
            sys.exit(1)

        loop = asyncio.get_event_loop()

        loop.run_until_complete(self.login())

    async def login(self):
        for number in phonelist:
            try:
                p = random.choices(proxies)
                print(p)
                username, password, ip, port  = p[0].split(",")

                print(f"Trying To Login {number}")
                print(f"Using Proxy {ip}:{port}")
                current_number = os.getcwd() + '/sessions/' + number
                client = TelegramClient(current_number, api_id, api_hash,)
                                        # p=(socks.SOCKS5, ip, 5501,True, username, password))
                # await client.connect()
                # if not await client.is_user_authorized():
                #     print(f"Not Logined {number}")
                #     continue
                #     print("Not Login")
                #     await client.send_code_request(number)
                #     try:
                #         me = await client.sign_in(number, input('Enter code Sened In Telegram: '))
                #     except SessionPasswordNeededError:
                #         password = input("Enter password: ")
                #         me = await client.sign_in(password=password)
                sessionString = StringSession.save(client.session)

                global index
                index = index + 1

#                 phoneObj = {
#                     "index": index,
#                     "phoneNumber": number,
#                     "sessionString": sessionString,
#                     "p": True,
#                     "ip": ip,
#                     "port": int(port),
#                     "username": username,
#                     "password": password
#                 }
#
#                 phones = json.load(open("sessions.json"))
#                 phones.append(phoneObj)
#
#                 # append phoneObj object in phones  Array in json file
#                 with open('sessions.json', 'w') as f:
#                     json.dump(phones, f, indent=4)

                print(f"Phone Number: {number}")
                print(f"Session String: {sessionString}")
                # save sessionString in csv
                # with open("sessions.csv", "a", encoding='UTF-8') as f:
                #     writer = csv.writer(f, delimiter=",", lineterminator="\n")
                #     writer.writerow([number, sessionString])

                print(f"{number} Logined")
                await client.disconnect()
            except KeyboardInterrupt:
                client.disconnect()
                print("KeyBoard Interrupted")
                exit(1)
            except:
                traceback.print_exc()
            await asyncio.sleep(3)
        print("All Number Has Been Logined")


try:
    print("""
    Developed By - https://www.fiverr.com/honey_devv
    Telegram - t.me/ronnekeren | t.me/blackhat_dev
    """)
    mani()
except KeyboardInterrupt:
    print("KeyBoard Interrupted")
except Exception as es:
    print("Unexpected Error")
    print(es)
    input("Press Enter To Exit.")
    sys.exit(1)
input("Press Enter to continue")
