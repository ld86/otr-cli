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

    def handle(self, s):
        while True:
            r, _, _ = select.select([sys.stdin, s], [], [])

            if r[0] == sys.stdin:
                message = raw_input()
                s.send(message)

            if r[0] == s:
                message = s.recv(1024)
                if not message:
                    break
                print(message)


class Server(PlainText):

    def __init__(self, s):
        self.s = s

    def serve(self):
        self.s.listen(1)
        while True:
            c, a = self.s.accept()
            self.handle(c)


class Client(PlainText):

    def __init__(self, s):
        self.s = s

    def serve(self):
        self.handle(self.s)


def main():
    host = parse_argv()

    s = try_to_bind(host)
    if s is not None:
        Server(s).serve()

    s = try_to_connect(host)
    if s is not None:
        Client(s).serve()


if __name__ == '__main__':
    main()
