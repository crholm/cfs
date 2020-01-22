
Instructions
* Use a linux dist, preferable ubuntu
* Install go, https://golang.org/doc/install

Creating a rootfs
```bash
$ sudo apt install debootstrap 
$ mkdir rootfs
$ sudo debootstrap eoan rootfs/
```

Build allocate
```bash
$ cd utils
$ go build allocate.go
$ cp ./allocate ../rootfs
```

Syscall informaiton
```bash
$ man 2 clone
$ man 2 unshare
```

Running the "container"
```bash
$ go run -exec sudo main.go run --env "hello=container" --memory 50000000 bash 
$ echo $$    ## shows current pid
$ ps 
$ ps aux
$ echo $hello
$ env
$ ./allocate 100
```

