
#!/bin/sh
ps ax |grep 'blog' | awk '{print $1}' |xargs kill -9

sleep 2

godep go build blog

chmod 777 blog

./blog &
#go run /apps/golang/src/blog/blog.go &
