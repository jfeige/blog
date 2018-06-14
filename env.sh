
#!/bin/sh
ps ax |grep 'blog' | awk '{print $1}' |xargs kill -9

sleep 2

go run /apps/golang/src/blog/blog.go &