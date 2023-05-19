#define _GNU_SOURCE

#include "server.h"

#include <arpa/inet.h>
#include <netinet/in.h>
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/socket.h>
#include <sys/types.h>
#include <unistd.h>
#include <stdbool.h>

#define BUF_MAXSIZE 1024

struct userinfo {
  char id[100];
  int playwith;
};

int fdt[MAX_CLIENT_COUNT] = {0};
char mes[1024];
int SendToClient(int fd, char *buf, int Size);
struct userinfo users[100];
int find_fd(char *name);
int win_dis[8][3] = {{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6},
                     {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}};

void exit_with_response(const char response[]) {
  perror(response);
  exit(EXIT_FAILURE);
}

void exit_on_error(const int value, const char response[]) {
  if (value < 0) {
    exit_with_response(response);
  }
}

void message_all_user(char *chatting) {
  int i = 0;
  for (i = 0; i < 100; i += 1) {
    if (users[i].id[0] != '\0') {
      printf("%s", chatting);
      send(i, chatting, strlen(chatting), 0);
    }
  }
}

void message_handler(char *mes, int sender) {
  int instruction = 0;
  sscanf(mes, "%d", &instruction);
  switch (instruction) {
  case 9: {
    message_all_user(mes);
    break;
  }
  case 1: { // add user
    char name[100];
    sscanf(mes, "1 %s", name);
    strncpy(users[sender].id, name, 100);
    send(sender, "1", 1, 0);
    printf("1:%s\n", name);
    break;
  }
  case 2: { // show
    char buf[BUF_MAXSIZE], tmp[100];
    int p = sprintf(buf, "2 ");
    for (int i = 0; i < 100; i += 1) {
      if (strcmp(users[i].id, "") != 0) {
        sscanf(users[i].id, "%s", tmp);
        p = sprintf(buf + p, "%s ", tmp) + p;
      }
    }
    printf("2:%s\n", buf);
    send(sender, buf, strlen(buf), 0);
    break;
  }
  case 3: { // invite
    char a[100], b[100];
    char buf[BUF_MAXSIZE];
    sscanf(mes, "3 %s %s", a, b);
    int b_fd = find_fd(b);
    sprintf(buf, "4 %s invite you. Accept?\n", a);
    send(b_fd, buf, strlen(buf), 0);
    printf("3:%s", buf);
    break;
  }
  case 5: { // agree(1) or not(0)
    int state;
    char inviter[100];
    sscanf(mes, "5 %d %s", &state, inviter);
    if (state == 1) {
      send(sender, "6\n", 2, 0);
      send(find_fd(inviter), "6\n", 2, 0);
      int fd = find_fd(inviter);
      users[sender].playwith = fd;
      users[fd].playwith = sender;
      printf("6:\n");
    }
    break;
  }
  case 7: {
    int board[9];
    char state[100];
    char buf[BUF_MAXSIZE];
    sscanf(mes, "7  %d %d %d %d %d %d %d %d %d", &board[0], &board[1],
           &board[2], &board[3], &board[4], &board[5], &board[6], &board[7],
           &board[8]);
    for (int i = 0; i < 100; i += 1) {
      state[i] = '\0';
    }

    memset(buf, '\0', BUF_MAXSIZE);
    memset(state, '\0', sizeof(state));
    strcat(state, users[sender].id);
    for (int i = 0; i < 8; i += 1) {
      if (board[win_dis[i][0]] == board[win_dis[i][1]] &&
          board[win_dis[i][1]] == board[win_dis[i][2]]) {
        if (board[win_dis[i][0]] != 0) {
          strcat(state, "_Win!\n");
          sprintf(buf, "8  %d %d %d %d %d %d %d %d %d %s\n", board[0], board[1],
                  board[2], board[3], board[4], board[5], board[6], board[7],
                  board[8], state);
          printf("7:%s", buf);
          send(sender, buf, sizeof(buf), 0);
          send(users[sender].playwith, buf, sizeof(buf), 0);
          return;
        }
      }
    }
    memset(buf, '\0', BUF_MAXSIZE);
    memset(state, '\0', sizeof(state));
    for (int i = 0; i < 9; i += 1) {
      if (i == 8) {
        strcat(state, "Tie!\n");
        sprintf(buf, "8  %d %d %d %d %d %d %d %d %d %s\n", board[0], board[1],
                board[2], board[3], board[4], board[5], board[6], board[7],
                board[8], state);
        printf("7:%s", buf);
        send(sender, buf, sizeof(buf), 0);
        send(users[sender].playwith, buf, sizeof(buf), 0);
        return;
      }
      if (board[i] == 0) {
        break;
      }
    }

    memset(buf, '\0', BUF_MAXSIZE);
    memset(state, '\0', sizeof(state));
    strcat(state, users[users[sender].playwith].id);
    strcat(state, "_your_tern!\n");
    sprintf(buf, "8  %d %d %d %d %d %d %d %d %d %s\n", board[0], board[1],
            board[2], board[3], board[4], board[5], board[6], board[7],
            board[8], state);
    printf("7:%s", buf);
    send(sender, buf, sizeof(buf), 0);
    send(users[sender].playwith, buf, sizeof(buf), 0);
    break;
  }
  }
}

void *pthread_service(void *socket) {
  int socket_ = *(int *)socket;
  while (true) {
    int numbytes;
    int i;
    numbytes = recv(socket_, mes, BUF_MAXSIZE, 0);
    printf("\n\n\n%s\n\n\n", mes);

    // close socket
    if (numbytes <= 0) {
      for (i = 0; i < MAX_CLIENT_COUNT; i += 1) {
        if (socket_ == fdt[i]) {
          fdt[i] = 0;
        }
      }
      memset(users[socket_].id, '\0', sizeof(users[socket_].id));
      users[socket_].playwith = -1;
      break;
    }
    message_handler(mes, socket_);
    bzero(mes, BUF_MAXSIZE);
  }
  close(socket_);
}

int main() {
  int server_socket;
  struct sockaddr_in server;
  struct sockaddr_in client;
  int sin_size;
  sin_size = sizeof(struct sockaddr_in);
  int socket_count = 0;

  for (int i = 0; i < 100; i += 1) {
    for (int j = 0; j < 100; j += 1) {
      users[i].id[j] = '\0';
    }
    users[i].playwith = -1;
  }

  server_socket = socknet(AF_INET, SOCK_STREAM, 0);
  exit_on_error(server_socket, "cannot create socket");
  int opt = SO_REUSEADDR;
  setsockopt(server_socket, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt));
  bzero(&server, sizeof(server));
  server.sin_family = AF_INET;
  server.sin_port = htons(SERVER_PORT);
  server.sin_addr.s_addr = htonl(INADDR_ANY);
  exit_on_error(bind(server_socket, (struct sockaddr *)&server, sizeof(struct sockaddr)), "bind failed");
  exit_on_error(listen(server_socket, LISTEN_BACKLOG), "listen failed");
  printf("Waiting for client....\n");

  int new_socket;
  while (true) {
    new_socket = accept(server_socket, (struct sockaddr *)&client, &sin_size);
    exit_on_error(new_socket, "accept failed");

    if (socket_count >= MAX_CLIENT_COUNT) {
      printf("no more client is allowed\n");
      close(new_socket);
    }

    for (int i = 0; i < MAX_CLIENT_COUNT; i += 1) {
      if (fdt[i] == 0) {
        fdt[i] = new_socket;
        break;
      }
    }
    pthread_t thread;
    pthread_create(&thread, NULL, (void *)pthread_service, &new_socket);
    socket_count += 1;
  }
  close(server_socket);
}

int find_fd(char *name) {
  for (int i = 0; i < 100; i += 1) {
    if (strcmp(name, users[i].id) == 0) {
      return i;
    }
  }
  return -1;
}