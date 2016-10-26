# Building containers from scratch using Go Lang

This demo was made as a prove of concept and to easily show that containers is nothing more than Linux processes running with some isolation from the processes running on the host.


## Inspiration

I started to study Go Lang couple weeks ago and I am a huge fan of Docker. Watching some videos from [Golang UK Conf. 2016](http://golanguk.com) I saw a talk from @lizrice: [What is a container, really](https://www.youtube.com/watch?v=HPuvDm8IC-4)

## Thanks

Need to thank Liz Rice (@lizrice) and as she mention at the talk, Julian Friedman (@julz), that with some lines of code could show details from it.

## Links

 - Gist from Liz Rice: https://gist.github.com/lizrice/a5ef4d175fd0cd3491c7e8d716826d27
 - Gist from Julian Friedman: https://gist.github.com/julz/c0017fa7a40de0543001
 - Julian post at InfoQ: https://www.infoq.com/articles/build-a-container-golang
 - Liz Rice talk at Golang UK Conference 2016: https://www.youtube.com/watch?v=HPuvDm8IC-4
 - Liz Rice talk at Container Camp UK 2016: https://www.youtube.com/watch?v=Utf-A4rODH8

## How to run

### 1. Bring the VM up

~~~bash
$ vagrant up
~~~

This first step will take a while at the first time. The clean Ubuntu trusty VM will:

 - Install debootstrap, and another tools to generate rootfs
 - Install and configure Go
 - Create rootfs for ubuntu trusty
 - Create rootfs for debian jessie

If it fails due to internet conection we can re-run with 

~~~bash
$ vagrant up --provision
~~~

### 2. Access the Ubuntu VM

Accessing it.

~~~bash
$ vagrant ssh
~~~

Going to `/container-demo` and run stuff as superuser.

~~~bash
vagrant@vagrant-ubuntu-trusty-64:~$ sudo -i
root@vagrant-ubuntu-trusty-64:~# cd /container-demo/
root@vagrant-ubuntu-trusty-64:/container-demo#
~~~

### 3. First example

run an echo command from inside a container

~~~bash
root@vagrant-ubuntu-trusty-64:/container-demo# go run demo.go run echo ola
Rodando [echo ola]
ola
Saindo do Container
root@vagrant-ubuntu-trusty-64:/container-demo#
~~~

Run a bash with no isolation from inside a container

~~~bash
root@vagrant-ubuntu-trusty-64:/container-demo# go run demo.go run /bin/bash
Rodando [/bin/bash]
root@vagrant-ubuntu-trusty-64:/container-demo# hostname
vagrant-ubuntu-trusty-64
root@vagrant-ubuntu-trusty-64:/container-demo# hostname demo
root@vagrant-ubuntu-trusty-64:/container-demo# hostname
demo
root@vagrant-ubuntu-trusty-64:/container-demo# exit
exit
Saindo do Container
root@vagrant-ubuntu-trusty-64:/container-demo# hostname
demo
root@vagrant-ubuntu-trusty-64:/container-demo#
~~~

We can notice that hostname was changed also in the host

### 4. Second example

Runnig a bash with UTS namespace isolation

~~~bash
root@vagrant-ubuntu-trusty-64:/container-demo#
root@vagrant-ubuntu-trusty-64:/container-demo# go run demo-uts.go run /bin/bash
Rodando [/bin/bash]
root@demo:/container-demo# hostname
demo
root@demo:/container-demo# hostname yada
root@demo:/container-demo# hostname
yada
root@demo:/container-demo# exit
exit
Saindo do Container
root@vagrant-ubuntu-trusty-64:/container-demo# hostname
demo
root@vagrant-ubuntu-trusty-64:/container-demo#
~~~

We can notice some isolation since hostname was changed inside the container but not in the host. 

### 5. Third example

Running a bash with UTS and PID namespaces isolated

~~~bash
root@demo:/container-demo# go run demo-pid.go run trusty ps aux
Rodando [ps aux] as PID 1 usando imagem trusty
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
0            1  0.0  0.1   3240   956 ?        Sl+  19:23   0:00 /proc/self/exe
0            4  0.0  0.2  15568  1160 ?        R+   19:23   0:00 ps aux
Saindo do Container
root@demo:/container-demo# go run demo-pid.go run trusty cat /etc/issue
Rodando [cat /etc/issue] as PID 1 usando imagem trusty
Ubuntu 14.04 LTS \n \l

Saindo do Container
root@demo:/container-demo#
root@demo:/container-demo# go run demo-pid.go run trusty ls -al
Rodando [ls -al] as PID 1 usando imagem trusty
total 72
drwxr-xr-x 19 0 0 4096 May 13  2013 .
drwxr-xr-x 19 0 0 4096 May 13  2013 ..
drwxr-xr-x  2 0 0 4096 Apr 16  2014 bin
drwxr-xr-x  2 0 0 4096 Apr 10  2014 boot
...
drwxr-xr-x 11 0 0 4096 Apr 11  2014 var
Saindo do Container
root@demo:/container-demo#
root@demo:/container-demo#
root@demo:/container-demo# go run demo-pid.go run jessie ps aux
Rodando [ps aux] as PID 1 usando imagem jessie
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.0  0.1   3240   952 ?        Sl+  19:25   0:00 /proc/self/exe
root         4  0.0  0.9 4137480 4712 ?        R+   18:55   0:00 /usr/bin/qemu-a
Saindo do Container
root@demo:/container-demo#
root@demo:/container-demo# go run demo-pid.go run jessie cat /etc/issue
Rodando [cat /etc/issue] as PID 1 usando imagem jessie
Debian GNU/Linux 8 \n \l

Saindo do Container
root@demo:/container-demo#
root@demo:/container-demo# go run demo-pid.go run jessie ls -al
Rodando [ls -al] as PID 1 usando imagem jessie
total 76
drwxr-xr-x 20 root root 4096 Oct 24 18:43 .
drwxr-xr-x 20 root root 4096 Oct 24 18:43 ..
drwxr-xr-x  2 root root 4096 Oct 24 18:40 bin
...
drwxr-xr-x 11 root root 4096 Oct 24 18:37 var
Saindo do Container
root@demo:/container-demo#
~~~

We can notice that the first proccess running inside the container have *PID=1* due to a new PID namespace

We can also show the rootfs from each image in the host (my VM in this case)

~~~bash
root@demo:/container-demo# ls -al / | grep rootfs
drwxr-xr-x 20 root    root     4096 Oct 24 18:43 rootfs-jessie
drwxr-xr-x 19 root    root     4096 May 13  2013 rootfs-trusty
root@demo:/container-demo#
~~~


### 6. To shut down

To leave VM shell

~~~bash
root@demo:/container-demo# exit
logout
vagrant@vagrant-ubuntu-trusty-64:~$ exit
logout
Connection to 127.0.0.1 closed.
~~~

Then shut the VM down

~~~bash 
$ vagrant halt
==> containerdemo: Attempting graceful shutdown of VM...
~~~