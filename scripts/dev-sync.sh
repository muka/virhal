#!/bin/sh

if ! [ -x "$(command -v inotifywait)" ]; then
  echo 'Error: inotifywait is not installed. Try installing `inotify-tools`, `sudo apt install inotify-tools` on ubuntu' >&2
  exit 1
fi

nodes="node1 node2 manager1 manager2"

appname=`basename $PWD`
app="./$appname"

sync_nodes() {
    for node in $nodes; do
        echo "Syncing $node"
        docker-machine ssh $node "killall $appname"
        docker-machine scp $app "$node:~/$appname"
        # docker-machine ssh $node "nohup ~/$appname agent --debug"
        docker-machine ssh $node "nohup ~/$appname agent --debug >> ~/$appname.log 2>&1 &"
    done
    echo ""
    echo "Watching $app"    
}

echo "Watching $app"

while true
do
    inotifywait -e close_write $app | sync_nodes
done
