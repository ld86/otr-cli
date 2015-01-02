import argparse
import socket
import select
import sys


def parse_argv():
    parser = argparse.ArgumentParser()
    parser.add_argument('host', help='host')
    args = parser.parse_args()
    if args.host is None:
        parser.print_help()
    return args.host


def parse_host(host):
    host = host.split(':')
    if host == '':
        host = '0.0.0.0'
    return (host[0], int(host[1]))


def try_to_bind(host):
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.bind(parse_host(host))
    except socket.error:
        return None
    return s


def try_to_connect(host):
    try:
        return socket.create_connection(parse_host(host), 1)
    except socket.error:
        return None


class PlainText:

    def write(self, message):
        self.s.send(message)

    def read(self):
        return self.s.recv(1024)


class Server:

    def __init__(self, s):
        self.s = s

    def serve(self, handler):
        self.s.listen(1)
        while True:
            c, a = self.s.accept()
            handler(Client(c))


STDIN = 0
SOCKET = 1


class Client(PlainText):

    def __init__(self, s):
        self.s = s

    def select(self):
        r, _, _ = select.select([sys.stdin, self.s], [], [])

        if r[0] == sys.stdin:
            return (sys.stdin, STDIN)

        if r[0] == self.s:
            return (self.s, SOCKET)


class RawTerminal:

    def read(self):
        return raw_input()

    def write(self, message):
        print(message)


def chat(client):
    if client is None:
        raise 'Can\'t chat with None'

    while True:
        terminal = RawTerminal()
        descriptor, who = client.select()

        if who == STDIN:
            client.write(terminal.read())
        if who == SOCKET:
            terminal.write(client.read())


def main():
    host = parse_argv()

    s = try_to_bind(host)
    if s is not None:
        Server(s).serve(chat)
    else:
        s = try_to_connect(host)
        chat(Client(s))

if __name__ == '__main__':
    main()
