### memory alert for cform

alert the memory resource for all containers

1 get total memory from /proc/memoinfo
2 read info from /sys/fs/cgroup/memory/docker/<containerip>/memory.stat
3 get the percentage of (cache+rss)/total
4 if the value > 80 than alert to page duty