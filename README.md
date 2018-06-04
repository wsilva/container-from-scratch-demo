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

This first step will take a while at the first time. The clean Ubuntu Xenial VM will:

 - Install and configure Go (1.10.2)
 - Untar the root file systems of each image on /
    - we can check with ```ls / | grep rootfs```

If it fails due to internet conection we can re-run with 

~~~bash
$ vagrant up --provision
~~~

If the box become outdated we can just run

~~~bash
$ vagrant box update
~~~

### 2. Access the Ubuntu VM

Accessing it.

~~~bash
$ vagrant ssh
~~~

Run stuff as superuser and go to `/demo`.

~~~bash
vagrant@ubuntu-xenial:~$ sudo -i
root@ubuntu-xenial:~# cd /demo/
root@ubuntu-xenial:/demo#
~~~

### 3. First example - No isolation

Build and run an echo command from inside a container

~~~bash
root@ubuntu-xenial:/demo# go build demo.go
root@ubuntu-xenial:/demo# ./demo run echo "Qualquer coisa / Anything"
--Entrando no conteiner / Get into container--
--Rodando comando [echo Qualquer coisa / Anything] / Running command %!v(MISSING) --
Qualquer coisa / Anything
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo#
~~~

Run a bash with no isolation from inside a container

~~~bash
root@ubuntu-xenial:/demo# ./demo run /bin/bash
--Entrando no conteiner / Get into container--
--Rodando comando [/bin/bash] / Running command [/bin/bash] --
root@ubuntu-xenial:/demo# hostname
ubuntu-xenial
root@ubuntu-xenial:/demo# hostname demo
root@ubuntu-xenial:/demo# hostname
demo
root@ubuntu-xenial:/demo# exit
exit
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo# hostname
demo
root@ubuntu-xenial:/demo#
~~~

We can notice that hostname was changed also in the host

### 4. Second Example - Isolating UTS

Let's rebuild and run bash again but now with UTS (Unix Time Sharing) isolated

~~~bash
root@ubuntu-xenial:/demo# go build demo.go
root@ubuntu-xenial:/demo# ./demo run /bin/bash
--Entrando no conteiner / Get into container--
--Rodando comando [/bin/bash] / Running command [/bin/bash] --
root@ubuntu-xenial:/demo# hostname
demo
root@ubuntu-xenial:/demo# hostname yabadabadoo
root@ubuntu-xenial:/demo# hostname
yabadabadoo
root@ubuntu-xenial:/demo# exit
exit
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo# hostname
demo
root@ubuntu-xenial:/demo#
~~~

We can notice that hostname was modified inside container but not on the host.

###  5. Third Example - Defining a container hostname

Lets' rebuild and run bash again. Now we define a new hostname. But if we define before the ```cmd.Run()``` call we will define the host process not the container process. That's why we fork and exec so we can create a container with the run and then define the hostname inside the exec command.

~~~bash
root@ubuntu-xenial:/demo# go build demo.go
root@ubuntu-xenial:/demo# ./demo run bash
--Entrando no conteiner / Get into container--
--Rodando comando [bash] / Running command [bash] --
root@container:/demo# exit
exit
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo#
~~~

### 6. Fourth Example - define another rootfs

Let's rebuild and run choosing one of the 5 rootfs available: alpine, centos, debian, fedora or ubuntu.

~~~bash
root@ubuntu-xenial:/demo# go build demo.go
root@ubuntu-xenial:/demo# ./demo run debian bash
--Entrando no conteiner / Get into container--
--Imagem usada debian / Image in use debian --
--Rodando comando [bash] / Running command [bash] --
root@container:/# cat /etc/issue
Debian GNU/Linux 8 \n \l

root@container:/# exit
exit
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo#
root@ubuntu-xenial:/demo# ./demo run alpine sh
--Entrando no conteiner / Get into container--
--Imagem usada alpine / Image in use alpine --
--Rodando comando [sh] / Running command [sh] --
/ # cat /etc/issue
Welcome to Alpine Linux 3.4
Kernel \r on an \m (\l)
~~~


Befor leaving alpine container let's write a file and check if is written on the host also.

~~~bash
/ # touch /random-file
/ # exit
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo# ls /rootfs-alpine/
bin  dev  etc  home  lib  linuxrc  media  mnt  proc  random-file  root  run  sbin  srv  sys  tmp  usr  var
root@ubuntu-xenial:/demo#
~~~

Yes, the random-file exists inside and outside of the container.

### 7. Fifth Example - Isolate file system

Let's rebuild and run choosing one of the 5 rootfs available: alpine, centos, debian, fedora or ubuntu.

One folder will be created inside the container, let's create one file inside this folder and check if it appears on the host.

~~~bash
root@ubuntu-xenial:/demo# go build demo.go
root@ubuntu-xenial:/demo# ./demo run ubuntu bash
--Entrando no conteiner / Get into container--
--Imagem usada ubuntu / Image in use ubuntu --
--Rodando comando [bash] / Running command [bash] --
root@container:/# touch /mytemp/file-inside-container
root@container:/# exit
exit
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo# ls -l /rootfs-ubuntu/ | grep mytemp
root@ubuntu-xenial:/demo# ls -l /rootfs-ubuntu/mytemp/
total 0
root@ubuntu-xenial:/demo#
~~~

Yes, file was created and contained with filesystem isolation
Tip for next example run a container again and look for it's pids table:

~~~bash
root@container:/# ps aux
Error, do this: mount -t proc proc /proc
root@container:/# exit
exit
panic: exit status 47

goroutine 1 [running]:
main.doStuff(0x4ea240, 0xc420082060)
        /demo/demo.go:63 +0x4a
main.fork()
        /demo/demo.go:52 +0x2c4
main.main()
        /demo/demo.go:16 +0x6e
--Saindo do conteiner / Exiting container--
panic: exit status 2

goroutine 1 [running]:
main.doStuff(0x4ea240, 0xc42000a0e0)
        /demo/demo.go:63 +0x4a
main.run()
        /demo/demo.go:36 +0x4d2
main.main()
        /demo/demo.go:14 +0x9f
root@ubuntu-xenial:/demo#
~~~

Yes, we need to remount ```/proc```

### 8. Sixth Example - Isolate PID table

Let's rebuild and run a container, then we check if ```ps aux``` works.

~~~bash
root@ubuntu-xenial:/demo# go build demo.go
root@ubuntu-xenial:/demo# ./demo run ubuntu bash
--Entrando no conteiner / Get into container--
--Imagem usada ubuntu / Image in use ubuntu --
--Rodando comando [bash] / Running command [bash] --
root@container:/# ps aux
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.6  0.0   3376   976 ?        Sl   14:09   0:00 /proc/self/exe fork ubunt
root         5  0.0  0.3  18292  3244 ?        S    14:09   0:00 bash
root        10  0.0  0.2  34424  2880 ?        R+   14:09   0:00 ps aux
root@container:/# exit
exit
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo#
~~~

Yes we have only the pids from process running inside the container.

### 9. Seventh Example - Limiting processes running with cgroups

Let's rebuild and run a container. This time we need 2 terminal windows. Both must be on the vagrant box, at the first we will run a container.

On firts window:

~~~bash
root@ubuntu-xenial:/demo# go build demo.go
root@ubuntu-xenial:/demo# ./demo run ubuntu bash
--Entrando no conteiner / Get into container--
--Imagem usada ubuntu / Image in use ubuntu --
--Rodando comando [bash] / Running command [bash] --
root@container:/#
~~~

On the second one we will check how many process is the cgroup limit and which process are running 

On second window:

~~~bash
root@ubuntu-xenial:/demo# cat /sys/fs/cgroup/pids/demo/cgroup.procs
3233
3237
root@ubuntu-xenial:/demo# cat /sys/fs/cgroup/pids/demo/pids.max
15
root@ubuntu-xenial:/demo#
~~~

Then we will run a fork bomb and hope that cgroups limits the amount of forked process.

First window:

~~~bash
root@ubuntu-xenial:/demo# ./demo run ubuntu bash
--Entrando no conteiner / Get into container--
--Imagem usada ubuntu / Image in use ubuntu --
--Rodando comando [bash] / Running command [bash] --
root@container:/# :() { : | : & }; :
[1] 11
bash: fork: retry: No child processes
bash: fork: retry: No child processes
...
~~~

### 10. Navigating between tags

Each example is in a git tag between *demo1* and *demo7*, we can check the tags with: 

~~~bash
git tag --list
demo1
demo2
demo3
demo4
demo5
demo6
demo7
~~~

We can change for a specific tag with the following:

~~~bash
git checkout demo1
~~~

Getting back to master branch:

~~~bash
git checkout master
~~~

We also can check the differences between 2 tags with the following commands: 

~~~bash
git diff demo1..demo2
git diff demo2..demo3
~~~




### 11. Clean up

To leave VM shell

~~~bash
root@container:/# exit
exit
--Saindo do conteiner / Exiting container--
root@ubuntu-xenial:/demo# exit
logout
vagrant@ubuntu-xenial:~$ exit
logout
Connection to 127.0.0.1 closed.
[wsilva@localhost ]$
~~~

Then shut the VM down

~~~bash 
$ vagrant halt
==> containerdemo: Attempting graceful shutdown of VM...
[wsilva@localhost ]$
~~~