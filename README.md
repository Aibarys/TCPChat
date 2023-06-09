# Net-cat

```sh
Welcome to TCP-Chat!

         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'
```

## Description

This project consists on recreating the NetCat in a Server-Client Architecture that can run in a server mode on a specified port listening for incoming connections, and it can be used in a client mode, trying to connect to a specified port and transmitting information to the server. NetCat, nc system command, is a command-line utility that reads and writes data across network connections using TCP or UDP.

## Usage: how to run

- Server

```bash
  go run .
  go run . $port
```

- Client (in new Terminal)

```bash
  nc localhost $port
```

## Features

To change your name you can use the following command:
```bash
  [2023-06-07 21:08:48][Aiba]:/change
  [ENTER YOUR NEW NAME]:Aibarys
  You have changed the name from Aiba to Aibarys
  [2023-06-07 21:08:52][Aibarys]:
```
