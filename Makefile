fix-bwrap:
	sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0

# 	pgrep -af db-viewer
# ps -p 14045 -o pid,ppid,%cpu,%mem,rss,vsz,etime,cmd
# 