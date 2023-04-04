# 用于重置集群后清空整个磁盘，无法复原，谨慎操作
DISK="/dev/sdc"

sgdisk --zap-all $DISK

dd if=/dev/zero of="$DISK" bs=1M count=100 oflag=direct,dsync

blkdiscard $DISK

partprobe $DISK
