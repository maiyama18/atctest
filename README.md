# atctest

`atctest` is a command line tool for [AtCoder](https://atcoder.jp/).
it checks if your program correctly solves the problem for the sample inputs provided on the problem page.

## installation

```bash
// with go get
$ go get -u github.com/mui87/atctest

// with homebrew
$ brew tap mui87/atctest
$ brew install atctest
```

## usage

### command

#### specify contest/problem/command

```bash
$ atctest -contest ABC087 -problem A -command 'ruby abc/087/a.rb'
```

#### specify problem url/command

```bash
$ atctest -url 'https://atcoder.jp/contests/abc087/tasks/abc087_a' -command 'ruby abc/087/a.rb'
```

#### multiple commands (useful when using compile languages)

```bash
$ atctest -contest ABC087 -problem A -command 'g++ abc/087/a.cpp; ./a.out'
```

### results

#### success case

![](https://user-images.githubusercontent.com/22269397/56220836-15505500-60a4-11e9-807b-26f0fff3d8c0.png)

#### failure case

![](https://user-images.githubusercontent.com/22269397/56220844-171a1880-60a4-11e9-883c-6211afc45d10.png)
