
```bash

$ sudo apt install debootstrap 
$ mkdir rootfs
$ sudo debootstrap eoan rootfs/


$ man 2 clone
$ man 2 unshare

$ echo $$

go run -exec sudo main.go run --env "hello=container" --memory 50000000 bash

```