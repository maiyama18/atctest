# atctest

`atctest` is a command line tool for [AtCoder](https://atcoder.jp/).
it checks if your program correctly solves the problem for the sample inputs provided on the problem page.

## installation

```bash
$ go get -u github.com/mui87/atctest
```

## usage

```bash
atctest -contest ABC087 -problem A -command 'ruby abc/087/a.rb'
```

#### success case

![](https://user-images.githubusercontent.com/22269397/56220836-15505500-60a4-11e9-807b-26f0fff3d8c0.png)

#### failure case

![](https://user-images.githubusercontent.com/22269397/56220844-171a1880-60a4-11e9-883c-6211afc45d10.png)
