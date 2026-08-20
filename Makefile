fix-bwrap:
	sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0

run: 
	wails3 dev

generate:
	wails3 generate bindings

# 	pgrep -af db-viewer
# ps -p 14045 -o pid,ppid,%cpu,%mem,rss,vsz,etime,cmd
# 
# Installation
# wails3 package GOOS=linux
# sudo apt install ./bin/*.deb