.PHONY: fix-bwrap run generate linux appimage deb windows macos macos-universal release

fix-bwrap:
	sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0

run: 
	wails3 dev

generate:
	wails3 generate bindings


appimage:
	wails3 task linux:create:appimage

deb:
	wails3 task linux:create:deb

linux: appimage deb

windows:
	wails3 task windows:build

macos:
	wails3 task darwin:package:universal

release: linux windows macos



# 	pgrep -af db-viewer
# ps -p 14045 -o pid,ppid,%cpu,%mem,rss,vsz,etime,cmd
# 
# Installation
# wails3 package GOOS=linux
# sudo apt install ./bin/*.deb