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


